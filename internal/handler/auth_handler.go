package handler

import (
	"errors"

	"github.com/H0wZy/authsys/internal/dto"
	"github.com/H0wZy/authsys/internal/response"
	"github.com/H0wZy/authsys/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var input dto.Auth

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(err)
		response.BadRequest(ctx, "invalid request body")
		return
	}

	pair, err := h.authService.Login(ctx.Request.Context(), &input)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			response.Unauthorized(ctx, "invalid credentials")

		default:
			ctx.Error(err)
			response.InternalServerError(ctx)
		}

		return
	}

	response.Ok(ctx, pair, "login successful!")
}

func (h *AuthHandler) Refresh(ctx *gin.Context) {
	var input dto.RefreshRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(err)
		response.BadRequest(ctx, "invalid request body")
		return
	}

	pair, err := h.authService.Refresh(ctx.Request.Context(), input.RefreshToken)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRefreshToken):
			response.Unauthorized(ctx, "invalid refresh token")

		default:
			ctx.Error(err)
			response.InternalServerError(ctx)
		}

		return
	}

	response.Ok(ctx, pair, "token refreshed!")
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	var input dto.RefreshRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(err)
		response.BadRequest(ctx, "invalid request body")
		return
	}

	if err := h.authService.Logout(ctx.Request.Context(), input.RefreshToken); err != nil {
		ctx.Error(err)
		response.InternalServerError(ctx)
		return
	}

	response.Ok(ctx, nil, "logout successful!")
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}
