package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const leeway = 30 * time.Second

var (
	ErrTokenInvalid = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type Manager interface {
	Generate(id uint, email, username string) (string, error)
	Validate(token string) (*Claims, error)
}

type Claims struct {
	UserID   uint
	Email    string
	Username string
}

type tokenClaims struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type manager struct {
	secret     []byte
	expiration time.Duration
	issuer     string
}

func (m *manager) Generate(id uint, email, username string) (string, error) {
	now := time.Now()

	claims := tokenClaims{
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   strconv.FormatUint(uint64(id), 10),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiration)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return token, nil
}

func (m *manager) Validate(token string) (*Claims, error) {
	var claims tokenClaims

	_, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(*jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(m.issuer),
		jwt.WithLeeway(leeway),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %w", ErrTokenInvalid, err)
	}

	id, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: malformed subject claim", ErrTokenInvalid)
	}

	return &Claims{
		UserID:   uint(id),
		Email:    claims.Email,
		Username: claims.Username,
	}, nil
}

func NewManager(secret []byte, expiration time.Duration, issuer string) Manager {
	return &manager{
		secret:     secret,
		expiration: expiration,
		issuer:     issuer,
	}
}
