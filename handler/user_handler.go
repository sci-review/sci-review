package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"sci-review/common"
	"sci-review/form"
	"sci-review/service"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (uh *UserHandler) Create(c *gin.Context) {
	userCreateForm := new(form.UserCreateForm)
	if err := c.ShouldBind(&userCreateForm); err != nil {
		slog.Warn("user create", "error", err.Error())
		c.JSON(400, common.InvalidJson())
		return
	}
	slog.Info("user create", "data", userCreateForm)

	if err := common.Validate(userCreateForm); len(err) > 0 {
		slog.Warn("user create", "error", "validation error")
		c.JSON(400, common.ProblemWithErrors(err))
		return
	}

	user, err := uh.UserService.Create(*userCreateForm)
	if err != nil {
		if errors.Is(common.DbInternalError, err) {
			c.JSON(500, common.NewProblemDetail("Database internal error", 500))
			return
		}
		if errors.Is(service.ErrorUserAlreadyExists, err) {
			c.JSON(409, common.NewProblemDetail("Account with this e-mail already exists", 409))
			return
		}
		c.JSON(500, common.NewProblemDetail("Internal error", 500))
		return
	}

	c.JSON(201, user)
}

func RegisterUserHandler(r *gin.Engine, userService *service.UserService) {
	slog.Info("user handler", "status", "registering")
	userHandle := NewUserHandler(userService)
	r.POST("/api/register", userHandle.Create)
	slog.Info("user handler", "status", "registered")
}
