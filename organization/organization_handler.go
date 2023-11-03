package organization

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"sci-review/auth"
	"sci-review/common"
)

type OrganizationHandler struct {
	OrganizationService *OrganizationService
}

func NewOrganizationHandler(organizationService *OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{OrganizationService: organizationService}
}

func (oh *OrganizationHandler) Create(c *gin.Context) {
	principal := c.MustGet("principal").(*auth.Principal)
	organizationCreateForm := new(OrganizationCreateForm)
	if err := c.ShouldBindJSON(&organizationCreateForm); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		slog.Warn("organization create", "error", err.Error())
		return
	}
	slog.Info("organization create", "data", organizationCreateForm)

	if err := common.Validate(organizationCreateForm); len(err) > 0 {
		c.JSON(400, gin.H{"errors": err})
		slog.Warn("organization create", "error", "validation error")
		return
	}

	user, err := oh.OrganizationService.Create(*organizationCreateForm, principal.Id)
	if err != nil {
		c.JSON(409, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, user)
}

func Register(r *gin.Engine, organizationService *OrganizationService, middleware gin.HandlerFunc) {
	slog.Info("organization handler", "status", "registering")
	handler := NewOrganizationHandler(organizationService)
	r.POST("/organizations", middleware, handler.Create)
	slog.Info("organization handler", "status", "registered")
}
