/*
Copyright © 2024 Spyros Mouchlianitis

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var storageLocation string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pwm",
	Short:   "A simple password management tool",
	Example: "pwm <my_secret>  Will print the decrypted secret",
	Args:    cobra.MinimumNArgs(0),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			cmd.HelpFunc()(cmd, args)
			os.Exit(0)
		case 1:
			secret, err := GetSecret(args[0])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(secret)
			os.Exit(0)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	homeDir, _ := os.UserHomeDir()
	storageLocation = homeDir + "/.pwm"

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pwm.yaml)")
	rootCmd.PersistentFlags().StringVar(&storageLocation, "location", storageLocation, "The directory to store the "+
		"secrets.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
