package organization

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (oh *OrganizationHandler) List(c *gin.Context) {
	principal := c.MustGet("principal").(*auth.Principal)

	organizations, err := oh.OrganizationService.List(principal.Id)
	if err != nil {
		c.JSON(409, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, organizations)
}

func (oh *OrganizationHandler) Get(c *gin.Context) {
	principal := c.MustGet("principal").(*auth.Principal)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	organization, err := oh.OrganizationService.Get(id, principal.Id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, organization)
}

func (oh *OrganizationHandler) Archive(c *gin.Context) {
	principal := c.MustGet("principal").(*auth.Principal)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = oh.OrganizationService.Archive(id, principal.Id)
	if err != nil {
		c.JSON(409, gin.H{"error": err.Error()})
		return
	}

	c.JSON(204, nil)
}

func Register(r *gin.Engine, organizationService *OrganizationService, middleware gin.HandlerFunc) {
	slog.Info("organization handler", "status", "registering")
	handler := NewOrganizationHandler(organizationService)
	r.POST("/organizations", middleware, handler.Create)
	r.GET("/organizations", middleware, handler.List)
	r.GET("/organizations/:id", middleware, handler.Get)
	r.POST("/organizations/:id/archive", middleware, handler.Archive)
	slog.Info("organization handler", "status", "registered")
}
