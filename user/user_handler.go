package user

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"sci-review/common"
)

type UserHandler struct {
	UserService *UserService
}

func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (uh *UserHandler) Create(c *gin.Context) {
	userCreateForm := new(UserCreateForm)
	if err := c.ShouldBindJSON(&userCreateForm); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		slog.Warn("user create", "error", err.Error())
		return
	}
	slog.Info("user create", "data", userCreateForm)

	if err := common.Validate(userCreateForm); len(err) > 0 {
		c.JSON(400, gin.H{"errors": err})
		slog.Warn("user create", "error", "validation error")
		return
	}

	user, err := uh.UserService.Create(*userCreateForm)
	if err != nil {
		c.JSON(409, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, user)
}

func Register(r *gin.Engine, userService *UserService) {
	slog.Info("user handler", "status", "registering")
	userHandle := NewUserHandler(userService)
	r.POST("/users", userHandle.Create)
	slog.Info("user handler", "status", "registered")
}
