/*
Copyright © 2024 Spyros Mouchlianitis
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists all secrets located in the default location.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(storageLocation)
		err := ListSecrets(storageLocation, 0)
		if err != nil {
			log.Fatalf(err.Error())
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
