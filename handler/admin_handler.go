package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"sci-review/common"
	"sci-review/model"
	"sci-review/service"
)

type AdminHandler struct {
	UserService *service.UserService
}

func NewAdminHandler(userService *service.UserService) *AdminHandler {
	return &AdminHandler{UserService: userService}
}

func (ah *AdminHandler) Index(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)

	pageData := common.PageData{
		Title:  "Users",
		Active: "admin",
		User:   principal,
	}

	users, err := ah.UserService.FindAll(principal.Id)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	c.HTML(200, "admin/users.html", gin.H{
		"pageData": pageData,
		"users":    users,
	})
}

func (ah *AdminHandler) Activate(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = ah.UserService.Activate(principal.Id, id)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	c.Redirect(302, "/users")
}

func (ah *AdminHandler) Deactivate(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = ah.UserService.Deactivate(principal.Id, id)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	c.Redirect(302, "/users")
}

func RegisterAdminHandler(
	r *gin.Engine,
	userService *service.UserService,
	authMiddleware gin.HandlerFunc,
	adminMiddleware gin.HandlerFunc,
) {
	slog.Info("admin handler", "status", "registering")
	adminHandler := NewAdminHandler(userService)
	r.GET("/users", authMiddleware, adminMiddleware, adminHandler.Index)
	r.POST("/users/:id/activate", authMiddleware, adminMiddleware, adminHandler.Activate)
	r.POST("/users/:id/deactivate", authMiddleware, adminMiddleware, adminHandler.Deactivate)
	slog.Info("admin handler", "status", "registered")
}
