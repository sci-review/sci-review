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

type InvestigationHandler struct {
	ReviewService        *service.ReviewService
	InvestigationService *service.InvestigationService
}

func NewInvestigationHandler(reviewService *service.ReviewService, investigationService *service.InvestigationService) *InvestigationHandler {
	return &InvestigationHandler{ReviewService: reviewService, InvestigationService: investigationService}
}

func (pi *InvestigationHandler) Create(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	review := c.MustGet("review").(*model.Review)

	investigationForm := new(form.InvestigationForm)
	if err := c.ShouldBind(&investigationForm); err != nil {
		slog.Warn("investigation create", "error", err.Error())
		c.JSON(400, common.InvalidJson())
		return
	}
	slog.Info("investigation create", "data", investigationForm)

	if err := common.Validate(investigationForm); len(err) > 0 {
		slog.Warn("investigation create", "error", "validation error")
		c.JSON(400, common.ProblemWithErrors(err))
		return
	}

	investigation, err := pi.InvestigationService.Create(*investigationForm, review.Id, principal.Id)
	if err != nil {
		c.JSON(409, common.NewProblemDetail(err.Error(), 409))
		return
	}

	c.JSON(201, investigation)
}

func (pi *InvestigationHandler) Index(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	review := c.MustGet("review").(*model.Review)

	investigations, err := pi.InvestigationService.FindAllByReviewID(review.Id, principal.Id)
	if err != nil {
		c.JSON(500, common.NewProblemDetail(err.Error(), 500))
	}

	c.JSON(200, investigations)
}

func (pi *InvestigationHandler) Show(c *gin.Context) {
	//principal := c.MustGet("principal").(*model.Principal)
	//review := c.MustGet("review").(*model.Review)
	investigation := c.MustGet("investigation").(*model.Investigation)

	//keywords, err := pi.InvestigationService.GetKeywordsByInvestigationId(investigation.Id)
	//if err != nil {
	//	return
	//}

	//pageData := common.PageData{
	//	Title:  "Investigation",
	//	Active: "reviews",
	//	User:   principal,
	//}

	//c.HTML(200, "investigations/show.html", gin.H{
	//	"pageData":      pageData,
	//	"review":        review,
	//	"investigation": investigation,
	//	"keywords":      keywords,
	//})

	c.JSON(200, investigation)
}

func (pi *InvestigationHandler) CreateKeyword(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	//review := c.MustGet("review").(*model.Review)
	investigation := c.MustGet("investigation").(*model.Investigation)

	keywordForm := new(form.KeywordForm)
	if err := c.ShouldBind(&keywordForm); err != nil {
		slog.Warn("investigation keyword create", "error", err.Error())
		c.JSON(400, common.InvalidJson())
		return
	}
	slog.Info("investigation keyword create", "data", keywordForm)

	if err := common.Validate(keywordForm); len(err) > 0 {
		slog.Warn("investigation keyword create", "error", "validation error")
		c.JSON(400, common.ProblemWithErrors(err))
		return
	}

	keyword, err := pi.InvestigationService.SaveKeyword(investigation.Id, principal.Id, *keywordForm)
	if err != nil {
		c.JSON(409, common.NewProblemDetail(err.Error(), 409))
		return
	}

	c.JSON(201, keyword)

}

func (pi *InvestigationHandler) GetAllKeywords(c *gin.Context) {
	//principal := c.MustGet("principal").(*model.Principal)
	//review := c.MustGet("review").(*model.Review)
	investigation := c.MustGet("investigation").(*model.Investigation)

	keywords, err := pi.InvestigationService.GetKeywordsByInvestigationId(investigation.Id)
	if err != nil {
		c.JSON(500, common.NewProblemDetail(err.Error(), 500))
		return
	}

	c.JSON(200, keywords)
}

func (pi *InvestigationHandler) DeleteKeyword(c *gin.Context) {
	//principal := c.MustGet("principal").(*model.Principal)
	//review := c.MustGet("review").(*model.Review)
	//investigation := c.MustGet("investigation").(*model.Investigation)

	keywordId, err := uuid.Parse(c.Param("keywordId"))
	if err != nil {
		c.JSON(400, common.NewProblemDetail(err.Error(), 400))
		return
	}

	err = pi.InvestigationService.DeleteKeyword(keywordId)
	if err != nil {
		c.JSON(500, common.NewProblemDetail(err.Error(), 500))
		return
	}

	c.JSON(200, gin.H{})
}

func (pi *InvestigationHandler) UpdateKeyword(c *gin.Context) {
	//principal := c.MustGet("principal").(*model.Principal)
	//review := c.MustGet("review").(*model.Review)
	//investigation := c.MustGet("investigation").(*model.Investigation)

	keywordId, err := uuid.Parse(c.Param("keywordId"))
	if err != nil {
		c.JSON(400, common.NewProblemDetail(err.Error(), 400))
		return
	}

	keywordForm := new(form.KeywordForm)
	if err := c.ShouldBind(&keywordForm); err != nil {
		slog.Warn("investigation keyword update", "error", err.Error())
		c.JSON(400, common.InvalidJson())
		return
	}
	slog.Info("investigation keyword update", "data", keywordForm)

	if err := common.Validate(keywordForm); len(err) > 0 {
		slog.Warn("investigation keyword update", "error", "validation error")
		c.JSON(400, common.ProblemWithErrors(err))
		return
	}

	keyword, err := pi.InvestigationService.UpdateKeyword(keywordId, *keywordForm)
	if err != nil {
		c.JSON(409, common.NewProblemDetail(err.Error(), 409))
		return
	}

	c.JSON(200, keyword)
}

func RegisterInvestigationHandler(
	r *gin.Engine,
	reviewService *service.ReviewService,
	investigationService *service.InvestigationService,
	tokenMiddleware gin.HandlerFunc,
	reviewMiddleware gin.HandlerFunc,
	investigationMiddleware gin.HandlerFunc,
) {
	investigationHandler := NewInvestigationHandler(reviewService, investigationService)

	r.POST(
		"/api/reviews/:reviewId/investigations",
		tokenMiddleware,
		reviewMiddleware,
		investigationHandler.Create,
	)
	r.GET(
		"/api/reviews/:reviewId/investigations",
		tokenMiddleware,
		reviewMiddleware,
		investigationHandler.Index,
	)
	r.GET(
		"/api/reviews/:reviewId/investigations/:investigationId",
		tokenMiddleware,
		reviewMiddleware,
		investigationMiddleware,
		investigationHandler.Show,
	)
	r.POST(
		"/api/reviews/:reviewId/investigations/:investigationId/keywords",
		tokenMiddleware,
		reviewMiddleware,
		investigationMiddleware,
		investigationHandler.CreateKeyword,
	)
	r.GET(
		"/api/reviews/:reviewId/investigations/:investigationId/keywords",
		tokenMiddleware,
		reviewMiddleware,
		investigationMiddleware,
		investigationHandler.GetAllKeywords,
	)

	r.DELETE(
		"/api/reviews/:reviewId/investigations/:investigationId/keywords/:keywordId",
		tokenMiddleware,
		reviewMiddleware,
		investigationMiddleware,
		investigationHandler.DeleteKeyword,
	)

	r.PUT(
		"/api/reviews/:reviewId/investigations/:investigationId/keywords/:keywordId",
		tokenMiddleware,
		reviewMiddleware,
		investigationMiddleware,
		investigationHandler.UpdateKeyword,
	)

}
