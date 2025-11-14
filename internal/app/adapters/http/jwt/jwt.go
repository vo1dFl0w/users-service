package jwt

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	jwt_usecase "github.com/vo1dFl0w/users-service/internal/app/usecase/jwt_usecase"
)

type JWTService struct {
	secret []byte
}

func New(secret []byte) *JWTService {
	return &JWTService{
		secret: secret,
	}
}

func (s *JWTService) GenerateAccessToken(userID uuid.UUID) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt_usecase.TokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	return t.SignedString(s.secret)
}

func (s *JWTService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", string(b)), nil
}

func (s *JWTService) ValidateAccessToken(ctx context.Context, token string) (*jwt_usecase.TokenClaims, error) {
	c := &jwt_usecase.TokenClaims{}

	t, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		var e *jwt.ValidationError
		if errors.As(err, &e) && e.Errors&jwt.ValidationErrorExpired != 0 {
			return nil, fmt.Errorf("token expired")
		} else {
			return nil, fmt.Errorf("invalid token")
		}
	}

	if !t.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return c, nil
}
