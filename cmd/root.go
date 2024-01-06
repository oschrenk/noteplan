/*
Copyright © 2024 Oliver Schrenk
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "noteplan",
	Short: "Fetches todocount from noteplan",
	Long:  `Fetches todocount from noteplan long`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// initialize here
}

func initConfig() {
	// read config file here
}
