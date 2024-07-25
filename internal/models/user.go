package models

import (
	"github/oxiginedev/gostarter/internal/database"
	"github/oxiginedev/gostarter/util"
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID              string     `json:"id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Email           string     `json:"email"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	Password        string     `json:"-"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (u *User) HasVerifiedEmail() bool {
	return u.EmailVerifiedAt != nil
}

func (u *User) MarkEmailAsVerified(tx *database.Connection) bool {
	return u.EmailVerifiedAt != nil
}

func CreateUser(tx *database.Connection, user *User) error {
	passwordHash, err := util.GenerateHashFromPassword(user.Password)
	if err != nil {
		return err
	}

	user.ID = uuid.Must(uuid.NewV7()).String()
	user.Password = passwordHash

	return nil
}

// FindUserByID finds a user with matching id
func FindUserByID(tx *database.Connection, id string) (*User, error) {
	user := User{}

	if err := tx.Find(&user, "id = ?", id); err != nil {
		return nil, err
	}

	return &user, nil
}

// FindUserByEmailAddress finds a user with matching email address
func FindUserByEmailAddress(tx *database.Connection, email string) (*User, error) {
	user := User{}

	if err := tx.Find(&user, "email = ?", email); err != nil {
		return nil, err
	}

	return &user, nil
}
