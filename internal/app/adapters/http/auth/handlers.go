package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/vo1dFl0w/users-service/internal/app/adapters/http/utils"
	"github.com/vo1dFl0w/users-service/internal/app/usecase/auth_usecase"
)

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
)

type AuthHandler struct {
	AuthService auth_usecase.Service
	Logger      *slog.Logger
}

func NewAuthHandler(ac auth_usecase.Service, log *slog.Logger) *AuthHandler {
	return &AuthHandler{
		AuthService: ac,
		Logger:      log,
	}
}

func (h *AuthHandler) Register() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		if r.Method != http.MethodPost {
			utils.ErrorFunc(w, r, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorFunc(w, r, http.StatusBadRequest, err)
			return
		}

		if err := h.AuthService.CreateUser(ctx, req.Email, req.Password); err != nil {
			utils.ErrorFunc(w, r, http.StatusBadRequest, err)
			return
		}

		utils.RespondFunc(w, r, http.StatusCreated, map[string]string{"status": "success"})
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		if r.Method != http.MethodPost {
			utils.ErrorFunc(w, r, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorFunc(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := h.AuthService.GetUser(ctx, req.Email, req.Password)
		if err != nil {
			utils.ErrorFunc(w, r, http.StatusBadRequest, err)
			return
		}

		accessToken, refreshToken, err := h.AuthService.IssueTokens(ctx, u.UserID)
		if err != nil {
			utils.ErrorFunc(w, r, http.StatusBadRequest, err)
			return
		}

		utils.RespondFunc(w, r, http.StatusOK, map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

/* TODO
func (h *AuthHandler) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		if r.Method != http.MethodDelete {
			utils.ErrorFunc(w, r, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			return
		}

	}
}
*/

/* TODO
func (h *AuthHandler) RefreshTocken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		if r.Method != http.MethodPost {
			utils.ErrorFunc(w, r, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			return
		}

	}
}
*/
