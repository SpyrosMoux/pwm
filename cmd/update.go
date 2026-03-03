/*
Copyright © 2026 Spyros Mouchlianitis
*/
package cmd

import (
	"os"

	"github.com/SpyrosMoux/pwm/internal/helpers"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <name|index>",
	Short: "Update an existing secret (URL, username, password, description)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.UsageString()
			os.Exit(1)
		}
		name, err := resolveSecretArg(args[0])
		if err != nil {
			helpers.PrintError(err.Error())
			return
		}

		err = secretsService.UpdateSecret(name)
		if err != nil {
			helpers.PrintError(err.Error())
			return
		}

		helpers.PrintInfo("Secret updated successfully: " + name)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
