package handlers

import (
	"github/oxiginedev/gostarter/config"
	"github/oxiginedev/gostarter/internal/database"
	"github/oxiginedev/gostarter/internal/pkg/jwt"
)

type Handler struct {
	DB     database.Database
	Config *config.Configuration
	JWT    *jwt.JWT
}
