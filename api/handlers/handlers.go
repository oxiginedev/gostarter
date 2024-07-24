package handlers

import (
	"github/oxiginedev/gostarter/config"
	"github/oxiginedev/gostarter/internal/database"
)

type Handler struct {
	DB     database.Database
	Config *config.Configuration
}
