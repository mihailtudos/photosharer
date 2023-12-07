package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/mihailtudos/photosharer/rand"
	"strings"
	"time"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when a PasswordReset is being created.
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each session token. If this value is not set or is less then the
	// MinBytesPerToken const it will be ignored
	BytesPerToken int
	// Duration is the amount of time that a PasswordReset is valid for.
	// Defaults to DefaultResetDuration
	Duration time.Duration
}

func (ps *PasswordResetService) Create(email string) (*PasswordReset, error) {
	email = strings.ToLower(email)
	var userId int
	row := ps.DB.QueryRow(`
		SELECT id FROM users WHERE email = $1;
	`, email)

	err := row.Scan(&userId)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	bytesPerToken := ps.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	token, err := rand.String(bytesPerToken)

	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	duration := ps.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}

	pwReset := PasswordReset{
		UserID:    userId,
		Token:     token,
		TokenHash: ps.hash(token),
		ExpiresAt: time.Now().Add(duration),
	}

	row = ps.DB.QueryRow(`
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3) ON CONFLICT (user_id) DO 
		    UPDATE 
		SET token_hash = $2, expires_at = $3 
		RETURNING id;`, pwReset.UserID, pwReset.TokenHash, pwReset.ExpiresAt)

	err = row.Scan(&pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &pwReset, nil
}

func (ps *PasswordResetService) Consume(token string) (*User, error) {
	tokenHash := ps.hash(token)
	var user User
	var pwReset PasswordReset

	row := ps.DB.QueryRow(`
		SELECT password_resets.id, 
		       password_resets.expires_at,
		       users.id,
		       users.email,
		       users.password_hash
		FROM password_resets
			JOIN users on users.id = password_resets.user_id
		WHERE password_resets.token_hash = $1;
	`, tokenHash)

	err := row.Scan(
		&pwReset.ID, &pwReset.ExpiresAt,
		&user.ID, &user.Email, &user.PasswordHash)

	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	if time.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("token expired: %v", token)
	}

	err = ps.delete(pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	return &user, nil
}

func (ps *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (ps *PasswordResetService) delete(id int) error {
	_, err := ps.DB.Exec(`
		DELETE FROM password_resets
		WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
