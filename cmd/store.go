package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Store a username and password",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Result: %s\n", Store())
	},
}

func init() {
	rootCmd.AddCommand(storeCmd)
}
