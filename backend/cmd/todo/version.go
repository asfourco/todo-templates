package main

import (
	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"

var versionsCmd = &cobra.Command{
	Use:   "version",
	Short: "version of the todo CLI tool",
	RunE:  versionE,
	Args:  cobra.ExactArgs(0),
}

func init() {}

func versionE(cmd *cobra.Command, args []string) error {
	zlog.Info("Todo CLI tool version", zap.String("version", VERSION))
	return nil
}
