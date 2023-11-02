package user

import (
	"github.com/gin-gonic/gin"
	"sci-review/common"
)

type UserHandler struct {
	UserService *UserService
}

func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

type UserCreateForm struct {
	Name     string `json:"name" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email,max=350"`
	Password string `json:"password" validate:"required,min=6,max=60"`
}

func (uh *UserHandler) Create(c *gin.Context) {
	userCreateForm := new(UserCreateForm)
	if err := c.ShouldBindJSON(&userCreateForm); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := common.Validate(userCreateForm); len(err) > 0 {
		c.JSON(400, gin.H{"errors": err})
		return
	}

	user, err := uh.UserService.Create(*userCreateForm)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}
