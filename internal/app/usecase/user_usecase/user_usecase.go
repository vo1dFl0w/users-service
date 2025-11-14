package user_usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/users-service/internal/app/domain/user_domain"
)

type Service interface {
	UserStatus(ctx context.Context, userID uuid.UUID) (*user_domain.User, error)
	Leaderboard(ctx context.Context) (map[int]map[string]interface{}, error)
	CompleteUserTask(ctx context.Context, userID uuid.UUID, task string) error
	Referrer(ctx context.Context, userID uuid.UUID, referrerID uuid.UUID, task string) error
}

type service struct {
	repository user_domain.UserRepository
}

func NewService(repository user_domain.UserRepository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) UserStatus(ctx context.Context, userID uuid.UUID) (*user_domain.User, error) {
	u := &user_domain.User{UserID: userID}

	if err := u.ValidateUUID(); err != nil {
		return nil, err
	}

	return s.repository.UserStatus(ctx, userID)
}

func (s *service) Leaderboard(ctx context.Context) (map[int]map[string]interface{}, error) {
	return s.repository.Leaderboard(ctx)
}

func (s *service) CompleteUserTask(ctx context.Context, userID uuid.UUID, task string) error {
	u := &user_domain.User{UserID: userID}

	if err := u.ValidateUUID(); err != nil {
		return err
	}

	if task == "" {
		return fmt.Errorf("empty task")
	}

	return s.repository.CompleteUserTask(ctx, userID, task)
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

	return s.repository.Referrer(ctx, userID, referrerID, task)
}
