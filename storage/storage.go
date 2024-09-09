package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"

	"github.com/v-starostin/goph-keeper/internal/model"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) AddUser(ctx context.Context, username, password string) error {
	query := "INSERT INTO users(username, password) VALUES ($1, $2)"
	_, err := s.db.ExecContext(ctx, query, username, hash(password))
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUser(ctx context.Context, username, password string) (*model.User, error) {
	query := "SELECT * FROM users WHERE username = $1 AND password = $2"
	var u model.User
	err := s.db.QueryRowContext(ctx, query, username, hash(password)).Scan(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Storage) GetTokenByUserID(ctx context.Context, userID int32) (string, error) {
	query := "SELECT token FROM tokens WHERE user_id = $1"
	var token string
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Storage) SaveToken(ctx context.Context, userID int32, token string) error {
	query := "INSERT INTO tokens (token, user_id) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET token=EXCLUDED.token"
	_, err := s.db.ExecContext(ctx, query, token, userID)
	if err != nil {
		return err
	}
	return nil
}

func hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
