package main

import (
	"fmt"
	"os"
	"time"

	"github.com/H0wZy/authsys/internal/db"
	"github.com/H0wZy/authsys/internal/handler"
	"github.com/H0wZy/authsys/internal/jwt"
	"github.com/H0wZy/authsys/internal/repository"
	"github.com/H0wZy/authsys/internal/router"
	"github.com/H0wZy/authsys/internal/service"
)

func main() {
	conn, err := db.Connect()
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	secret := []byte(os.Getenv("JWT_SECRET"))
	jwtManager := jwt.NewJwtManager(secret, 24*time.Hour)

	userRepository := repository.NewUserRepository(conn)

	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	authService := service.NewAuthService(userRepository, jwtManager)
	authHandler := handler.NewAuthHandler(authService)

	h := router.Handlers{
		User: userHandler,
		Auth: authHandler,
	}

	r := router.New(h)
	if err := r.Run(); err != nil {
		panic(fmt.Errorf("failed to start server: %w", err))
	}
}
