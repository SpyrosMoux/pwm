/*
Copyright © 2026 Spyros Mouchlianitis
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/SpyrosMoux/pwm/internal/helpers"
	"github.com/SpyrosMoux/pwm/internal/models"
	"golang.design/x/clipboard"
)

type Secreter interface {
	Encrypt([]byte) error
	Decrypt([]byte) error
}

func CreateSecret(secretName string) string {
	url := helpers.StringInput("Enter a url for your secret: ")
	username := helpers.StringInput("Enter username: ")
	password := helpers.SecretInput("Enter password ('a' to autogenerate): ")
	description := helpers.StringInput("Enter a description: ")

	secret := models.Secret{
		Name:        secretName,
		Url:         url,
		Username:    username,
		Password:    password,
		Description: description,
	}

	err := Secreter.Encrypt(&secret, []byte(cipherKey))
	if err != nil {
		helpers.PrintError(err.Error())
	}

	jsonSecret, err := json.Marshal(secret)
	if err != nil {
		helpers.PrintError(err.Error())
	}

	dstPath, err := storeFile(secret.Name, jsonSecret)
	if err != nil {
		helpers.PrintError(err.Error())
	}

	return "Secret created at " + dstPath
}

// ListSecrets prints a tree with the files stored in the
// default or user-defined directory provided by the --location flag
func ListSecrets(path string, level int) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for i, file := range files {
		prefix := strings.Repeat("│   ", level)
		if i == len(files)-1 {
			fmt.Printf("%s└── %s\n", prefix, file.Name())
		} else {
			fmt.Printf("%s├── %s\n", prefix, file.Name())
		}

		if file.IsDir() {
			err := ListSecrets(filepath.Join(storageLocation, file.Name()), level+1)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// storeFile stores a secret in a default or user-defined directory
// provided by the --location flag
func storeFile(secretName string, hex []byte) (string, error) {
	_, err := os.Stat(storageLocation)
	if err != nil {
		err := os.Mkdir(storageLocation, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	dstPath := storageLocation + "/" + secretName

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}

	_, err = dstFile.Write([]byte(hex))
	if err != nil {
		return "", err
	}

	defer dstFile.Close()

	return dstPath, nil
}

// GetSecret reads a given secret, decrypts it and returns it
// Currently does not support subdirectories
// TODO(spyrosmoux) should be able to get secrets in subdirectories
func GetSecret(secret string) (decryptedSecret string, err error) {
	hex, err := os.ReadFile(storageLocation + "/" + secret)
	if err != nil {
		return
	}

	var jsonSecret models.Secret
	err = json.Unmarshal(hex, &jsonSecret)
	if err != nil {
		return
	}

	err = Secreter.Decrypt(&jsonSecret, []byte(cipherKey))
	if err != nil {
		return "", err
	}

	return jsonSecret.String(), nil
}

// resolveSecretArg accepts either a secret name or a 1-based index string
// and returns the resolved secret filename.
func resolveSecretArg(arg string) (string, error) {
	// try parse as integer index
	idx, err := strconv.Atoi(arg)
	if err != nil {
		// not an integer, assume it's a name (possibly with subdir)
		return arg, nil
	}

	if idx <= 0 {
		return "", errors.New("index must be >= 1")
	}

	files, err := getFilesSortedByModTime(storageLocation)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", errors.New("no secrets available")
	}

	if idx > len(files) {
		return "", errors.New("index out of range")
	}

	return files[idx-1], nil
}

// collectFiles walks the storage directory and returns a slice of file paths
// relative to the root. Directories are skipped.
func collectFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// getFilesSortedByModTime returns a list of files (relative paths) under root
// sorted by modification time ascending (oldest first). If mod times are equal
// it falls back to lexical order.
func getFilesSortedByModTime(root string) ([]string, error) {
	type fileEntry struct {
		rel     string
		modTime time.Time
	}

	var entries []fileEntry
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		entries = append(entries, fileEntry{rel: rel, modTime: info.ModTime()})
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].modTime.Equal(entries[j].modTime) {
			return entries[i].rel < entries[j].rel
		}
		return entries[i].modTime.Before(entries[j].modTime)
	})

	var files []string
	for _, e := range entries {
		files = append(files, e.rel)
	}
	return files, nil
}

func RemoveSecret(secret string) error {
	_, err := os.Stat(storageLocation + "/" + secret)
	if err != nil {
		return err
	}

	err = os.Remove(storageLocation + "/" + secret)
	if err != nil {
		return err
	}

	return nil
}

// ListSecretsNumbered prints a tree with numeric indices for files.
// Indices correspond to the sorted list used by resolveSecretArg.
func ListSecretsNumbered(root string, level int) error {
	files, err := getFilesSortedByModTime(root)
	if err != nil {
		return err
	}

	indexMap := make(map[string]int)
	for i, f := range files {
		indexMap[f] = i + 1
	}

	// Build dirMin map: smallest index for files under each directory
	dirMin := make(map[string]int)
	const maxInt = int(^uint(0) >> 1)
	dirMin["."] = maxInt
	for rel, idx := range indexMap {
		if idx < dirMin["."] {
			dirMin["."] = idx
		}
		dir := filepath.Dir(rel)
		for dir != "." && dir != string(filepath.Separator) && dir != "" {
			cur, ok := dirMin[dir]
			if !ok || idx < cur {
				dirMin[dir] = idx
			}
			dir = filepath.Dir(dir)
		}
	}

	var printTree func(path string, level int) error
	printTree = func(path string, level int) error {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		type node struct {
			entry os.DirEntry
			key   int
		}
		nodes := make([]node, 0, len(entries))
		for _, entry := range entries {
			if entry.IsDir() {
				relDir, _ := filepath.Rel(root, filepath.Join(path, entry.Name()))
				key := dirMin[relDir]
				if key == 0 {
					key = maxInt
				}
				nodes = append(nodes, node{entry: entry, key: key})
			} else {
				relFile, _ := filepath.Rel(root, filepath.Join(path, entry.Name()))
				key := indexMap[relFile]
				if key == 0 {
					key = maxInt
				}
				nodes = append(nodes, node{entry: entry, key: key})
			}
		}

		sort.Slice(nodes, func(i, j int) bool {
			if nodes[i].key == nodes[j].key {
				return nodes[i].entry.Name() < nodes[j].entry.Name()
			}
			return nodes[i].key < nodes[j].key
		})

		for i, n := range nodes {
			entry := n.entry
			prefix := strings.Repeat("│   ", level)
			isLast := i == len(nodes)-1
			connector := "├── "
			if isLast {
				connector = "└── "
			}

			if entry.IsDir() {
				fmt.Printf("%s%s%s\n", prefix, connector, entry.Name())
				err := printTree(filepath.Join(path, entry.Name()), level+1)
				if err != nil {
					return err
				}
			} else {
				rel, err := filepath.Rel(root, filepath.Join(path, entry.Name()))
				if err != nil {
					return err
				}
				idx := indexMap[rel]
				if idx > 0 {
					fmt.Printf("%s%s{%d} %s\n", prefix, connector, idx, entry.Name())
				} else {
					fmt.Printf("%s%s%s\n", prefix, connector, entry.Name())
				}
			}
		}

		return nil
	}

	return printTree(root, level)
}

func CopySecret(secretName string) error {
	helpers.PrintInfo("Copying secret " + secretName)

	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		return err
	}

	secret, err := readSecretIntoStruct(secretName)
	if err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtText, []byte(secret.Password))

	return nil
}

func readSecretIntoStruct(secret string) (models.Secret, error) {
	hex, err := os.ReadFile(storageLocation + "/" + secret)
	if err != nil {
		return models.Secret{}, err
	}

	var jsonSecret models.Secret
	err = json.Unmarshal(hex, &jsonSecret)
	if err != nil {
		return models.Secret{}, err
	}

	err = Secreter.Decrypt(&jsonSecret, []byte(cipherKey))
	if err != nil {
		return models.Secret{}, err
	}

	return jsonSecret, nil
}
