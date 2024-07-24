package jwt

import (
	"github/oxiginedev/gostarter/config"
	"github/oxiginedev/gostarter/internal/database"
	"github/oxiginedev/gostarter/util"
	"time"

	"github.com/golang-jwt/jwt"
)

var DefaultJWTConfiguration = &JWT{
	Secret:        "starter-secret",
	Expiry:        1800,
	RefreshSecret: "starter-refresh-secret",
	RefreshExpiry: 86400,
}

type JWT struct {
	Secret        string
	Expiry        int
	RefreshSecret string
	RefreshExpiry int
}

type Token struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

func NewJWT(cfg *config.JWTConfiguration) *JWT {
	jt := DefaultJWTConfiguration

	if !util.IsStringEmpty(cfg.Secret) {
		jt.Secret = cfg.Secret
	}

	if cfg.Expiry == 0 {
		jt.Expiry = cfg.Expiry
	}

	if !util.IsStringEmpty(cfg.RefreshSecret) {
		jt.Secret = cfg.RefreshSecret
	}

	if cfg.RefreshExpiry == 0 {
		jt.RefreshExpiry = cfg.RefreshExpiry
	}

	return jt
}

func (j *JWT) GenerateAccessToken(user *database.User) (*Token, error) {
	aclaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Second * time.Duration(j.Expiry)).Unix(),
	})

	token := &Token{}

	accessToken, err := aclaims.SignedString([]byte(j.Secret))
	if err != nil {
		return token, err
	}

	rclaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Second * time.Duration(j.RefreshExpiry)).Unix(),
	})

	refreshToken, err := rclaims.SignedString([]byte(j.RefreshSecret))
	if err != nil {
		return token, err
	}

	token.Access = accessToken
	token.Refresh = refreshToken

	return token, nil
}
