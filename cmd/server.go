package cmd

import (
	"context"
	"fmt"
	"github/oxiginedev/gostarter/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func serverCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"S"},
		Run: func(cmd *cobra.Command, args []string) {
			var h http.Handler

			a := api.NewAPI(&api.APIOptions{
				DB:     app.database,
				Config: app.config,
			})
			h = a.BuildAPIRoutes()

			srv := &http.Server{
				Addr:              fmt.Sprintf(":%d", app.config.HTTP.Port),
				Handler:           h,
				ReadHeaderTimeout: time.Second * 2,
				ReadTimeout:       time.Second * 15,
				WriteTimeout:      time.Second * 15,
			}

			log.Printf("server listening on port [:%d]\n", app.config.HTTP.Port)

			go func() {
				err := srv.ListenAndServe()
				if err != nil && err != http.ErrServerClosed {
					log.Fatalf("server failed to start - %v", err)
				}
			}()

			stop := make(chan os.Signal, 1)
			signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

			<-stop

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := srv.Shutdown(ctx)
			if err != nil {
				log.Fatalf("server failed to shutdown - %v", err)
			}

			log.Println("server shutdown gracefully")
		},
	}

	return cmd
}
