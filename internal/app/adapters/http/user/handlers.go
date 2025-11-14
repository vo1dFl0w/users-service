package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/users-service/internal/app/adapters/http/middlewares"
	"github.com/vo1dFl0w/users-service/internal/app/adapters/http/utils"
	"github.com/vo1dFl0w/users-service/internal/app/usecase/user_usecase"
)

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
)

type UserHandler struct {
	UserService user_usecase.Service
	Logger      *slog.Logger
}

func NewUserHandler(us user_usecase.Service, log *slog.Logger) *UserHandler {
	return &UserHandler{
		UserService: us,
		Logger:      log,
	}
}

func (h *UserHandler) GetUserStatus(userID uuid.UUID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
		defer cancel()

		if r.Method != http.MethodGet {
			utils.ErrorFunc(w, r, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			return
		}

		authUser, ok := getUserID(ctx)
		if !ok {
			utils.ErrorFunc(w, r, http.StatusUnauthorized, fmt.Errorf("access denied"))
			return
		}

		if err := compareUserID(authUser, userID); err != nil {
			utils.ErrorFunc(w, r, http.StatusUnauthorized, err)
			return
		}

		u, err := h.UserService.UserStatus(ctx, userID)
		if err != nil {
			utils.ErrorFunc(w, r, http.StatusBadRequest, err)
			return
		}

		utils.RespondFunc(w, r, http.StatusOK, map[string]interface{}{
			"status":  "success",
			"user_id": u.UserID,
			"score":   u.Score,
		})
	}
}

func (h *UserHandler) Leaderboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
		defer cancel()

		if r.Method != http.MethodGet {
			utils.ErrorFunc(w, r, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			return
		}

		_, ok := getUserID(ctx)
		if !ok {
			utils.ErrorFunc(w, r, http.StatusUnauthorized, fmt.Errorf("access denied"))
			return
		}

		res, err := h.UserService.Leaderboard(ctx)
		if err != nil {
			utils.ErrorFunc(w, r, http.StatusBadRequest, err)
			return
		}

		utils.RespondFunc(w, r, http.StatusOK, res)
	}
}

func (h *UserHandler) CompleteTask(userID uuid.UUID) http.HandlerFunc {
	type request struct {
		Task string `json:"task"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
		defer cancel()

		if r.Method != http.MethodPatch {
			utils.ErrorFunc(w, r, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			return
		}

		authUser, ok := getUserID(ctx)
		if !ok {
			utils.ErrorFunc(w, r, http.StatusUnauthorized, fmt.Errorf("access denied"))
			return
		}

		if err := compareUserID(authUser, userID); err != nil {
			utils.ErrorFunc(w, r, http.StatusUnauthorized, err)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorFunc(w, r, http.StatusBadRequest, err)
			return
		}

		if err := h.UserService.CompleteUserTask(ctx, userID, req.Task); err != nil {
			utils.ErrorFunc(w, r, http.StatusBadRequest, err)
			return
		}

		utils.RespondFunc(w, r, http.StatusOK, map[string]string{"status": "success"})
	}
}

func (h *UserHandler) Refferer(userID uuid.UUID) http.HandlerFunc {
	type request struct {
		Task       string    `json:"task"`
		ReferrerID uuid.UUID `json:"referrer_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
		defer cancel()

		if r.Method != http.MethodPatch {
			utils.ErrorFunc(w, r, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			return
		}

		authUser, ok := getUserID(ctx)
		if !ok {
			utils.ErrorFunc(w, r, http.StatusUnauthorized, fmt.Errorf("access denied"))
			return
		}

		if err := compareUserID(authUser, userID); err != nil {
			utils.ErrorFunc(w, r, http.StatusUnauthorized, err)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorFunc(w, r, http.StatusBadRequest, err)
			return
		}

		if err := h.UserService.Referrer(ctx, userID, req.ReferrerID, req.Task); err != nil {
			utils.ErrorFunc(w, r, http.StatusBadRequest, err)
			return
		}

		utils.RespondFunc(w, r, http.StatusOK, map[string]string{"status": "success"})
	}
}

func getUserID(ctx context.Context) (uuid.UUID, bool) {
	v := ctx.Value(middlewares.CtxKeyUser)
	s, ok := v.(uuid.UUID)
	return s, ok
}

func compareUserID(idStr uuid.UUID, userID uuid.UUID) error {
	if idStr != userID {
		return fmt.Errorf("access denied")
	}

	return nil
}
