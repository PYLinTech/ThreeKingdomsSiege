package token

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrNotFound = errors.New("token not found")

// Record stores a server-issued bearer token and its lifecycle timestamps.
type Record struct {
	Value     string    `json:"token"`
	UserID    string    `json:"userId"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// Store is the durable token storage abstraction (SQLite in production).
type Store interface {
	Save(context.Context, Record) error
}

// Cache is the in-process lookup cache used to validate stateful tokens.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]Record
}

func NewCache() *Cache { return &Cache{entries: make(map[string]Record)} }

func (c *Cache) Put(record Record) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[record.Value] = record
}

// Get compares token bytes in constant time after an indexed lookup.
func (c *Cache) Get(value string, now time.Time) (Record, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	record, ok := c.entries[value]
	if !ok || subtle.ConstantTimeCompare([]byte(record.Value), []byte(value)) != 1 || !now.Before(record.ExpiresAt) {
		return Record{}, false
	}
	return record, true
}

// Issuer creates cryptographically random, opaque, stateful bearer tokens.
type Issuer struct {
	cache  *Cache
	store  Store
	ttl    time.Duration
	now    func() time.Time
	random func([]byte) (int, error)
}

func NewIssuer(cache *Cache, store Store, ttl time.Duration) (*Issuer, error) {
	if cache == nil {
		return nil, errors.New("token cache is required")
	}
	if ttl <= 0 {
		return nil, errors.New("token ttl must be positive")
	}
	return &Issuer{cache: cache, store: store, ttl: ttl, now: time.Now, random: rand.Read}, nil
}

func (i *Issuer) Issue(ctx context.Context) (Record, error) { return i.IssueFor(ctx, "") }

func (i *Issuer) IssueFor(ctx context.Context, userID string) (Record, error) {
	buf := make([]byte, 32) // 256 bits of entropy
	if _, err := i.random(buf); err != nil {
		return Record{}, fmt.Errorf("generate token: %w", err)
	}
	now := i.now().UTC()
	record := Record{
		Value:     base64.RawURLEncoding.EncodeToString(buf),
		UserID:    userID,
		IssuedAt:  now,
		ExpiresAt: now.Add(i.ttl),
	}
	// Persist first so a database failure never returns a token that cannot be recovered.
	if i.store != nil {
		if err := i.store.Save(ctx, record); err != nil {
			return Record{}, fmt.Errorf("persist token: %w", err)
		}
	}
	i.cache.Put(record)
	return record, nil
}
