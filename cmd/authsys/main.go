package main

import (
	"log/slog"
	"os"

	"github.com/H0wZy/authsys/internal/config"
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

	cfg, err := config.Load()
	if err != nil {
		logger.Error("invalid configuration", slog.Any("error", err))
		os.Exit(1)
	}

	conn, err := db.Connect(cfg.DatabaseDSN)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("database connected!")

	jwtManager := jwt.NewJwtManager(cfg.JWTSecret, cfg.JWTExpiration)

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

	if err := r.Run(cfg.ServerAddr); err != nil {
		logger.Error("failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
