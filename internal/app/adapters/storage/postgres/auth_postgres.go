package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/users-service/internal/app/domain/auth_domain"
)

var (
	subscrubeTelegramTask  = "subscribe to 'telegram' channel/group"
	subscrubeInstagramTask = "subscribe to 'instagram' account"
	subscrubeVKTask        = "subscribe to 'vkontakte' group"
	reward                 = 150
)

type Auth struct {
	DB *sql.DB
}

func (a *Auth) CreateUser(ctx context.Context, email string, encryptedPassword string) error {
	tx, err := a.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("failed to start create user transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var userID uuid.UUID
	if err := tx.QueryRowContext(ctx,
		"INSERT INTO users (email, encrypted_password, created_at) VALUES ($1, $2, NOW()) RETURNING user_id",
		email, encryptedPassword,
	).Scan(&userID); err != nil {
		return fmt.Errorf("failed to create new user: %w", err)
	}

	task1, err := tx.ExecContext(ctx,
		"INSERT INTO users_tasks (user_id, task, reward) VALUES ($1, $2, $3)",
		userID, subscrubeTelegramTask, reward,
	)

	if err != nil {
		return fmt.Errorf("failed to create new task to user: %w", err)
	}

	r, err := task1.RowsAffected()
	if err == nil {
		if r == 0 {
			return fmt.Errorf("no row added")
		}
	}

	task2, err := tx.ExecContext(ctx,
		"INSERT INTO users_tasks (user_id, task, reward) VALUES ($1, $2, $3)",
		userID, subscrubeInstagramTask, reward,
	)

	if err != nil {
		return fmt.Errorf("failed to create new task to user: %w", err)
	}

	r, err = task2.RowsAffected()
	if err == nil {
		if r == 0 {
			return fmt.Errorf("no row added")
		}
	}

	task3, err := tx.ExecContext(ctx,
		"INSERT INTO users_tasks (user_id, task, reward) VALUES ($1, $2, $3)",
		userID, subscrubeVKTask, reward,
	)

	if err != nil {
		return fmt.Errorf("failed to create new task to user: %w", err)
	}

	r, err = task3.RowsAffected()
	if err == nil {
		if r == 0 {
			return fmt.Errorf("no row added")
		}
	}

	row, err := tx.ExecContext(ctx,
		"INSERT INTO users_scoreboard (user_id) VALUES ($1)",
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to create new user in users_scoreboard table: %w", err)
	}

	r, err = row.RowsAffected()
	if err == nil {
		if r == 0 {
			return fmt.Errorf("no row added")
		}
	}

	return nil
}

func (a *Auth) GetUser(ctx context.Context, email string) (*auth_domain.User, error) {
	u := &auth_domain.User{}

	err := a.DB.QueryRowContext(ctx,
		"SELECT user_id, encrypted_password FROM users WHERE email = $1",
		email,
	).Scan(&u.UserID, &u.EncryptedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		} else {
			return nil, err
		}
	}

	return u, nil
}

// TODO func (a *Auth) UpdateUser(ctx context.Context, user *auth.User) error {return nil}

// TODO func (a *Auth) DeleteUser(ctx context.Context, user *auth.User) error {return nil}

func (a *Auth) SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiry time.Time) error {
	tx, err := a.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("failed to start 'save refresh token' transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx,
		"INSERT INTO users_tokens (user_id, refresh_token, refresh_token_expiry) VALUES($1, $2, $3)",
		userID, token, expiry,
	)
	if err != nil {
		return err
	}

	return nil
}

/* TODO
func (a *Auth) GetRefreshToken(ctx context.Context, userID uuid.UUID) (*auth_domain.User, error) {
	u := &auth_domain.User{}

	if err := a.DB.QueryRowContext(ctx,
		"SELECT user_id, refresh_token, refresh_token_expiry FROM users_tokens WHERE user_id = $1",
		userID,
	).Scan(&u.UserID, &u.RefreshToken, &u.RefreshTokenExpire); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		} else {
			return nil, err
		}
	}

	return u, nil
}
*/
