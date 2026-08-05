package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtManager interface {
	Generate(id uint, email string, username string) (string, error)
	Validate(token string) (*Claims, error)
}

type Claims struct {
	Id       uint   `json:"sub"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type jwtManager struct {
	secret     []byte
	expiration time.Duration
}

func (j *jwtManager) Generate(id uint, email string, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      id,
		"email":    email,
		"username": username,
		"exp":      time.Now().Add(j.expiration).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *jwtManager) Validate(token string) (*Claims, error) {
	panic("unimplemented")
}

func NewJwtManager(secret []byte, expiration time.Duration) JwtManager {
	return &jwtManager{
		secret:     secret,
		expiration: expiration,
	}
}
