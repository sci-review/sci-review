package handler

import (
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
	pageData := common.PageData{
		Title:  "Register",
		Active: "register",
	}
	userCreateForm := new(form.UserCreateForm)
	if err := c.ShouldBind(&userCreateForm); err != nil {
		slog.Warn("user create", "error", err.Error())
		pageData.Message = "Invalid form data"
		c.HTML(200, "users/register.html", gin.H{
			"pageData":       pageData,
			"userCreateForm": userCreateForm,
		})
		return
	}
	slog.Info("user create", "data", userCreateForm)

	if err := common.Validate(userCreateForm); len(err) > 0 {
		slog.Warn("user create", "error", "validation error")
		pageData.Errors = err
		c.HTML(200, "users/register.html", gin.H{
			"pageData":       pageData,
			"userCreateForm": userCreateForm,
		})
		return
	}

	_, err := uh.UserService.Create(*userCreateForm)
	if err != nil {
		pageData.Message = "Email already exists"
		c.HTML(409, "users/register.html", gin.H{
			"pageData":       pageData,
			"userCreateForm": userCreateForm,
		})
		return
	}
	c.Redirect(302, "/login?from=register")
}

func (uh *UserHandler) CreateForm(c *gin.Context) {
	pageData := common.PageData{
		Title:  "Register",
		Active: "register",
	}
	c.HTML(200, "users/register.html", gin.H{
		"pageData": pageData,
	})
}

func RegisterUserHandler(r *gin.Engine, userService *service.UserService) {
	slog.Info("user handler", "status", "registering")
	userHandle := NewUserHandler(userService)
	r.GET("/register", userHandle.CreateForm)
	r.POST("/register", userHandle.Create)
	slog.Info("user handler", "status", "registered")
}
