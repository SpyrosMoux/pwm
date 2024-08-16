package helpers

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"
	"syscall"
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
