package user

import "github.com/gin-gonic/gin"

type UserHandler struct {
	UserService *UserService
}

func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

type UserCreateForm struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uh *UserHandler) Create(c *gin.Context) {
	userCreateForm := new(UserCreateForm)
	if err := c.ShouldBindJSON(&userCreateForm); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := uh.UserService.Create(*userCreateForm)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}
