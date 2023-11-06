package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"sci-review/common"
	"sci-review/form"
	"sci-review/model"
	"sci-review/service"
)

type OrganizationHandler struct {
	OrganizationService *service.OrganizationService
}

func NewOrganizationHandler(organizationService *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{OrganizationService: organizationService}
}

func (oh *OrganizationHandler) Create(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	pageData := common.PageData{
		Title:  "Create Organization",
		Active: "organizations",
		User:   principal,
	}

	organizationCreateForm := new(form.OrganizationCreateForm)
	if err := c.ShouldBind(&organizationCreateForm); err != nil {
		slog.Warn("organization create", "error", err.Error())
		pageData.Message = "Invalid form data"
		c.HTML(200, "organizations/create.html", gin.H{
			"pageData":         pageData,
			"organizationForm": organizationCreateForm,
		})
		return
	}
	slog.Info("organization create", "data", organizationCreateForm)

	if err := common.Validate(organizationCreateForm); len(err) > 0 {
		slog.Warn("organization create", "error", "validation error")
		pageData.Errors = err
		c.HTML(200, "organizations/create.html", gin.H{
			"pageData":         pageData,
			"organizationForm": organizationCreateForm,
		})
		return
	}

	_, err := oh.OrganizationService.Create(*organizationCreateForm, principal.Id)
	if err != nil {
		pageData.Message = "Service error"
		c.HTML(409, "organizations/create.html", gin.H{
			"pageData":         pageData,
			"organizationForm": organizationCreateForm,
		})
		return
	}

	c.Redirect(302, "/organizations")
}

func (oh *OrganizationHandler) List(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	pageData := common.PageData{
		Title:  "Organizations",
		Active: "organizations",
		User:   principal,
	}

	organizations, err := oh.OrganizationService.List(principal.Id)
	if err != nil {
		c.JSON(409, gin.H{"error": err.Error()})
		return
	}

	c.HTML(200, "organizations/index.html", gin.H{
		"pageData":      pageData,
		"organizations": organizations,
	})

}

func (oh *OrganizationHandler) Get(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)

	pageData := common.PageData{
		Title:  "Organization",
		Active: "organizations",
		User:   principal,
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Redirect(302, "/organizations")
		return
	}

	organization, err := oh.OrganizationService.Get(id, principal.Id)
	if err != nil {
		c.Redirect(302, "/organizations")
		return
	}

	c.HTML(200, "organizations/show.html", gin.H{
		"pageData":     pageData,
		"organization": organization,
	})
}

func (oh *OrganizationHandler) Archive(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)

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

func (oh *OrganizationHandler) CreateForm(c *gin.Context) {
	pageData := common.PageData{
		Title:  "Create Organization",
		Active: "organizations",
		User:   c.MustGet("principal").(*model.Principal),
	}
	c.HTML(200, "organizations/create.html", gin.H{
		"pageData": pageData,
	})
}

func RegisterOrganizationHandler(r *gin.Engine, organizationService *service.OrganizationService, middleware gin.HandlerFunc) {
	slog.Info("organization handler", "status", "registering")
	handler := NewOrganizationHandler(organizationService)
	r.GET("/organizations/new", middleware, handler.CreateForm)
	r.POST("/organizations", middleware, handler.Create)
	r.GET("/organizations", middleware, handler.List)
	r.GET("/organizations/:id", middleware, handler.Get)
	r.POST("/organizations/:id/archive", middleware, handler.Archive)
	slog.Info("organization handler", "status", "registered")
}
