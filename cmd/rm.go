/*
Copyright © 2026 Spyros Mouchlianitis
*/
package cmd

import (
	"log"
	"os"

	"github.com/SpyrosMoux/pwm/internal/helpers"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm <name|index>",
	Short: "Removes a secret by name or numeric index (from `pwm ls`)",
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

		err = secretsService.RemoveSecret(name)
		if err != nil {
			log.Fatal(err)
		}

		helpers.PrintInfo("Removed secret: " + name)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
