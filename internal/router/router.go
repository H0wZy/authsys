package router

import (
	"log/slog"

	"github.com/H0wZy/authsys/internal/handler"
	"github.com/H0wZy/authsys/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	User *handler.UserHandler
	Auth *handler.AuthHandler
}

func New(h Handlers, logger *slog.Logger) *gin.Engine {
	r := gin.New()

	r.Use(middleware.Logger(logger))
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.POST("/user", h.User.Create)
		v1.POST("/auth/login", h.Auth.Login)
		v1.POST("/auth/refresh", h.Auth.Refresh)
		v1.POST("/auth/logout", h.Auth.Logout)
	}

	return r
}
