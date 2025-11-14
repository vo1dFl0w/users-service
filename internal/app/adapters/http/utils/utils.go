package utils

import (
	"encoding/json"
	"net/http"
)

func ErrorFunc(w http.ResponseWriter, r *http.Request, code int, err error) {
	RespondFunc(w, r, code, map[string]string{"error": err.Error()})
}

func RespondFunc(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}