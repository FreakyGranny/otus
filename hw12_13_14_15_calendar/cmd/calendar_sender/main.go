package main

import (
	"os"
)

var configFile string

func main() {
	rootCmd := NewSenderCmd()
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "/etc/calendar/sender_config.yaml", "Path to configuration file")

	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}
