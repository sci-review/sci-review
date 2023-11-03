package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")

		if authorization == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		if strings.HasPrefix(authorization, "Bearer ") {
			splits := strings.Split(authorization, " ")
			if len(splits) != 2 {
				c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
				return
			}

			token, err := ParseToken(splits[1])
			if err != nil {
				c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
				return
			}

			sub := token.Claims.(jwt.MapClaims)["sub"].(string)
			role := token.Claims.(jwt.MapClaims)["role"].(string)
			principal := NewPrincipal(sub, role)
			c.Set("principal", principal)
			return
		}

		c.Next()
	}
}