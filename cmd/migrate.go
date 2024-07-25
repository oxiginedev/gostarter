package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func migrateCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
	}

	cmd.AddCommand(createCommand())

	return cmd
}

func createCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"C"},
		Short:   "create new sql migration file",
		Run: func(cmd *cobra.Command, args []string) {
			migrationFile := fmt.Sprintf("migrations/%v.sql", time.Now().Unix())
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
