package auth_usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/users-service/internal/app/adapters/storage"
	"github.com/vo1dFl0w/users-service/internal/app/domain/auth_domain"
	"github.com/vo1dFl0w/users-service/internal/app/usecase/jwt_usecase"
)

type Service interface {
	CreateUser(ctx context.Context, email string, password string) error
	GetUser(ctx context.Context, email string, password string) (*auth_domain.User, error)
	// TODO UpdateUser(ctx context.Context, user *auth.User) error
	// TODO DeleteUser(ctx context.Context, user *auth.User) error
	IssueTokens(ctx context.Context, userID uuid.UUID) (accessToken string, refreshToken string, err error)
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiry time.Time) error
	// TODO GetRefreshToken(ctx context.Context, userID uuid.UUID) (*auth_domain.User, error)
}

type service struct {
	storage storage.Storage
	jwt jwt.Service
}

func NewService(storage storage.Storage, jwtService jwt.Service) Service {
	return &service{
		storage: storage,
		jwt: jwtService,
	}
}

func (s *service) CreateUser(ctx context.Context, email string, password string) error {
	u, err := auth_domain.NewUser(email, password)
	if err != nil {
		return fmt.Errorf("failed to create new user: %w", err)
	}

	return s.storage.Auth().CreateUser(ctx, u.Email, u.EncryptedPassword)
}

func (s *service) GetUser(ctx context.Context, email string, password string) (*auth_domain.User, error) {
	u := &auth_domain.User{
		Email: email,
		Password: password,
	}

	if err := u.ValidateUser(); err != nil {
		return nil, err
	}

	u, err := s.storage.Auth().GetUser(ctx, u.Email)
	if err != nil {
		return nil, err
	}

	if !u.ComparePassword(password) {
		return nil, fmt.Errorf("wrong email or password")
	}

	u.Password = ""
	u.EncryptedPassword = ""

	return u, nil
}

// TODO func (s *service) UpdateUser(ctx context.Context, user *auth.User) error {return nil}

// TODO func (s *service) DeleteUser(ctx context.Context, user *auth.User) error {return nil}

func (s *service) IssueTokens(ctx context.Context, userID uuid.UUID) (accessToken string, refreshToken string, err error) {
	accessToken, err = s.jwt.GenerateAccessToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access tocken: %w", err)
	}

	refreshToken, err = s.jwt.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access tocken: %w", err)
	}

	hashedToken := sha256.Sum256([]byte(refreshToken))
	expiry := time.Now().Add(30 * 24 * time.Hour)

	if err := s.storage.Auth().SaveRefreshToken(ctx, userID, hex.EncodeToString(hashedToken[:]), expiry); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *service) SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiry time.Time) error {
	return s.storage.Auth().SaveRefreshToken(ctx, userID, token , expiry)
}

/* TODO func (s *service) GetRefreshToken(ctx context.Context, userID uuid.UUID) (*auth_domain.User, error) {
	return s.storage.Auth().GetRefreshToken(ctx, userID)
}
*/


