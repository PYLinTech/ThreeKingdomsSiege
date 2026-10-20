package token

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type memoryStore struct{ saved []Record }

func (s *memoryStore) Save(_ context.Context, r Record) error {
	s.saved = append(s.saved, r)
	return nil
}

func TestIssueCreatesDistinctCachedToken(t *testing.T) {
	cache := NewCache()
	store := &memoryStore{}
	issuer, err := NewIssuer(cache, store, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	first, err := issuer.Issue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	second, err := issuer.Issue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if first.Value == second.Value {
		t.Fatal("tokens must be unique")
	}
	if len(first.Value) != 43 {
		t.Fatalf("expected 256-bit URL-safe token, got %d chars", len(first.Value))
	}
	if _, ok := cache.Get(first.Value, time.Now()); !ok {
		t.Fatal("issued token missing from cache")
	}
	if len(store.saved) != 2 {
		t.Fatalf("expected 2 persisted tokens, got %d", len(store.saved))
	}
}

func TestHandlerRequiresPost(t *testing.T) {
	issuer, err := NewIssuer(NewCache(), nil, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/auth/token", nil)
	res := httptest.NewRecorder()
	(Handler{Issuer: issuer}).ServeHTTP(res, req)
	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got status %d", res.Code)
	}
}
