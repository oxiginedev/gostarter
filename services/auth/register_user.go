package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/oxiginedev/gostarter/api/types"
	"github.com/oxiginedev/gostarter/internal/database"
	"github.com/oxiginedev/gostarter/internal/email"
	"github.com/oxiginedev/gostarter/internal/models"
	"github.com/oxiginedev/gostarter/internal/pkg/jwt"
	"github.com/oxiginedev/gostarter/internal/worker"
	"github.com/oxiginedev/gostarter/internal/worker/task"
	"github.com/oxiginedev/gostarter/services"
)

type RegisterUserService struct {
	DB    *database.Connection
	JWT   *jwt.JWT
	Queue *worker.Queue
	Data  *types.RegisterUserRequest
}

func (rs *RegisterUserService) Run(ctx context.Context) (*models.User, *jwt.Token, error) {
	exists, err := models.CheckUserExistsByEmailAddress(rs.DB, rs.Data.Email)
	if err != nil {
		return nil, nil, err
	}

	if exists {
		return nil, nil, services.NewServiceError(http.StatusBadRequest, "account with email exists", nil)
	}

	user := &models.User{
		FirstName:                  rs.Data.FirstName,
		LastName:                   rs.Data.LastName,
		Email:                      rs.Data.Email,
		EmailVerificationToken:     uuid.Must(uuid.NewV4()).String(),
		EmailVerificationExpiresAt: time.Now().Add(time.Minute * 30),
		Password:                   rs.Data.Password,
	}

	err = models.CreateUser(rs.DB, user)
	if err != nil {
		return nil, nil, err
	}

	token, err := rs.JWT.GenerateAccessToken(user)
	if err != nil {
		return nil, nil, err
	}

	err = rs.sendVerificationEmail(user)
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}

func (rs *RegisterUserService) sendVerificationEmail(user *models.User) error {
	em := &email.Message{
		Email:    user.Email,
		Subject:  "GoStarter: Verify your Email",
		Template: email.EmailVerifyTemplate,
		Params: map[string]interface{}{
			"first_name": user.FirstName,
			"email":      user.Email,
			"token":      user.EmailVerificationToken,
			"expires_at": user.EmailVerificationExpiresAt,
		},
	}

	job := &worker.Job{
		ID:      uuid.Must(uuid.NewV7()),
		Payload: em,
		Delay:   0,
	}

	return rs.Queue.Enqueue(task.TaskEmailDelivery, job)
}
