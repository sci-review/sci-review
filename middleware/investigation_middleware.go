package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sci-review/model"
	"sci-review/service"
)

func InvestigationMiddleware(investigationService *service.PreliminaryInvestigationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		principal := c.MustGet("principal").(*model.Principal)

		investigationId, err := uuid.Parse(c.Param("investigationId"))
		if err != nil {
			review := c.MustGet("review").(*model.Review)
			c.Redirect(302, "/reviews/"+review.Id.String())
			c.Abort()
			return
		}
		investigation, err := investigationService.GetById(investigationId, principal.Id)
		if err != nil {
			review := c.MustGet("review").(*model.Review)
			c.Redirect(302, "/reviews/"+review.Id.String())
			c.Abort()
			return
		}

		c.Set("investigation", investigation)
		c.Next()
	}
}
