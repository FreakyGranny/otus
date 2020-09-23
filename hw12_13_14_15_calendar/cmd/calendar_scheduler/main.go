package main

import (
	"os"
)

var configFile string

func main() {
	rootCmd := NewSchedulerCmd()
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "/etc/calendar/scheduler_config.yaml", "Path to configuration file")

	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}
