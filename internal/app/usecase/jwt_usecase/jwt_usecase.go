package jwt

import (
	"context"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type TokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

type Service interface {
	GenerateAccessToken(userID uuid.UUID) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateAccessToken(ctx context.Context, token string) (*TokenClaims, error)
}
