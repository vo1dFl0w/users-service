package http_adaptor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/users-service/internal/app/adapters/http/auth"
	"github.com/vo1dFl0w/users-service/internal/app/adapters/http/middlewares"
	"github.com/vo1dFl0w/users-service/internal/app/adapters/http/user"
	"github.com/vo1dFl0w/users-service/internal/app/adapters/http/utils"
)

func (h *Handler) Routes() http.Handler {
	authHandler := auth.NewAuthHandler(h.AuthService, h.Logger)

	userHandler := user.NewUserHandler(h.UserService, h.Logger)

	h.Root = middlewares.LoggerMiddleware(h.Logger)(h.Router)

	h.Router.HandleFunc("/register", authHandler.Register())
	h.Router.HandleFunc("/login", authHandler.Login())

	authorized := http.NewServeMux()
	authorized.Handle("/users/", middlewares.AuthMiddleware(h.JWTService)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parts := parseURL(r.URL.Path)

			if len(parts) == 2 && parts[0] == "users" && parts[1] == "leaderboard" {
				userHandler.Leaderboard()(w, r)
				return
			}

			if len(parts) == 3 && parts[0] == "users" {
				userID, err := parseUUID(parts[1])
				if err != nil {
					utils.ErrorFunc(w, r, http.StatusUnprocessableEntity, err)
					return
				}

				switch parts[2] {
				case "status":
					userHandler.GetUserStatus(userID)(w, r)
					return
				case "referrer":
					userHandler.Refferer(userID)(w, r)
					return
				default:
					utils.ErrorFunc(w, r, http.StatusBadRequest, fmt.Errorf("unknown endpoint"))
					return
				}
			}

			if len(parts) == 4 && parts[0] == "users" && parts[2] == "task" && parts[3] == "complete" {
				userID, err := parseUUID(parts[1])
				if err != nil {
					utils.ErrorFunc(w, r, http.StatusUnprocessableEntity, err)
					return
				}

				userHandler.CompleteTask(userID)(w, r)
				return
			}
		}),
	))
	h.Router.Handle("/users/", authorized)

	return h.Router
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Root.ServeHTTP(w, r)
}

func parseURL(url string) []string {
	path := strings.Trim(url, "/")
	if path == "" {
		return nil
	}
	return strings.FieldsFunc(path, func(r rune) bool {
		return r == '/'
	})
}

func parseUUID(userID string) (uuid.UUID, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user_id")
	}

	return id, nil
}
