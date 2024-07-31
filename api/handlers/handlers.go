package handlers

import (
	"github.com/oxiginedev/gostarter/config"
	"github.com/oxiginedev/gostarter/internal/database"
	"github.com/oxiginedev/gostarter/internal/pkg/jwt"
)

type Handler struct {
	Config *config.Configuration
	DB     *database.Connection
	JWT    *jwt.JWT
}
