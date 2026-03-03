/*
Copyright © 2026 Spyros Mouchlianitis
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <your_secret>",
	Short: "Create a new secret",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.UsageString()
			os.Exit(1)
		}
		msg := CreateSecret(args[0])
		if msg != "" {
			println(msg)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
