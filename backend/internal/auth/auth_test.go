package token

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRegisterAuthenticateAndDuplicate(t *testing.T) {
	store, err := OpenSQLite(t.TempDir() + "/auth.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	ctx := context.Background()
	user, err := store.Register(ctx, "demo_user", "correct-horse")
	if err != nil {
		t.Fatal(err)
	}
	if user.ID == "" || user.Username != "demo_user" {
		t.Fatalf("unexpected user: %+v", user)
	}
	loggedIn, err := store.Authenticate(ctx, "DEMO_USER", "correct-horse")
	if err != nil || loggedIn.ID != user.ID {
		t.Fatalf("login failed: user=%+v err=%v", loggedIn, err)
	}
	if _, err = store.Authenticate(ctx, "demo_user", "wrong-password"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
	if _, err = store.Register(ctx, "DEMO_USER", "another-pass"); !errors.Is(err, ErrUsernameTaken) {
		t.Fatalf("expected duplicate username, got %v", err)
	}
}

func TestAuthHandlerRegisterAndLoginIssueUserToken(t *testing.T) {
	store, err := OpenSQLite(t.TempDir() + "/handler.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	cache := NewCache()
	issuer, err := NewIssuer(cache, store, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	h := AuthHandler{Store: store, Issuer: issuer}
	register := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"username":"tester","password":"good-password"}`))
	register.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	h.ServeHTTP(response, register)
	if response.Code != http.StatusCreated {
		t.Fatalf("register status=%d body=%s", response.Code, response.Body.String())
	}
	var issued struct {
		Token string `json:"token"`
		User  User   `json:"user"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &issued); err != nil {
		t.Fatalf("decode register response: %v", err)
	}
	cached, ok := cache.Get(issued.Token, time.Now())
	if !ok || cached.UserID != issued.User.ID {
		t.Fatalf("register token not linked to account: cached=%+v user=%+v", cached, issued.User)
	}
	login := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"tester","password":"good-password"}`))
	response = httptest.NewRecorder()
	h.ServeHTTP(response, login)
	if response.Code != http.StatusOK {
		t.Fatalf("login status=%d body=%s", response.Code, response.Body.String())
	}
}
