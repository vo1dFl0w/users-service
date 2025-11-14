package user_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/users-service/internal/app/adapters/storage"
	"github.com/vo1dFl0w/users-service/internal/app/domain/user_domain"
)

type Service interface {
	UserStatus(ctx context.Context, userID uuid.UUID) (*user_domain.User, error)
	Leaderboard(ctx context.Context) (map[int]map[string]interface{}, error)
	CompleteUserTask(ctx context.Context, userID uuid.UUID, task string) error
	Referrer(ctx context.Context, userID uuid.UUID, referrerID uuid.UUID, task string) error
}

type service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage, ) Service {
	return &service{
		storage: storage,
	}
}

func (s *service) UserStatus(ctx context.Context, userID uuid.UUID) (*user_domain.User, error) {
	u := &user_domain.User{UserID: userID}

	if err := u.ValidateUUID(); err != nil {
		return nil, err
	}

	return s.storage.User().UserStatus(ctx, userID)
}

func (s *service) Leaderboard(ctx context.Context) (map[int]map[string]interface{}, error) {
	return s.storage.User().Leaderboard(ctx)
}

func (s *service) CompleteUserTask(ctx context.Context, userID uuid.UUID, task string) error {
	u := &user_domain.User{UserID: userID}

	if err := u.ValidateUUID(); err != nil {
		return err
	}

	if task == "" {
		return fmt.Errorf("empty task")
	}

	return s.storage.User().CompleteUserTask(ctx, userID, task)
}

func (s *service) Referrer(ctx context.Context, userID uuid.UUID, referrerID uuid.UUID, task string) error {
	u := user_domain.User{UserID: userID}

	if err := u.ValidateUUID(); err != nil {
		return err
	}

	if userID == referrerID {
		return fmt.Errorf("referrer_id cannot be the same as user_id")
	}

	if task == "" {
		return fmt.Errorf("empty task")
	}

	return s.storage.User().Referrer(ctx, userID, referrerID, task)
}
