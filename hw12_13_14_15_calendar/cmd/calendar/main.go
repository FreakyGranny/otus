package main

import (
	"os"

	"github.com/spf13/cobra"
)

var configFile string

func main() {
	rootCmd := &cobra.Command{Use: "calendar"}
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "/etc/calendar/config.yaml", "Path to configuration file")

	rootCmd.AddCommand(NewAPICmd(), NewMigrateCmd())
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}
