package storage

import (
	"github.com/vo1dFl0w/users-service/internal/app/domain/auth_domain"
	"github.com/vo1dFl0w/users-service/internal/app/domain/user_domain"
)

type Storage interface {
	Auth() auth_domain.AuthRepository
	User() user_domain.UserRepository
}