package postgres

import (
	"fmt"
	"github/oxiginedev/gostarter/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	*gorm.DB
}

func NewPostgresRepository(cfg *config.DatabaseConfiguration) (*PostgresRepository, error) {
	db, err := gorm.Open(postgres.Open(cfg.BuildDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database - %v", err)
	}

	return &PostgresRepository{db}, nil
}

func (p *PostgresRepository) GetDB() *gorm.DB {
	return p.DB
}

func (p *PostgresRepository) CloseDB() error {
	return nil
}
