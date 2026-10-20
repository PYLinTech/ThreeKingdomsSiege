package token

import (
	"context"
	"database/sql"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// SQLiteStore persists token records across server restarts. Only a SHA-256
// digest is stored; raw bearer credentials remain in the issuing client/cache.
type SQLiteStore struct{ db *sql.DB }

func OpenSQLite(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	const schema = `CREATE TABLE IF NOT EXISTS auth_tokens (
		token_hash TEXT PRIMARY KEY,
		user_id TEXT NOT NULL DEFAULT '',
		issued_at TEXT NOT NULL,
		expires_at TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL COLLATE NOCASE UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TEXT NOT NULL
	);`
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}
	// Migrate databases created by the token-only first phase.
	if _, err := db.Exec(`ALTER TABLE auth_tokens ADD COLUMN user_id TEXT NOT NULL DEFAULT ''`); err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate column") {
		db.Close()
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Save(ctx context.Context, record Record) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO auth_tokens(token_hash, user_id, issued_at, expires_at) VALUES(?, ?, ?, ?)`,
		digest(record.Value), record.UserID, record.IssuedAt.UTC().Format(time.RFC3339Nano), record.ExpiresAt.UTC().Format(time.RFC3339Nano),
	)
	return err
}

func (s *SQLiteStore) Close() error { return s.db.Close() }
