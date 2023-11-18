package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"sci-review/common"
	"sci-review/form"
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
			return
		}
		if errors.Is(service.ErrorUserNotActive, err) {
			c.JSON(409, common.NewProblemDetail("User not active", 409))
			return
		}
		c.JSON(500, common.NewProblemDetail("Internal error", 500))
		return
	}

	c.JSON(200, user)
}

func (h ProfileHandler) ChangePassword(c *gin.Context) {
	principal := c.MustGet("principal").(*model.Principal)
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(400, common.InvalidJson())
		return
	}

	passwordForm := new(form.ChangePasswordForm)
	if err := c.ShouldBind(&passwordForm); err != nil {
		slog.Warn("change password", "error", err.Error())
		c.JSON(400, common.InvalidJson())
		return
	}
	slog.Info("change password", "data", passwordForm)

	if err := common.Validate(passwordForm); len(err) > 0 {
		slog.Warn("change password", "error", "validation error")
		c.JSON(400, common.ProblemWithErrors(err))
		return
	}

	err = h.userService.ChangePassword(principal.Id, userId, passwordForm)
	if err != nil {
		if errors.Is(service.ErrorUserNotFound, err) {
			c.JSON(404, common.NewProblemDetail("User not found", 404))
			return
		}
		if errors.Is(service.ErrorUserNotActive, err) {
			c.JSON(409, common.NewProblemDetail("User not active", 409))
			return
		}
		if errors.Is(service.ErrorPasswordIsNotValid, err) {
			c.JSON(409, common.NewProblemDetail("Password is not valid", 409))
			return
		}
		c.JSON(500, common.NewProblemDetail("Internal error", 500))
		return
	}

	c.JSON(204, nil)
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
	r.PUT("/api/user/:userId/password", tokenMiddleware, profileHandler.ChangePassword)
	slog.Info("profile handler", "status", "registered")
}
