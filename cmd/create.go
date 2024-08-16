package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var userPassRecipe bool
var emailPassRecipe bool
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new secret",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if userPassRecipe {
			CreateUserPassSecret()
			os.Exit(0)
		}
		if emailPassRecipe {
			CreateEmailPassSecret()
			os.Exit(0)
		}

		err := cmd.Usage()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	createCmd.Flags().BoolVarP(&userPassRecipe, "user", "u", false, "Create user/password secret")
	createCmd.Flags().BoolVarP(&emailPassRecipe, "email", "e", false, "Create email/password secret")
	rootCmd.AddCommand(createCmd)
}
