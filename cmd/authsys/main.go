package main

import (
	"fmt"
	"os"
	"time"

	"github.com/H0wZy/authsys/internal/db"
	"github.com/H0wZy/authsys/internal/jwt"
	"github.com/H0wZy/authsys/internal/repository"
	"github.com/H0wZy/authsys/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	conn, err := db.Connect()
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	secret := []byte(os.Getenv("JWT_SECRET"))
	jwtManager := jwt.NewJwtManager(secret, 24*time.Hour)

	userRepository := repository.NewUserRepository(conn)
	authService := service.NewAuthService(userRepository, jwtManager)
	userService := service.NewUserService(userRepository)

	r := gin.Default()

	r.Run(":8080")
}
