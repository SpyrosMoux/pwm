/*
Copyright © 2026 Spyros Mouchlianitis
*/
package cmd

import (
	"github.com/SpyrosMoux/pwm/internal/helpers"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists all secrets located in the default location.",
	Run: func(cmd *cobra.Command, args []string) {
		helpers.PrintInfo(storageLocation)
		err := ListSecrets(storageLocation, 0)
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
