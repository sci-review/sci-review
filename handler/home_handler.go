package handler

import "github.com/gin-gonic/gin"

type HomeHandler struct {
}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (hh *HomeHandler) Index(c *gin.Context) {
	_, exists := c.Get("principal")

	if exists {
		c.Redirect(302, "/reviews")
		return
	}

	c.Redirect(302, "/login")
}

func RegisterHomeHandler(r *gin.Engine, middleware gin.HandlerFunc) {
	hh := NewHomeHandler()

	r.GET("/", middleware, hh.Index)
}
