package auth

import (
	"github.com/gin-gonic/gin"
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

func (ah *AuthHandler) Login(c *gin.Context) {
	loginForm := new(LoginForm)
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := common.Validate(loginForm); len(err) > 0 {
		c.JSON(400, gin.H{"errors": err})
		return
	}

	loginAttemptData := LoginAttemptData{
		Email:     loginForm.Email,
		Password:  loginForm.Password,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}

	tokenResponse, err := ah.AuthService.Login(loginAttemptData)
	if err != nil {
		c.JSON(409, gin.H{"error": "Invalid credentials"})
		return
	}
	c.JSON(200, tokenResponse)
}

func Register(r *gin.Engine, authService *AuthService) {
	authHandler := NewAuthHandler(authService)
	r.POST("/login", authHandler.Login)
}
