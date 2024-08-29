/*
Copyright © 2024 Spyros Mouchlianitis
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// cpCmd represents the cp command
var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "Copies the password of the specified secret to the clipboard.",
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			fmt.Println(cmd.UsageString())
		case 1:
			err := CopySecret(args[0])
			if err != nil {
				log.Fatal(err)
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
