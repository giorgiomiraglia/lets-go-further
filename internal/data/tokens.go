package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"greenlight/internal/validator"
	"time"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type TokenModel struct {
	DB *sql.DB
}

type Token struct {
	PlainText string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {

	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	// Fills randomBytes with random bytes using the system's
	// criptographically secure random number generator (CSRNG)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// Encodes de random bytes to a base-32-encoded string, trimming the '=' from padding
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Generates the hash
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}

func (m *TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {

	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, nil
}

func (m *TokenModel) Insert(t *Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)
		;
	`
	args := []any{t.Hash, t.UserID, t.Expiry, t.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(
		ctx,
		query,
		args...,
	)

	return err
}

func (m *TokenModel) DeleteForAllUsers(scope string, userID int64) error {
	query := `
		DELETE FROM tokens
		WHERE scope = $1 AND user_id = $2
		;
	`
	args := []any{scope, userID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(
		ctx,
		query,
		args...,
	)

	return err
}

func ValidateTokenPlainText(v *validator.Validator, token string) {
	v.Check(token != "", "token", "must be provided")
	v.Check(len(token) == 26, "token", "must be 26 bytes long")
}
