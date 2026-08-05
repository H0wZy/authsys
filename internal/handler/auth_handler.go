package handler

import (
	"github.com/H0wZy/authsys/internal/service"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authService service.AuthService
}

func (h *authHandler) Login(ctx *gin.Context) {
	panic("unimplemented")
}

func (h *authHandler) Logout(ctx *gin.Context) {
	panic("unimplemented")
}

func NewAuthHandler(authService service.AuthService) *authHandler {
	return &authHandler{authService: authService}
}
