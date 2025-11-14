package user

import (
	"net/http"

	"github.com/google/uuid"
)

type TaskHandler interface {
	GetUserStatus(userID uuid.UUID) http.HandlerFunc
	Leaderboard() http.HandlerFunc
	CompleteTask(userID uuid.UUID) http.HandlerFunc
	Refferer(userID uuid.UUID) http.HandlerFunc
}
