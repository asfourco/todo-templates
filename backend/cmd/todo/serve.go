package main

import (
	"fmt"

	"os"
	"time"

	"github.com/asfourco/todo-templates/backend/api"
	"github.com/asfourco/todo-templates/backend/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/streamingfast/derr"
	"github.com/streamingfast/shutter"
)

var serveCmd = &cobra.Command{Use: "serve", Short: "starts the HTTP API server", RunE: serveE}

func init() {
	serveCmd.Flags().String("port", ":8080", "HTTP port")
}

func serveE(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true
	ctx := cmd.Context()

	port := viper.GetString("serve-port")
	dbHost := viper.GetString("global-db-host")
	dbPort := viper.GetString("global-db-port")
	dbUser := viper.GetString("global-db-user")
	dbName := viper.GetString("global-db-name")
	dbPassword := viper.GetString("global-db-password")

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	postgresClient, err := db.NewPostgresClient(ctx, dbUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	server, err := api.NewServer(ctx, port, postgresClient)
	if err != nil {
		return fmt.Errorf("failed to create HTTP API server: %w", err)
	}

	app := shutter.New()
	app.OnTerminating(func(_ error) {})

	zlog.Info("starting HTTP server")
	err = server.Start(app)
	if err != nil {
		return fmt.Errorf("Unable to start HTTP API server: %w", err)
	}

	signalHandler := derr.SetupSignalHandler(0 * time.Second)

	zlog.Info("read, waiting for signal to quit")

	select {
	case <-signalHandler:
		zlog.Info("received signal, quitting")
		go app.Shutdown(nil)
	case <-app.Terminating():
		if app.Err() != nil {
			fmt.Printf("app terminated with error: %s\n", app.Err())
			os.Exit(1)
		}
	}

	zlog.Info("waiting for app termination")
	select {
	case <-app.Terminated():
	case <-ctx.Done():
	case <-time.After(10 * time.Second):
		zlog.Error("app did not terminate within 10s, forcing exit")
	}

	zlog.Info("app terminated")
	return nil
}
