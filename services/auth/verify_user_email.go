package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/oxiginedev/gostarter/internal/database"
	"github.com/oxiginedev/gostarter/internal/models"
	"github.com/oxiginedev/gostarter/services"
)

type VerifyUserEmailService struct {
	DB    *database.Connection
	Token string
}

func (vs *VerifyUserEmailService) Run(ctx context.Context) error {
	code := http.StatusBadRequest

	user, err := models.FindUserByVerificationToken(vs.DB, vs.Token)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			return services.NewServiceError(code, "invalid verification token", nil)
		}

		return err
	}

	if user.HasVerifiedEmail() {
		return nil
	}

	if time.Now().After(user.EmailVerificationExpiresAt) {
		return services.NewServiceError(code, "verification token expired", nil)
	}

	return user.MarkEmailAsVerified(vs.DB)
}
