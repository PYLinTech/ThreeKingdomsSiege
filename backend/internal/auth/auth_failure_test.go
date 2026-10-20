package token

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestLoginFailuresUseOneGenericResponse(t *testing.T) {
	store, err := OpenSQLite(t.TempDir() + "/login-failure.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	issuer, err := NewIssuer(NewCache(), store, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	h := AuthHandler{Store: store, Issuer: issuer}

	responses := make([]string, 0, 2)
	for _, body := range []string{`{"username":"missing","password":"wrong-password"}`, `{not-json`} {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
		res := httptest.NewRecorder()
		h.ServeHTTP(res, req)
		if res.Code != http.StatusUnauthorized {
			t.Fatalf("status=%d, want 401", res.Code)
		}
		responses = append(responses, strings.TrimSpace(res.Body.String()))
	}
	if responses[0] != responses[1] {
		t.Fatalf("login errors differ: %q vs %q", responses[0], responses[1])
	}
	if !strings.Contains(responses[0], `"message":"登录失败"`) {
		t.Fatalf("unexpected failure response: %s", responses[0])
	}
}
