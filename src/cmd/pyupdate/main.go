package main

import (
	"github.com/ashishb/pyupdate/src/internal/util/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func main() {
	logger.ConfigureLogging(true)
	var rootCmd = &cobra.Command{
		Use:   "pyupdate",
		Short: "A tool to update Python packages",
		Long:  "pyupdate is a command-line tool that helps you update Python packages in your environment.",
	}

	dirPath := rootCmd.PersistentFlags().StringP("directory", "d", ".", "Path to directory containing pyproject.toml")

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		log.Info().
			Str("directory", *dirPath).
			Msg("Updating Python packages")
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to execute command")
	}
}
