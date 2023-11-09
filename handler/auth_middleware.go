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
			c.Abort()
			return
		}

		user := model.NewPrincipal(userId.(string), userRole.(string))
		c.Set("principal", user)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userRole := session.Get("userRole")

		if userRole == nil {
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		if userRole != model.UserAdmin {
			c.Redirect(302, "/reviews")
			c.Abort()
			return
		}

		c.Next()
	}
}
