package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "throttler",
	Short: "throttler service",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("throttler service")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
