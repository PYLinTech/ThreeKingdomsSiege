package token

import (
	"encoding/json"
	"net/http"
)

// Handler implements POST /api/auth/token. Each call issues a distinct token.
type Handler struct{ Issuer *Issuer }

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/auth/token" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	record, err := h.Issuer.Issue(r.Context())
	if err != nil {
		http.Error(w, "token issuance failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(struct {
		Token     string `json:"token"`
		TokenType string `json:"tokenType"`
		ExpiresAt string `json:"expiresAt"`
	}{Token: record.Value, TokenType: "Bearer", ExpiresAt: record.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")})
}
