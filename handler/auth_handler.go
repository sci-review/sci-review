package handler

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"sci-review/common"
	"sci-review/form"
	"sci-review/service"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func (ah *AuthHandler) Login(c *gin.Context) {
	pageData := common.PageData{
		Title:   "Login",
		Active:  "login",
		Message: "",
	}

	loginForm := new(form.LoginForm)
	if err := c.ShouldBind(&loginForm); err != nil {
		pageData.Message = "Invalid form data"
		c.HTML(200, "users/login.html", gin.H{
			"pageData":  pageData,
			"loginForm": loginForm,
		})
		return
	}

	if err := common.Validate(loginForm); len(err) > 0 {
		slog.Warn("login", "error", "validation error")
		pageData.Errors = err
		c.HTML(200, "users/login.html", gin.H{
			"pageData":  pageData,
			"loginForm": loginForm,
		})
		return
	}

	loginAttemptData := form.LoginAttemptData{
		Email:     loginForm.Email,
		Password:  loginForm.Password,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}
	slog.Info("login", "data", loginAttemptData)

	userLogged, err := ah.AuthService.Login(loginAttemptData)
	if err != nil {
		slog.Warn("login", "error", err.Error())

		switch {
		case errors.Is(err, service.ErrorUserNotFound):
			pageData.Message = "Invalid email or password"
		case errors.Is(err, service.ErrorUserNotActive):
			pageData.Message = "User not active"
		case errors.Is(err, common.DbInternalError):
			pageData.Message = "Db Internal Error"
		default:
			pageData.Message = "Internal Error"
		}

		//if errors.Is(err, service.ErrorUserNotActive) {
		//	pageData.Message = "User not active"
		//} else {
		//	pageData.Message = "Invalid email or password"
		//}

		c.HTML(409, "users/login.html", gin.H{
			"pageData":  pageData,
			"loginForm": loginForm,
		})
		return
	}
	slog.Info("login", "result", "success")

	session := sessions.Default(c)
	session.Set("userId", userLogged.Id.String())
	session.Set("userRole", string(userLogged.Role))
	session.Save()

	c.Redirect(302, "reviews")
}

func (ah *AuthHandler) LoginForm(c *gin.Context) {
	from := c.Query("from")
	pageData := common.PageData{
		Title:   "Login",
		Active:  "login",
		Message: "",
	}
	c.HTML(200, "users/login.html", gin.H{
		"pageData": pageData,
		"from":     from,
	})
}

func (ah *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		slog.Error("logout", "error", err.Error())
		c.AbortWithStatus(500)
	}
	c.Redirect(302, "/login")
}

func RegisterAuthHandler(r *gin.Engine, authService *service.AuthService) {
	slog.Info("middleware handler", "status", "registering")
	authHandler := NewAuthHandler(authService)
	r.GET("/login", authHandler.LoginForm)
	r.POST("/login", authHandler.Login)
	r.GET("/logout", authHandler.Logout)
	slog.Info("middleware handler", "status", "registered")
}
