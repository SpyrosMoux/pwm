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
