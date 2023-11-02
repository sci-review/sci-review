package auth

import "github.com/gin-gonic/gin"

type AuthHandler struct {
	AuthService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ah *AuthHandler) Login(c *gin.Context) {
	loginForm := new(LoginForm)
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	tokenResponse, err := ah.AuthService.Login(*loginForm)
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
