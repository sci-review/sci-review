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

func (rh *ReviewHandler) CreateForm(c *gin.Context) {
	pageData := common.PageData{
		Title:  "Create Review",
		Active: "reviews",
		User:   c.MustGet("principal").(*model.Principal),
	}
	c.HTML(200, "reviews/create.html", gin.H{
		"pageData": pageData,
	})
}

func (rh *ReviewHandler) Create(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	pageData := common.PageData{
		Title:  "Create Review",
		Active: "reviews",
		User:   principal,
	}

	reviewForm := new(form.ReviewCreateForm)
	if err := c.ShouldBind(&reviewForm); err != nil {
		slog.Warn("review create", "error", err.Error())
		pageData.Message = "Invalid form data"
		c.HTML(200, "reviews/create.html", gin.H{
			"pageData":   pageData,
			"reviewForm": reviewForm,
		})
		return
	}
	slog.Info("review create", "data", reviewForm)

	if err := common.Validate(reviewForm); len(err) > 0 {
		slog.Warn("review create", "error", "validation error")
		pageData.Errors = err
		c.HTML(400, "reviews/create.html", gin.H{
			"pageData":   pageData,
			"reviewForm": reviewForm,
		})
		return
	}

	_, err := rh.ReviewService.Create(*reviewForm, principal.Id)
	if err != nil {
		pageData.Message = err.Error()
		c.HTML(409, "reviews/create.html", gin.H{
			"pageData":   pageData,
			"reviewForm": reviewForm,
		})
		return
	}

	c.Redirect(302, "/reviews")
}

func (rh *ReviewHandler) Index(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)

	pageData := common.PageData{
		Title:  "Reviews",
		Active: "reviews",
		User:   principal,
	}

	reviews, err := rh.ReviewService.FindAll(principal.Id)
	if err != nil {
		return
	}

	c.HTML(200, "reviews/index.html", gin.H{
		"pageData": pageData,
		"reviews":  reviews,
	})
}

func (rh *ReviewHandler) Show(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	review := c.MustGet("review").(*model.Review)

	investigations, err := rh.InvestigationService.FindAllByReviewID(review.Id)
	if err != nil {
		return
	}

	pageData := common.PageData{
		Title:  "Review",
		Active: "reviews",
		User:   principal,
	}

	c.HTML(200, "reviews/show.html", gin.H{
		"pageData":       pageData,
		"review":         review,
		"investigations": investigations,
	})
}

func RegisterReviewHandler(
	r *gin.Engine,
	reviewService *service.ReviewService,
	investigationService *service.InvestigationService,
	authMiddleware gin.HandlerFunc,
	reviewMiddleware gin.HandlerFunc,
	investigationMiddleware gin.HandlerFunc,
) {
	reviewHandler := NewReviewHandler(reviewService, investigationService)
	r.GET("/reviews", authMiddleware, reviewHandler.Index)
	r.GET("/reviews/new", authMiddleware, reviewHandler.CreateForm)
	r.POST("/reviews/new", authMiddleware, reviewHandler.Create)
	r.GET("/reviews/:reviewId", authMiddleware, reviewMiddleware, reviewHandler.Show)
}
