package jwt

import (
	"errors"
	"fmt"
	"github/oxiginedev/gostarter/config"
	"github/oxiginedev/gostarter/internal/database"
	"github/oxiginedev/gostarter/util"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("expired token")
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

type ValidatedToken struct {
	UserID string
	Expiry int64
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
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Second * time.Duration(j.Expiry)).Unix(),
	})

	token := &Token{}

	accessToken, err := claims.SignedString([]byte(j.Secret))
	if err != nil {
		return token, err
	}

	claims = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Second * time.Duration(j.RefreshExpiry)).Unix(),
	})

	refreshToken, err := claims.SignedString([]byte(j.RefreshSecret))
	if err != nil {
		return token, err
	}

	token.Access = accessToken
	token.Refresh = refreshToken

	return token, nil
}

func (j *JWT) ValidateAccessToken(accessToken string) (*ValidatedToken, error) {
	return j.validateToken(accessToken, j.Secret)
}

func (j *JWT) ValidateRefreshToken(accessToken string) (*ValidatedToken, error) {
	return j.validateToken(accessToken, j.RefreshSecret)
}

func (j *JWT) validateToken(accessToken, secret string) (*ValidatedToken, error) {
	var userId string
	var expiry float64

	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method - %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		v, ok := err.(*jwt.ValidationError)
		if ok && v.Errors == jwt.ValidationErrorExpired {
			if payload, ok := token.Claims.(jwt.MapClaims); ok {
				expiry = payload["exp"].(float64)
			}

			return &ValidatedToken{Expiry: int64(expiry)}, ErrTokenExpired
		}

		return nil, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId = payload["sub"].(string)
		expiry = payload["exp"].(float64)

		v := &ValidatedToken{UserID: userId, Expiry: int64(expiry)}
		return v, nil
	}

	return nil, err
}
