package token

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type AuthHandler struct {
	Store  *SQLiteStore
	Issuer *Issuer
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	isLogin := path == "/api/auth/login"
	if r.Method != http.MethodPost {
		if isLogin {
			writeLoginFailure(w)
			return
		}
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var input credentials
	r.Body = http.MaxBytesReader(w, r.Body, 8<<10)
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		if isLogin {
			writeLoginFailure(w)
			return
		}
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	var user User
	var err error
	switch strings.TrimSuffix(r.URL.Path, "/") {
	case "/api/auth/register":
		user, err = h.Store.Register(r.Context(), input.Username, input.Password)
		if errors.Is(err, ErrUsernameTaken) {
			http.Error(w, "username already exists", http.StatusConflict)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case "/api/auth/login":
		user, err = h.Store.Authenticate(r.Context(), input.Username, input.Password)
		if err != nil {
			writeLoginFailure(w)
			return
		}
	default:
		http.NotFound(w, r)
		return
	}
	record, err := h.Issuer.IssueFor(r.Context(), user.ID)
	if err != nil {
		if isLogin {
			writeLoginFailure(w)
			return
		}
		http.Error(w, "token issuance failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	status := http.StatusOK
	if strings.TrimSuffix(r.URL.Path, "/") == "/api/auth/register" {
		status = http.StatusCreated
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(struct {
		Success   bool   `json:"success"`
		User      User   `json:"user"`
		Token     string `json:"token"`
		TokenType string `json:"tokenType"`
		ExpiresAt string `json:"expiresAt"`
	}{Success: true, User: user, Token: record.Value, TokenType: "Bearer", ExpiresAt: record.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")})
}

func writeLoginFailure(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{Success: false, Message: "登录失败"})
}
