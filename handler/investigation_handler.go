package handler

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"sci-review/common"
	"sci-review/form"
	"sci-review/model"
	"sci-review/service"
)

type InvestigationHandler struct {
	ReviewService        *service.ReviewService
	InvestigationService *service.InvestigationService
}

func NewInvestigationHandler(reviewService *service.ReviewService, investigationService *service.InvestigationService) *InvestigationHandler {
	return &InvestigationHandler{ReviewService: reviewService, InvestigationService: investigationService}
}

func (pi *InvestigationHandler) CreateForm(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	review := c.MustGet("review").(*model.Review)

	pageData := common.PageData{
		Title:  "Create Investigation",
		Active: "reviews",
		User:   principal,
	}
	c.HTML(200, "investigations/create.html", gin.H{
		"pageData": pageData,
		"review":   review,
	})
}

func (pi *InvestigationHandler) Create(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	review := c.MustGet("review").(*model.Review)

	pageData := common.PageData{
		Title:  "Create Investigation",
		Active: "reviews",
		User:   principal,
	}

	investigationForm := new(form.InvestigationForm)
	if err := c.ShouldBind(&investigationForm); err != nil {
		slog.Warn("investigation create", "error", err.Error())
		pageData.Message = "Invalid form data"
		c.HTML(200, "investigations/create.html", gin.H{
			"pageData":          pageData,
			"investigationForm": investigationForm,
			"review":            review,
		})
		return
	}
	slog.Info("investigation create", "data", investigationForm)

	if err := common.Validate(investigationForm); len(err) > 0 {
		slog.Warn("investigation create", "error", "validation error")
		pageData.Errors = err
		c.HTML(400, "investigations/create.html", gin.H{
			"pageData":          pageData,
			"investigationForm": investigationForm,
			"review":            review,
		})
		return
	}

	_, err := pi.InvestigationService.Create(*investigationForm, review.Id, principal.Id)
	if err != nil {
		pageData.Message = err.Error()
		c.HTML(409, "investigations/create.html", gin.H{
			"pageData":          pageData,
			"investigationForm": investigationForm,
			"review":            review,
		})
		return
	}

	c.Redirect(302, "/reviews/"+review.Id.String())
}

func (pi *InvestigationHandler) Show(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	review := c.MustGet("review").(*model.Review)
	investigation := c.MustGet("investigation").(*model.Investigation)

	keywords, err := pi.InvestigationService.GetKeywordsByInvestigationId(investigation.Id)
	if err != nil {
		return
	}

	pageData := common.PageData{
		Title:  "Investigation",
		Active: "reviews",
		User:   principal,
	}

	c.HTML(200, "investigations/show.html", gin.H{
		"pageData":      pageData,
		"review":        review,
		"investigation": investigation,
		"keywords":      keywords,
	})
}

func (pi *InvestigationHandler) CreateKeyword(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	review := c.MustGet("review").(*model.Review)
	investigation := c.MustGet("investigation").(*model.Investigation)

	keywords, err := pi.InvestigationService.GetKeywordsByInvestigationId(investigation.Id)
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
		c.HTML(200, "investigations/show.html", gin.H{
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
		c.HTML(400, "investigations/show.html", gin.H{
			"pageData":      pageData,
			"keywordForm":   keywordForm,
			"review":        review,
			"investigation": investigation,
			"keywords":      keywords,
		})
		return
	}

	err = pi.InvestigationService.SaveKeyword(investigation.Id, principal.Id, *keywordForm)
	if err != nil {
		pageData.Message = err.Error()
		c.HTML(409, "investigations/show.html", gin.H{
			"pageData":      pageData,
			"keywordForm":   keywordForm,
			"review":        review,
			"investigation": investigation,
			"keywords":      keywords,
		})
		return
	}

	c.Redirect(302, "/reviews/"+review.Id.String()+"/investigations/"+investigation.Id.String())

}

func RegisterInvestigationHandler(
	r *gin.Engine,
	reviewService *service.ReviewService,
	investigationService *service.InvestigationService,
	authMiddleware gin.HandlerFunc,
	reviewMiddleware gin.HandlerFunc,
	investigationMiddleware gin.HandlerFunc,
) {
	investigationHandler := NewInvestigationHandler(reviewService, investigationService)
	r.GET(
		"/reviews/:reviewId/investigations/create",
		authMiddleware,
		reviewMiddleware,
		investigationHandler.CreateForm,
	)
	r.POST(
		"/reviews/:reviewId/investigations/create",
		authMiddleware,
		reviewMiddleware,
		investigationHandler.Create,
	)
	r.GET(
		"/reviews/:reviewId/investigations/:investigationId",
		authMiddleware,
		reviewMiddleware,
		investigationMiddleware,
		investigationHandler.Show,
	)
	r.POST(
		"/reviews/:reviewId/investigations/:investigationId/keywords",
		authMiddleware,
		reviewMiddleware,
		investigationMiddleware,
		investigationHandler.CreateKeyword,
	)
}
