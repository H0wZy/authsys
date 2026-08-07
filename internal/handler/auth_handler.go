package handler

import (
	"github.com/H0wZy/authsys/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	panic("unimplemented")
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	panic("unimplemented")
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}
