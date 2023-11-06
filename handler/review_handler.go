package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"sci-review/common"
	"sci-review/form"
	"sci-review/model"
	"sci-review/service"
)

type ReviewHandler struct {
	ReviewService *service.ReviewService
}

func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{ReviewService: reviewService}
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

	fmt.Println(reviews)

	c.HTML(200, "reviews/index.html", gin.H{
		"pageData": pageData,
		"reviews":  reviews,
	})
}

func RegisterReviewHandler(r *gin.Engine, reviewService *service.ReviewService, middleware gin.HandlerFunc) {
	reviewHandler := NewReviewHandler(reviewService)
	r.GET("/reviews", middleware, reviewHandler.Index)
	r.GET("/reviews/new", middleware, reviewHandler.CreateForm)
	r.POST("/reviews/new", middleware, reviewHandler.Create)
}
