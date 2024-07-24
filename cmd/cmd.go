package cmd

import (
	"github/oxiginedev/gostarter/config"
	"github/oxiginedev/gostarter/internal/database"
	"github/oxiginedev/gostarter/internal/database/migrations"
	"github/oxiginedev/gostarter/internal/database/postgres"
	"log"
	"time"

	"github.com/spf13/cobra"
)

type App struct {
	cmd      *cobra.Command
	config   *config.Configuration
	database database.Database
}

func Run() error {
	var configFile string

	// Force server time to be in UTC
	time.Local = time.UTC

	app := &App{
		cmd: &cobra.Command{
			Use:   "gostarter",
			Short: "Starter kit for golang applications",
		},
	}

	app.cmd.PersistentFlags().StringVar(&configFile, "config", ".env", "Path to configuration file")

	app.cmd.PersistentPreRunE = persistentPreRunE(app)
	app.cmd.PersistentPostRunE = persistentPostRunE(app)

	app.cmd.AddCommand(serverCommand(app))

	err := app.cmd.Execute()
	if err != nil {
		return err
	}

	return nil
}

func persistentPreRunE(app *App) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}

		config, err := config.Load(configFile)
		if err != nil {
			return err
		}

		database, err := postgres.NewPostgresRepository(&config.Database)
		if err != nil {
			return err
		}

		log.Println("connected to database")

		err = migrations.AutoMigrate(database)
		if err != nil {
			return err
		}

		app.config = config
		app.database = database

		return nil
	}
}

func persistentPostRunE(app *App) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}
}
