/*
Copyright © 2026 Spyros Mouchlianitis
*/
package cmd

import (
	"os"

	"github.com/SpyrosMoux/pwm/internal/helpers"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists all secrets located in the default location (prints numeric indices).",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			cmd.UsageString()
			os.Exit(1)
		}
		helpers.PrintInfo(storageLocation)
		err := secretsService.ListSecretsNumbered(storageLocation, 0)
		if err != nil {
			helpers.PrintError(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
