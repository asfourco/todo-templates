package main

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/asfourco/todo-templates/backend/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dbMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database to the latest version",
	RunE:  dbMigrateE,
	Args:  cobra.ExactArgs(0),
}

func init() {
	dbCmd.AddCommand(dbMigrateCmd)
}

func dbMigrateE(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	dbHost := viper.GetString("global-db-host")
	dbPort := viper.GetString("global-db-port")
	dbUser := viper.GetString("global-db-user")
	dbPassword := viper.GetString("global-db-password")
	dbName := viper.GetString("global-db-name")
	migrationsDir := viper.GetString("db-migrations-dir")

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	migrator := db.NewMigration(dbUrl, migrationsDir, nil)

	zlog.Info("Migrating database", zap.String("dbUrl", dbUrl))
	err := migrator.Migrate()
	if err != nil {
		zlog.Error("Migration failed", zap.Error(err))
		return err
	}

	zlog.Info("Migration successful")
	return nil
}
