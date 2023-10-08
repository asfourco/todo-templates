package main

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func autoBind(root *cobra.Command, prefix string) {
	viper.SetEnvPrefix(strings.ToUpper(prefix))
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	recurseCommands(root, nil)
}

func recurseCommands(root *cobra.Command, segments []string) {
	var segmentPrefix string
	if len(segments) > 0 {
		segmentPrefix = strings.Join(segments, "-") + "-"
	}

	zlog.Debug("re-binding flags", zap.String("cmd", root.Name()), zap.String("prefix", segmentPrefix))
	defer func() {
		zlog.Debug("re-binding flags terminated", zap.String("cmd", root.Name()))
	}()

	root.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		newVar := segmentPrefix + "global-" + flag.Name
		if err := viper.BindPFlag(newVar, flag); err != nil {
			panic(fmt.Errorf("unable to bind global flag: %w", err))
		}
		zlog.Debug("binding persistent flag", zap.String("flag", flag.Name), zap.String("new_var", newVar))
	})

	root.Flags().VisitAll(func(flag *pflag.Flag) {
		newVar := segmentPrefix + flag.Name
		if err := viper.BindPFlag(newVar, flag); err != nil {
			panic(fmt.Errorf("unable to bind namspaced flag: %w", err))
		}
		zlog.Debug("binding flag", zap.String("flag", flag.Name), zap.String("new_var", newVar))
	})

	for _, cmd := range root.Commands() {
		recurseCommands(cmd, append(segments, cmd.Name()))
	}
}
