package auth

import "net/http"

type Auth interface {
	Register() http.HandlerFunc
	Login() http.HandlerFunc
	DeleteUser() http.HandlerFunc
	RefreshToken() http.HandlerFunc
}
