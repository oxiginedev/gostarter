package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/spf13/cobra"
)

func migrateCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
	}

	cmd.AddCommand(upCommand(app))
	cmd.AddCommand(createCommand())

	return cmd
}

func upCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Apply pending migrations",
		Run: func(cmd *cobra.Command, args []string) {
			migrationFiles := app.config.Database.MigrationsPath
			migrator, err := pop.NewFileMigrator(migrationFiles, app.database.Connection)
			if err != nil {
				log.Fatalf("error creating migrator - %v", err)
			}

			migrator.SchemaPath = ""
			err = migrator.Status(os.Stdout)
			if err != nil {
				log.Fatal(err)
			}

			if err := migrator.Up(); err != nil {
				log.Fatalf("migrate up command failed - %v", err)
			}

			log.Println("database migration applied!")
		},
	}

	return cmd
}

func createCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"C"},
		Short:   "create new sql migration file",
		Run: func(cmd *cobra.Command, args []string) {
			migrationFile := fmt.Sprintf("migrations/%v_%s.up.sql", time.Now().Unix(), args[0])
			file, err := os.Create(migrationFile)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			lines := []string{"-- +migrate Up", "-- +migrate Down"}
			for _, line := range lines {
				_, err := file.WriteString(line + "\n\n")
				if err != nil {
					log.Fatal(err)
				}
			}

			path, err := filepath.Abs(migrationFile)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("migration [%s] created successfully", path)
		},
	}

	return cmd
}
