package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"sci-review/common"
	"sci-review/model"
	"sci-review/service"
)

type ProfileHandler struct {
	userService *service.UserService
}

func (h ProfileHandler) Show(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(400, common.InvalidJson())
		return
	}

	user, err := h.userService.FindById(principal.Id, userId)
	if err != nil {
		if errors.Is(service.ErrorUserNotFound, err) {
			c.JSON(404, common.NewProblemDetail("User not found", 404))
		}
		if errors.Is(service.ErrorUserNotActive, err) {
			c.JSON(409, common.NewProblemDetail("User not active", 409))
		}
		c.JSON(500, common.NewProblemDetail("Internal error", 500))
		return
	}

	c.JSON(200, user)
}

func NewProfileHandler(userService *service.UserService) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
	}
}

func RegisterProfileHandler(r *gin.Engine, userService *service.UserService, tokenMiddleware gin.HandlerFunc) {
	slog.Info("profile handler", "status", "registering")
	profileHandler := NewProfileHandler(userService)
	r.GET("/api/user/:userId", tokenMiddleware, profileHandler.Show)
	slog.Info("profile handler", "status", "registered")
}
