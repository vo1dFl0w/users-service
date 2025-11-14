package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/users-service/internal/app/domain/user_domain"
)

var (
	refReward = 100
	usrReward = 50
)

type User struct {
	DB *sql.DB
}

func (u *User) UserStatus(ctx context.Context, userID uuid.UUID) (*user_domain.User, error) {
	usr := &user_domain.User{}

	if err := u.DB.QueryRowContext(ctx,
		"SELECT score FROM users_scoreboard WHERE user_id = $1",
		userID,
	).Scan(&usr.Score); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		} else {
			return nil, err
		}
	}

	usr.UserID = userID

	return usr, nil
}

func (u *User) Leaderboard(ctx context.Context) (map[int]map[string]interface{}, error) {
	users := make(map[int]map[string]interface{})

	rows, err := u.DB.QueryContext(ctx,
		"SELECT user_id, score FROM users_scoreboard ORDER BY score DESC LIMIT 10",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	defer rows.Close()

	i := 1
	for rows.Next() {
		var usr user_domain.User
		err = rows.Scan(&usr.UserID, &usr.Score)
		if err != nil {
			return nil, err
		}
		users[i] = map[string]interface{}{"user_id": usr.UserID, "score": usr.Score}
		i++
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}

func (u *User) CompleteUserTask(ctx context.Context, userID uuid.UUID, task string) error {
	tx, err := u.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return fmt.Errorf("failed to start 'complete task' transaction: %w", err)
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

	reward := 0
	if err = tx.QueryRowContext(ctx,
		"UPDATE users_tasks SET complete = true WHERE user_id = $1 AND task = $2 RETURNING reward",
		userID, task,
	).Scan(&reward); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("task not found")
		}
		return fmt.Errorf("failed to update users_task: %w", err)
	}

	row, err := tx.ExecContext(ctx,
		"UPDATE users_scoreboard SET score = score + $1 WHERE user_id = $2",
		reward, userID,
	)
	r, err := row.RowsAffected()
	if err == nil {
		if r == 0 {
			return fmt.Errorf("no rows updated")
		}
	}

	return nil
}

func (u *User) Referrer(ctx context.Context, userID uuid.UUID, referrerID uuid.UUID, task string) error {
	tx, err := u.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return fmt.Errorf("failed to start 'refferer' transaction: %w", err)
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

	var ref sql.NullString

	if err = tx.QueryRowContext(ctx,
		"SELECT referrer_id FROM users_tasks WHERE user_id = $1 AND task = $2",
		userID, task,
	).Scan(&ref); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to select referrer_id by user_id: %w", err)
	}

	if ref.Valid && ref.String != "" {
		return fmt.Errorf("cannot use refer")
	}

	row, err := tx.ExecContext(ctx,
		"UPDATE users_tasks SET referrer_id = $1 WHERE user_id = $2 AND task = $3",
		referrerID, userID, task,
	)
	if err != nil {
		return fmt.Errorf("failed to update referrer_id by user_id")
	}

	r, err := row.RowsAffected()
	if err == nil {
		if r == 0 {
			return fmt.Errorf("no row updated")
		}
	}

	row, err = tx.ExecContext(ctx,
		"UPDATE users_scoreboard SET score = score + $1 WHERE user_id = $2",
		refReward, referrerID,
	)
	if err != nil {
		return fmt.Errorf("failed to update users_scoreboard: %w", err)
	}

	r, err = row.RowsAffected()
	if err == nil {
		if r == 0 {
			return fmt.Errorf("no row updated")
		}
	}

	row, err = tx.ExecContext(ctx,
		"UPDATE users_scoreboard SET score = score + $1 WHERE user_id = $2",
		usrReward, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update users_scoreboard: %w", err)
	}

	r, err = row.RowsAffected()
	if err == nil {
		if r == 0 {
			return fmt.Errorf("no row updated")
		}
	}

	return nil
}
