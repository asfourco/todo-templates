package main

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/asfourco/todo-templates/backend/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dbRollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback the database to the previous version",
	RunE:  dbRollbackE,
	Args:  cobra.ExactArgs(0),
}

func init() {
	dbCmd.AddCommand(dbRollbackCmd)
}

func dbRollbackE(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	dbHost := viper.GetString("global-db-host")
	dbPort := viper.GetString("global-db-port")
	dbUser := viper.GetString("global-db-user")
	dbPassword := viper.GetString("global-db-password")
	dbName := viper.GetString("global-db-name")
	migrationsDir := viper.GetString("db-migrations-dir")

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	migrator := db.NewMigration(dbUrl, migrationsDir, nil)

	zlog.Info("Rolling back database", zap.String("dbUrl", dbUrl))

	err := migrator.Rollback()
	if err != nil {
		zlog.Error("Rollback failed", zap.Error(err))
		return err
	}
	zlog.Info("Rollback successful")
	return nil
}
