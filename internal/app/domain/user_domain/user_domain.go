package user_domain

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	UserStatus(ctx context.Context, userID uuid.UUID) (*User, error)
	Leaderboard(ctx context.Context) (map[int]map[string]interface{}, error)
	CompleteUserTask(ctx context.Context, userID uuid.UUID, task string) error
	Referrer(ctx context.Context, userID uuid.UUID, referrerID uuid.UUID, task string) error
}
