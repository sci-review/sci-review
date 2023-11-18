package handler

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"sci-review/common"
	"sci-review/form"
	"sci-review/model"
	"sci-review/service"
)

type ReviewHandler struct {
	ReviewService        *service.ReviewService
	InvestigationService *service.InvestigationService
}

func NewReviewHandler(reviewService *service.ReviewService, investigationService *service.InvestigationService) *ReviewHandler {
	return &ReviewHandler{ReviewService: reviewService, InvestigationService: investigationService}
}

func (rh *ReviewHandler) Create(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)

	reviewForm := new(form.ReviewCreateForm)
	if err := c.ShouldBind(&reviewForm); err != nil {
		slog.Warn("review create", "error", err.Error())
		c.JSON(400, common.InvalidJson())
		return
	}
	slog.Info("review create", "data", reviewForm)

	if err := common.Validate(reviewForm); len(err) > 0 {
		slog.Warn("review create", "error", "validation error")
		c.JSON(400, common.ProblemWithErrors(err))
		return
	}

	review, err := rh.ReviewService.Create(*reviewForm, principal.Id)
	if err != nil {
		c.JSON(409, common.NewProblemDetail(err.Error(), 409))
		return
	}

	c.JSON(201, review)
}

func (rh *ReviewHandler) Index(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)

	reviews, err := rh.ReviewService.FindAll(principal.Id)
	if err != nil {
		c.JSON(500, common.NewProblemDetail(err.Error(), 500))
	}
	c.JSON(200, reviews)
}

func (rh *ReviewHandler) Show(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	review := c.MustGet("review").(*model.Review)

	_, err := rh.InvestigationService.FindAllByReviewID(review.Id, principal.Id)
	if err != nil {
		c.JSON(500, common.NewProblemDetail(err.Error(), 500))
	}

	//c.JSON(200, gin.H{
	//	"review":         review,
	//	"investigations": investigations,
	//})

	c.JSON(200, review)
}

func (rh *ReviewHandler) Update(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	review := c.MustGet("review").(*model.Review)

	reviewForm := new(form.ReviewCreateForm)
	if err := c.ShouldBind(&reviewForm); err != nil {
		slog.Warn("review update", "error", err.Error())
		c.JSON(400, common.InvalidJson())
		return
	}
	slog.Info("review update", "data", reviewForm)

	if err := common.Validate(reviewForm); len(err) > 0 {
		slog.Warn("review update", "error", "validation error")
		c.JSON(400, common.ProblemWithErrors(err))
		return
	}

	review, err := rh.ReviewService.Update(review.Id, *reviewForm, principal.Id)
	if err != nil {
		c.JSON(409, common.NewProblemDetail(err.Error(), 409))
		return
	}

	c.JSON(200, review)
}

func RegisterReviewHandler(
	r *gin.Engine,
	reviewService *service.ReviewService,
	investigationService *service.InvestigationService,
	tokenMiddleware gin.HandlerFunc,
	reviewMiddleware gin.HandlerFunc,
	investigationMiddleware gin.HandlerFunc,
) {
	reviewHandler := NewReviewHandler(reviewService, investigationService)
	r.GET("/api/reviews", tokenMiddleware, reviewHandler.Index)
	r.POST("/api/reviews", tokenMiddleware, reviewHandler.Create)
	r.GET("/api/reviews/:reviewId", tokenMiddleware, reviewMiddleware, reviewHandler.Show)
	r.PUT("/api/reviews/:reviewId", tokenMiddleware, reviewMiddleware, reviewHandler.Update)
}
