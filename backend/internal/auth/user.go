package token

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameTaken      = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
	usernamePattern       = regexp.MustCompile(`^[A-Za-z0-9_]{3,32}$`)
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

// Register validates credentials and stores only a bcrypt password hash.
func (s *SQLiteStore) Register(ctx context.Context, username, password string) (User, error) {
	username = strings.TrimSpace(username)
	if !usernamePattern.MatchString(username) {
		return User{}, errors.New("username must be 3-32 characters using letters, numbers, or underscore")
	}
	if len(password) < 8 || len(password) > 72 {
		return User{}, errors.New("password must be 8-72 bytes")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}
	id, err := randomID()
	if err != nil {
		return User{}, err
	}
	created := time.Now().UTC()
	_, err = s.db.ExecContext(ctx, `INSERT INTO users(id, username, password_hash, created_at) VALUES(?, ?, ?, ?)`, id, username, string(hash), created.Format(time.RFC3339Nano))
	if err != nil {
		var sqliteErr interface{ Code() int }
		if errors.As(err, &sqliteErr) && sqliteErr.Code()&0xff == 19 {
			return User{}, ErrUsernameTaken
		}
		return User{}, fmt.Errorf("create user: %w", err)
	}
	return User{ID: id, Username: username, CreatedAt: created}, nil
}

func (s *SQLiteStore) Authenticate(ctx context.Context, username, password string) (User, error) {
	var user User
	var hash, created string
	err := s.db.QueryRowContext(ctx, `SELECT id, username, password_hash, created_at FROM users WHERE username = ?`, strings.TrimSpace(username)).Scan(&user.ID, &user.Username, &hash, &created)
	if err != nil {
		return User{}, ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return User{}, ErrInvalidCredentials
	}
	user.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	return user, nil
}

func randomID() (string, error) {
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return "", fmt.Errorf("generate user id: %w", err)
	}
	return hex.EncodeToString(id[:]), nil
}
