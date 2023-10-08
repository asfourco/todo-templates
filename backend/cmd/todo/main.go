package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "cmd", Short: "Todo Api Server"}
var dbCmd = &cobra.Command{Use: "db", Short: "Database commands"}

func init() {
	rootCmd.PersistentFlags().String("db-host", getEnv("DB_HOST", "localhost"), "Database host")
	rootCmd.PersistentFlags().String("db-port", getEnv("DB_PORT", "5432"), "Database port")
	rootCmd.PersistentFlags().String("db-user", getEnv("DB_USER", "postgres"), "Database user")
	rootCmd.PersistentFlags().String("db-password", getEnv("DB_PASSWORD", "postgres"), "Database password")

	dbCmd.PersistentFlags().String("migrations-dir", getEnv("DB_MIGRATIONS_DIR", "backend/db/migrations"), "Migrations directory")

	cobra.OnInitialize(func() {
		autoBind(rootCmd, "cmd")
	})

	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionsCmd)
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
