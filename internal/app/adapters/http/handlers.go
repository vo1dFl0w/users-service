package http_adaptor

import (
	"log/slog"
	"net/http"

	"github.com/vo1dFl0w/users-service/internal/app/usecase/auth_usecase"
	"github.com/vo1dFl0w/users-service/internal/app/usecase/jwt_usecase"
	"github.com/vo1dFl0w/users-service/internal/app/usecase/user_usecase"
)

type Handler struct {
	Router      *http.ServeMux
	Root  http.Handler
	Logger      *slog.Logger
	JWTService  jwt.Service
	AuthService auth_usecase.Service
	UserService user_usecase.Service
}

func NewHandler(log *slog.Logger, token jwt.Service, auth auth_usecase.Service, user user_usecase.Service) *Handler {
	h := &Handler{
		Router:      http.NewServeMux(),
		Logger:      log,
		JWTService:  token,
		AuthService: auth,
		UserService: user,
	}

	h.Routes()

	return h
}
