package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "tsdw",
	Short: "Terraform state dependency walker",
}

// RootCmd.SetHelpFunc()

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
