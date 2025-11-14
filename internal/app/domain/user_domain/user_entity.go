package user_domain

import (
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	UserID   uuid.UUID `json:"user_id"`
	Score    int64     `json:"score"`
	Task     string    `json:"task"`
	Reward   int64     `json:"reward"`
	Complete bool      `json:"complete"`
}

func (u *User) ValidateUUID() error {
	if u.UserID == uuid.Nil {
		return fmt.Errorf("empty user_id")
	}

	_, err := uuid.Parse(u.UserID.String())
	if err != nil {
		return fmt.Errorf("invalid user_id")
	}

	return nil
}
