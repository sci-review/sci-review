package handler

import (
	"errors"
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
	loginForm := new(form.LoginForm)
	if err := c.ShouldBind(&loginForm); err != nil {
		slog.Warn("login", "error", err.Error())
		c.JSON(400, common.InvalidJson())
		return
	}
	slog.Info("login", "data", loginForm)

	if err := common.Validate(loginForm); len(err) > 0 {
		slog.Warn("login", "error", "validation error")
		c.JSON(400, common.ProblemWithErrors(err))
		return
	}

	loginAttemptData := form.LoginAttemptData{
		Email:     loginForm.Email,
		Password:  loginForm.Password,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}
	slog.Info("login", "data", loginAttemptData)

	tokenResponse, err := ah.AuthService.Login(loginAttemptData)
	if err != nil {
		slog.Warn("login", "error", err.Error())
		switch {
		case errors.Is(err, service.ErrorUserNotFound):
			c.JSON(409, common.NewProblemDetail("Invalid email or password", 409))
		case errors.Is(err, service.ErrorPasswordIsNotValid):
			c.JSON(409, common.NewProblemDetail("Invalid email or password", 409))
		case errors.Is(err, service.ErrorUserNotActive):
			c.JSON(409, common.NewProblemDetail("User not active", 409))
		case errors.Is(err, common.DbInternalError):
			c.JSON(500, common.NewInternalError())
		default:
			c.JSON(500, common.NewInternalError())
		}
		return
	}
	slog.Info("login", "result", "success")

	c.JSON(201, tokenResponse)
}

func (ah *AuthHandler) Logout(c *gin.Context) {
	logoutForm := new(form.LogoutForm)
	if err := c.ShouldBind(&logoutForm); err != nil {
		slog.Warn("logout", "error", err.Error())
		c.JSON(400, common.InvalidJson())
		return
	}

	if err := common.Validate(logoutForm); len(err) > 0 {
		slog.Warn("logout", "error", "validation error")
		c.JSON(400, common.ProblemWithErrors(err))
		return
	}

	err := ah.AuthService.Logout(logoutForm)
	if err != nil {
		slog.Warn("logout", "error", err.Error())
		switch {
		case errors.Is(err, service.ErrorParsingToken):
			c.JSON(400, common.NewProblemDetail("Error parsing token", 400))
		case errors.Is(err, service.ErrorRefreshTokenNotFound):
			c.JSON(409, common.NewProblemDetail("Refresh token not found", 409))
		case errors.Is(err, common.DbInternalError):
			c.JSON(500, common.NewInternalError())
		default:
			c.JSON(500, common.NewInternalError())
		}
		return
	}
	slog.Info("login", "result", "success")

	c.JSON(204, nil)
}

func RegisterAuthHandler(r *gin.Engine, authService *service.AuthService) {
	slog.Info("middleware handler", "status", "registering")
	authHandler := NewAuthHandler(authService)
	r.POST("/api/login", authHandler.Login)
	r.POST("/api/logout", authHandler.Logout)
	slog.Info("middleware handler", "status", "registered")
}
