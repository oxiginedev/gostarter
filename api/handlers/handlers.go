package handlers

import (
	"github/oxiginedev/gostarter/config"
	"github/oxiginedev/gostarter/internal/database"
	"github/oxiginedev/gostarter/internal/pkg/jwt"
)

type Handler struct {
	Config *config.Configuration
	DB     *database.Connection
	JWT    *jwt.JWT
}
