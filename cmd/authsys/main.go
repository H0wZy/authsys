package main

import (
	"log/slog"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	conn, err := db.Connect()
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("database connected!")

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

	r := router.New(h, logger)

	if err := r.Run(); err != nil {
		logger.Error("failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
