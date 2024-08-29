/*
Copyright © 2024 Spyros Mouchlianitis
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm <your_secret>",
	Short: "Removes a secret",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := RemoveSecret(args[0])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Removed secret:", args[0])
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
