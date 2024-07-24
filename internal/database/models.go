package database

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Email           string       `json:"email"`
	EmailVerifiedAt sql.NullTime `json:"email_verified_at"`
	Password        string       `json:"password"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type Password struct {
	PlainText string
	Hash      string
}

func (p *Password) GenerateHash() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.PlainText), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.Hash = string(hash)

	return nil
}

func (p *Password) Matches() (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(p.Hash), []byte(p.PlainText))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
