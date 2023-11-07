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

type ReviewHandler struct {
	ReviewService                   *service.ReviewService
	PreliminaryInvestigationService *service.PreliminaryInvestigationService
}

func NewReviewHandler(reviewService *service.ReviewService, preliminaryInvestigationService *service.PreliminaryInvestigationService) *ReviewHandler {
	return &ReviewHandler{ReviewService: reviewService, PreliminaryInvestigationService: preliminaryInvestigationService}
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

	reviews, err := rh.ReviewService.GetByUserId(principal.Id)
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

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Redirect(302, "/reviews")
		return
	}

	review, err := rh.ReviewService.GetById(id, principal.Id)
	if err != nil {
		c.Redirect(302, "/reviews")
		return
	}

	investigations, err := rh.PreliminaryInvestigationService.GetAllByReviewID(review.Id)
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
	preliminaryInvestigationService *service.PreliminaryInvestigationService,
	middleware gin.HandlerFunc,
) {
	reviewHandler := NewReviewHandler(reviewService, preliminaryInvestigationService)
	r.GET("/reviews", middleware, reviewHandler.Index)
	r.GET("/reviews/new", middleware, reviewHandler.CreateForm)
	r.POST("/reviews/new", middleware, reviewHandler.Create)
	r.GET("/reviews/:id", middleware, reviewHandler.Show)
}
