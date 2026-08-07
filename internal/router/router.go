package router

import (
	"github.com/H0wZy/authsys/internal/handler"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	User *handler.UserHandler
	Auth *handler.AuthHandler
}

func New(h Handlers) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/user", h.User.Create)
		v1.POST("/auth/login", h.Auth.Login)
		v1.POST("/auth/logout", h.Auth.Logout)
	}

	return r
}
