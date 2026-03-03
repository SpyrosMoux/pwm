/*
Copyright © 2026 Spyros Mouchlianitis
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/SpyrosMoux/pwm/internal/helpers"
	"github.com/spf13/cobra"
)

// cpCmd represents the cp command
var cpCmd = &cobra.Command{
	Use:   "cp <name|index>",
	Short: "Copies the password of the specified secret (name or index) to the clipboard.",
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			cmd.UsageString()
			os.Exit(1)
		case 1:
			name, err := resolveSecretArg(args[0])
			if err != nil {
				helpers.PrintError(err.Error())
				return
			}

			err = CopySecret(name)
			if err != nil {
				helpers.PrintError(err.Error())
			}
		default:
			fmt.Println(cmd.UsageString())
		}
	},
}

func init() {
	rootCmd.AddCommand(cpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
