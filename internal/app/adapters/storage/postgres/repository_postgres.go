package postgres

import (
	"database/sql"

	"github.com/vo1dFl0w/users-service/internal/app/domain/auth_domain"
	"github.com/vo1dFl0w/users-service/internal/app/domain/user_domain"
)

type Storage struct {
	DB *sql.DB
	authRepository auth_domain.AuthRepository
	userRepository user_domain.UserRepository
}

func New(db *sql.DB) *Storage {
	return &Storage{
		DB: db,
	}
}

func (s *Storage) Auth() auth_domain.AuthRepository {
	if s.authRepository != nil {
		return s.authRepository
	}

	s.authRepository = &Auth{
		DB: s.DB,
	}

	return s.authRepository
}

func (s *Storage) User() user_domain.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &User{
		DB: s.DB,
	}

	return s.userRepository
}
