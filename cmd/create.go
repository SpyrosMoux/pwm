/*
Copyright © 2026 Spyros Mouchlianitis
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <your_secret>",
	Short: "Create a new secret",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		msg := CreateSecret(args[0])
		if msg != "" {
			println(msg)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
