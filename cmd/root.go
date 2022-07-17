/*
Copyright © 2022 codeschooldropout code@cay.io

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	Version = "0.0.3"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "3gophers",
	Short: "My own personal bloomberg terminal or something",
	Long:  `3gopher is a terminal application that allows you to interact 
	with tradingview webhooks and act on them with other apis. Additionally,
	it can be used to show various information about the market and your portfolio.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = Version
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.3gophers.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
