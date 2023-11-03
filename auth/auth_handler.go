package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"sci-review/common"
)

type AuthHandler struct {
	AuthService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

type LoginForm struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginAttemptData struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	IPAddress string `json:"ipAddress"`
	UserAgent string `json:"userAgent"`
}

func (la LoginAttemptData) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("email", la.Email),
		slog.String("ipAddress", la.IPAddress),
		slog.String("userAgent", la.UserAgent),
	)
}

type RotateTokenData struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

func (ah *AuthHandler) Login(c *gin.Context) {
	loginForm := new(LoginForm)
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		slog.Warn("login", "error", "error binding json")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := common.Validate(loginForm); len(err) > 0 {
		slog.Warn("login", "error", "validation error")
		c.JSON(400, gin.H{"errors": err})
		return
	}

	loginAttemptData := LoginAttemptData{
		Email:     loginForm.Email,
		Password:  loginForm.Password,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}
	slog.Info("login", "data", loginAttemptData)

	tokenResponse, err := ah.AuthService.Login(loginAttemptData)
	if err != nil {
		slog.Warn("login", "error", err.Error())
		c.JSON(409, gin.H{"error": "Invalid credentials"})
		return
	}
	slog.Info("login", "result", "success")
	c.JSON(200, tokenResponse)
}

func (ah *AuthHandler) RotateToken(c *gin.Context) {
	rotateTokenData := new(RotateTokenData)
	if err := c.ShouldBindJSON(&rotateTokenData); err != nil {
		slog.Warn("rotate token", "error", "error binding json")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := common.Validate(rotateTokenData); len(err) > 0 {
		slog.Warn("rotate token", "error", "validation error")
		c.JSON(400, gin.H{"errors": err})
		return
	}

	tokenResponse, err := ah.AuthService.RotateToken(rotateTokenData)
	if err != nil {
		slog.Warn("rotate token", "error", err.Error())
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	slog.Info("rotate token", "result", "success")
	c.JSON(200, tokenResponse)
}

func Register(r *gin.Engine, authService *AuthService) {
	slog.Info("auth handler", "status", "registering")
	authHandler := NewAuthHandler(authService)
	r.POST("/login", authHandler.Login)
	r.POST("/rotate-token", authHandler.RotateToken)
	slog.Info("auth handler", "status", "registered")
}
