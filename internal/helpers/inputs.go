/*
Copyright © 2024 Spyros Mouchlianitis
*/
package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/SpyrosMoux/passwdgen"
	"golang.org/x/term"
)

func StringInput(inputLabel string) string {
	var input string
	r := bufio.NewReader(os.Stdin)

	for {
		_, err := fmt.Fprint(os.Stderr, inputLabel+" ")
		if err != nil {
			panic(err)
		}

		input, _ = r.ReadString('\n')
		if input != "" {
			break
		}
	}

	return strings.TrimSpace(input)
}

func SecretInput(inputLabel string) string {
	var input string

	for {
		_, err := fmt.Fprint(os.Stderr, inputLabel+" ")
		if err != nil {
			panic(err)
		}

		i, _ := term.ReadPassword(int(syscall.Stdin))
		input = string(i)

		options := passwdgen.NewRandomStringOptions()

		if input == "a" {
			input = passwdgen.RandomStringNumbersSymbols(&options)
			fmt.Printf("\nGenerated password: %s", input)
			break
		}

		if input != "" {
			break
		}
	}

	fmt.Println()
	return strings.TrimSpace(input)
}
