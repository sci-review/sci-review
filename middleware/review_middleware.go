package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sci-review/model"
	"sci-review/service"
)

func ReviewMiddleware(reviewService *service.ReviewService) gin.HandlerFunc {
	return func(c *gin.Context) {
		principal := c.MustGet("principal").(*model.Principal)

		id, err := uuid.Parse(c.Param("reviewId"))
		if err != nil {
			c.Redirect(302, "/reviews")
			c.Abort()
			return
		}

		review, err := reviewService.FindById(id, principal.Id)
		if err != nil {
			c.Redirect(302, "/reviews")
			c.Abort()
			return
		}

		c.Set("review", review)
		c.Next()
	}
}
