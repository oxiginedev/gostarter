package postgres

import (
	"context"
	"errors"
	"github/oxiginedev/gostarter/internal/database"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (gu *userRepository) CreateUser(ctx context.Context, user *database.User) error {
	result := gu.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return database.ErrRecordNotCreated
	}

	return nil
}

func (gu *userRepository) FindUserByEmail(ctx context.Context, email string) (*database.User, error) {
	var user database.User

	tx := gu.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, database.ErrRecordNotFound
		}
		return nil, tx.Error
	}

	return &user, nil
}
