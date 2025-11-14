package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vo1dFl0w/users-service/internal/app/adapters/http/utils"
	jwt_usecase "github.com/vo1dFl0w/users-service/internal/app/usecase/jwt_usecase"
)

func AuthMiddleware(jwtServ jwt_usecase.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.ErrorFunc(w, r, http.StatusUnauthorized, fmt.Errorf("missing authorization"))
				return
			}

			parts := strings.Fields(authHeader)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				utils.ErrorFunc(w, r, http.StatusUnauthorized, fmt.Errorf("invalid authorization header"))
				return
			}
			token := parts[1]

			claims, err := jwtServ.ValidateAccessToken(ctx, token)
			if err != nil {
				utils.ErrorFunc(w, r, http.StatusUnauthorized, fmt.Errorf("invalid or expired access token"))
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, CtxKeyUser, claims.UserID)))
		})
	}
}
