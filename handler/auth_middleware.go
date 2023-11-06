package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"sci-review/model"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userId := session.Get("userId")
		userRole := session.Get("userRole")

		if userId == nil || userRole == nil {
			c.Redirect(302, "/login")
			return
		}

		user := model.NewPrincipal(userId.(string), userRole.(string))
		c.Set("principal", user)
		c.Next()
	}
}
