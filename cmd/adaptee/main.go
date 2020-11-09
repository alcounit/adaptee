package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alcounit/adaptee"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//command ...
func command() *cobra.Command {
	var (
		address         string
		selenosisURL    string
		shutdownTimeout time.Duration
	)
	cmd := &cobra.Command{
		Use:   "adaptee",
		Short: "adaptee is a adaptor sidecar for selenoid ui",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logrus.New()
			logger.Infof("starting selenoid ui adaptor: %s", address)

			_, err := url.Parse(selenosisURL)
			if err != nil {
				logger.Fatalf("failed to parse selenosis url: %v", err)
			}

			app := adaptee.New(logger, adaptee.Configuration{
				SelenosisURL: selenosisURL,
			})

			router := mux.NewRouter()
			router.HandleFunc("/status", app.HandleStatus)
			router.HandleFunc("/vnc/{sessionId}", app.HandleWs)
			router.HandleFunc("/logs/{sessionId}", app.HandleWs)

			srv := &http.Server{
				Addr:    address,
				Handler: router,
			}

			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

			e := make(chan error)
			go func() {
				e <- srv.ListenAndServe()
			}()

			select {
			case err := <-e:
				logger.Fatalf("failed to start adaptee: %v", err)
			case <-stop:
				logger.Warn("stopping adaptee")
			}

			ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				logger.Fatalf("failed to stop adaptee", err)
			}
		},
	}

	cmd.Flags().StringVar(&address, "port", ":4444", "adaptee port")
	cmd.Flags().StringVar(&selenosisURL, "selenosis-url", "http://selenosis-local:4444", "selenosis url")
	cmd.Flags().DurationVar(&shutdownTimeout, "graceful-shutdown-timeout", 30*time.Second, "time in seconds  gracefull shutdown timeout")

	return cmd
}

func main() {
	if err := command().Execute(); err != nil {
		os.Exit(1)
	}
}
