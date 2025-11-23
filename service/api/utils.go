package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	// "git.sapienzaapps.it/.../reqcontext" <--- QUESTO DEVE ESSERE RIMOSSO O COMMENTATO

	"github.com/google/uuid"
)

func (rt *_router) sendJSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

func (rt *_router) sendJSONError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errResp := Error{
		Code:    code,
		Message: message,
	}
	_ = json.NewEncoder(w).Encode(errResp)
}

func (rt *_router) extractUser(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("auth token missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("auth token missing")
	}

	token := parts[1]

	if _, err := uuid.Parse(token); err != nil {
		return "", errors.New("invalid auth token format")
	}

	return token, nil
}
