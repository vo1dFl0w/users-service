package auth_domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, email string, password string) error
	GetUser(ctx context.Context, email string) (*User, error)
	// TODO UpdateUser(ctx context.Context, user *User) error
	// TODO DeleteUser(ctx context.Context, user *User) error
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiry time.Time) error
	// TODO GetRefreshToken(ctx context.Context, userID uuid.UUID) (*User, error)
}
