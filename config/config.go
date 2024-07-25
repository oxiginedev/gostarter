package config

import (
	"fmt"
	"github/oxiginedev/gostarter/util"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type DatabaseConfiguration struct {
	Driver         string `envconfig:"DB_DRIVER"`
	Host           string `envconfig:"DB_HOST"`
	Username       string `envconfig:"DB_USERNAME"`
	Password       string `envconfig:"DB_PASSWORD"`
	Database       string `envconfig:"DB_DATABASE"`
	Port           uint   `envconfig:"DB_PORT"`
	MaxPool        int    `envconfig:"DB_MAX_POOL"`
	MaxIdlePool    int    `envconfig:"DB_MAX_IDLE_POOL"`
	MaxLifetime    int    `envconfig:"DB_MAX_LIFETIME"`
	MaxIdleTime    int    `envconfig:"DB_MAX_IDLE_TIME"`
	Options        string `envconfig:"DB_OPTIONS"`
	MigrationsPath string `envconfig:"DB_MIGRATIONS_PATH"`
}

func (dc DatabaseConfiguration) BuildDSN() string {
	if dc.Driver == "" {
		return ""
	}

	authPart := ""
	if dc.Username != "" || dc.Password != "" {
		authPrefix := url.UserPassword(dc.Username, dc.Password)
		authPart = fmt.Sprintf("%s@", authPrefix)
	}

	dbPart := ""
	if dc.Database != "" {
		dbPart = fmt.Sprintf("/%s", dc.Database)
	}

	return fmt.Sprintf("%s://%s%s:%d%s", dc.Driver, authPart, dc.Host, dc.Port, dbPart)
}

type JWTConfiguration struct {
	Secret        string `envconfig:"JWT_SECRET"`
	Expiry        int    `envconfig:"JWT_EXPIRY"`
	RefreshSecret string `envconfig:"JWT_REFRESH_SECRET"`
	RefreshExpiry int    `envconfig:"JWT_REFRESH_EXPIRY"`
}

type HTTPConfiguration struct {
	Port uint32 `envconfig:"HTTP_PORT"`
}

type Configuration struct {
	HTTP     HTTPConfiguration
	Database DatabaseConfiguration
	JWT      JWTConfiguration
}

func Load(p string) (*Configuration, error) {
	err := loadFromEnvironment(p)
	if err != nil {
		return nil, err
	}

	config := &Configuration{}

	err = envconfig.Process("", config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func loadFromEnvironment(p string) error {
	var err error

	if !util.IsStringEmpty(p) {
		err = godotenv.Load(p)
	} else {
		err = godotenv.Load()

		// check if .env file does not exist, this is OK
		if os.IsNotExist(err) {
			return nil
		}
	}

	return err
}
