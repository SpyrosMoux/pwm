/*
Copyright © 2026 Spyros Mouchlianitis
*/
package cmd

import (
	"errors"
	"fmt"
	"strconv"
)

// resolveSecretArg turns an argument into a secret filename. Numeric
// strings are always treated as indices; everything else is returned
// as user typed it.
func resolveSecretArg(arg string) (string, error) {
	idx, err := strconv.ParseInt(arg, 10, 64)
	if err == nil {
		// parsed as integer, treat strictly as index
		if idx < 1 {
			return "", errors.New("index must be >= 1")
		}

		files, err := secretsService.Store.GetFilesSortedByModTime(storageLocation)
		if err != nil {
			return "", err
		}

		if len(files) == 0 {
			return "", errors.New("no secrets available")
		}

		if int(idx) > len(files) {
			return "", fmt.Errorf("index out of range (1..%d)", len(files))
		}

		return files[idx-1], nil
	}

	// parsing failed; inspect error
	if numErr, ok := err.(*strconv.NumError); ok {
		switch numErr.Err {
		case strconv.ErrRange:
			return "", errors.New("invalid index: number too large")
		case strconv.ErrSyntax:
			// not an integer at all, treat as a name
			return arg, nil
		}
	}

	// fallback to name (covers any other unexpected error)
	return arg, nil
}
