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

	token, err := h.authService.Login(ctx.Request.Context(), &input)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			response.Unauthorized(ctx, err.Error())

		default:
			ctx.Error(err)
			response.InternalServerError(ctx)
		}

		return
	}

	response.Ok(ctx, dto.AuthResponse{Token: token}, "login successfull!")
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	panic("unimplemented")
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}
