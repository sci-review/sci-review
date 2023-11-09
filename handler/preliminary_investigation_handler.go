package handler

import (
	"github.com/gin-gonic/gin"
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
	review := c.MustGet("review").(*model.Review)

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
	review := c.MustGet("review").(*model.Review)

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

	_, err := pi.PreliminaryInvestigationService.Create(*preliminaryInvestigationForm, review.Id, principal.Id)
	if err != nil {
		pageData.Message = err.Error()
		c.HTML(409, "preliminary_investigations/create.html", gin.H{
			"pageData":                     pageData,
			"preliminaryInvestigationForm": preliminaryInvestigationForm,
			"review":                       review,
		})
		return
	}

	c.Redirect(302, "/reviews/"+review.Id.String())
}

func (pi *PreliminaryInvestigationHandler) Show(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	review := c.MustGet("review").(*model.Review)
	investigation := c.MustGet("investigation").(*model.PreliminaryInvestigation)

	keywords, err := pi.PreliminaryInvestigationService.GetKeywordsByInvestigationId(investigation.Id)
	if err != nil {
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
		"keywords":      keywords,
	})
}

func (pi *PreliminaryInvestigationHandler) CreateKeyword(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	review := c.MustGet("review").(*model.Review)
	investigation := c.MustGet("investigation").(*model.PreliminaryInvestigation)

	keywords, err := pi.PreliminaryInvestigationService.GetKeywordsByInvestigationId(investigation.Id)
	if err != nil {
		return
	}

	pageData := common.PageData{
		Title:  "Preliminary Investigation",
		Active: "reviews",
		User:   principal,
	}

	keywordForm := new(form.KeywordForm)
	if err := c.ShouldBind(&keywordForm); err != nil {
		slog.Warn("investigation keyword create", "error", err.Error())
		pageData.Message = "Invalid form data"
		c.HTML(200, "preliminary_investigations/show.html", gin.H{
			"pageData":      pageData,
			"keywordForm":   keywordForm,
			"review":        review,
			"investigation": investigation,
			"keywords":      keywords,
		})
		return
	}
	slog.Info("investigation keyword create", "data", keywordForm)

	if err := common.Validate(keywordForm); len(err) > 0 {
		slog.Warn("investigation keyword create", "error", "validation error")
		pageData.Errors = err
		c.HTML(400, "preliminary_investigations/show.html", gin.H{
			"pageData":      pageData,
			"keywordForm":   keywordForm,
			"review":        review,
			"investigation": investigation,
			"keywords":      keywords,
		})
		return
	}

	err = pi.PreliminaryInvestigationService.SaveKeyword(investigation.Id, principal.Id, *keywordForm)
	if err != nil {
		pageData.Message = err.Error()
		c.HTML(409, "preliminary_investigations/show.html", gin.H{
			"pageData":      pageData,
			"keywordForm":   keywordForm,
			"review":        review,
			"investigation": investigation,
			"keywords":      keywords,
		})
		return
	}

	c.Redirect(302, "/reviews/"+review.Id.String()+"/preliminary_investigations/"+investigation.Id.String())

}

func RegisterPreliminaryInvestigationHandler(
	r *gin.Engine,
	reviewService *service.ReviewService,
	preliminaryInvestigationService *service.PreliminaryInvestigationService,
	authMiddleware gin.HandlerFunc,
	reviewMiddleware gin.HandlerFunc,
	investigationMiddleware gin.HandlerFunc,
) {
	preliminaryInvestigationHandler := NewPreliminaryInvestigationHandler(reviewService, preliminaryInvestigationService)
	r.GET(
		"/reviews/:reviewId/preliminary_investigations/create",
		authMiddleware,
		reviewMiddleware,
		preliminaryInvestigationHandler.CreateForm,
	)
	r.POST(
		"/reviews/:reviewId/preliminary_investigations/create",
		authMiddleware,
		reviewMiddleware,
		preliminaryInvestigationHandler.Create,
	)
	r.GET(
		"/reviews/:reviewId/preliminary_investigations/:investigationId",
		authMiddleware,
		reviewMiddleware,
		investigationMiddleware,
		preliminaryInvestigationHandler.Show,
	)
	r.POST(
		"/reviews/:reviewId/preliminary_investigations/:investigationId/keywords",
		authMiddleware,
		reviewMiddleware,
		investigationMiddleware,
		preliminaryInvestigationHandler.CreateKeyword,
	)
}
