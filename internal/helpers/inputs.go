/*
Copyright © 2026 Spyros Mouchlianitis
*/
package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

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

		if input != "" {
			break
		}
	}

	fmt.Println()
	return strings.TrimSpace(input)
}

func OptionalSecretInput(inputLabel string) string {
	var input string

	_, err := fmt.Fprint(os.Stderr, inputLabel+" ")
	if err != nil {
		panic(err)
	}

	i, _ := term.ReadPassword(int(syscall.Stdin))
	input = string(i)

	fmt.Println()
	return strings.TrimSpace(input)
}
