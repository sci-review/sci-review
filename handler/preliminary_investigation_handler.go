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

type PreliminaryInvestigationHandler struct {
	ReviewService                   *service.ReviewService
	PreliminaryInvestigationService *service.PreliminaryInvestigationService
}

func NewPreliminaryInvestigationHandler(reviewService *service.ReviewService, preliminaryInvestigationService *service.PreliminaryInvestigationService) *PreliminaryInvestigationHandler {
	return &PreliminaryInvestigationHandler{ReviewService: reviewService, PreliminaryInvestigationService: preliminaryInvestigationService}
}

func (pi *PreliminaryInvestigationHandler) CreateForm(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)

	reviewId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	review, err := pi.ReviewService.GetById(reviewId, principal.Id)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	pageData := common.PageData{
		Title:  "Create Preliminary Investigation",
		Active: "reviews",
		User:   principal,
	}
	c.HTML(200, "preliminary_investigations/create.html", gin.H{
		"pageData": pageData,
		"review":   review,
	})
}

func (pi *PreliminaryInvestigationHandler) Create(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)

	reviewId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	review, err := pi.ReviewService.GetById(reviewId, principal.Id)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	pageData := common.PageData{
		Title:  "Create Preliminary Investigation",
		Active: "reviews",
		User:   principal,
	}

	preliminaryInvestigationForm := new(form.PreliminaryInvestigationForm)
	if err := c.ShouldBind(&preliminaryInvestigationForm); err != nil {
		slog.Warn("preliminary investigation create", "error", err.Error())
		pageData.Message = "Invalid form data"
		c.HTML(200, "preliminary_investigations/create.html", gin.H{
			"pageData":                     pageData,
			"preliminaryInvestigationForm": preliminaryInvestigationForm,
			"review":                       review,
		})
		return
	}
	slog.Info("preliminary investigation create", "data", preliminaryInvestigationForm)

	if err := common.Validate(preliminaryInvestigationForm); len(err) > 0 {
		slog.Warn("preliminary investigation create", "error", "validation error")
		pageData.Errors = err
		c.HTML(400, "preliminary_investigations/create.html", gin.H{
			"pageData":                     pageData,
			"preliminaryInvestigationForm": preliminaryInvestigationForm,
			"review":                       review,
		})
		return
	}

	_, err = pi.PreliminaryInvestigationService.Create(*preliminaryInvestigationForm, reviewId, principal.Id)
	if err != nil {
		pageData.Message = err.Error()
		c.HTML(409, "preliminary_investigations/create.html", gin.H{
			"pageData":                     pageData,
			"preliminaryInvestigationForm": preliminaryInvestigationForm,
			"review":                       review,
		})
		return
	}

	c.Redirect(302, "/reviews/"+reviewId.String())
}

func (pi *PreliminaryInvestigationHandler) Show(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)

	reviewId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	investigationId, err := uuid.Parse(c.Param("investigationId"))
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	review, err := pi.ReviewService.GetById(reviewId, principal.Id)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	investigation, err := pi.PreliminaryInvestigationService.GetById(investigationId, principal.Id)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	pageData := common.PageData{
		Title:  "Preliminary Investigation",
		Active: "reviews",
		User:   principal,
	}

	c.HTML(200, "preliminary_investigations/show.html", gin.H{
		"pageData":      pageData,
		"review":        review,
		"investigation": investigation,
	})
}

func RegisterPreliminaryInvestigationHandler(
	r *gin.Engine,
	reviewService *service.ReviewService,
	preliminaryInvestigationService *service.PreliminaryInvestigationService,
	middleware gin.HandlerFunc,
) {
	preliminaryInvestigationHandler := NewPreliminaryInvestigationHandler(reviewService, preliminaryInvestigationService)
	r.GET("/reviews/:id/preliminary_investigations/create", middleware, preliminaryInvestigationHandler.CreateForm)
	r.POST("/reviews/:id/preliminary_investigations/create", middleware, preliminaryInvestigationHandler.Create)
	r.GET("/reviews/:id/preliminary_investigations/:investigationId", middleware, preliminaryInvestigationHandler.Show)
}
