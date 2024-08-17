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
	"github.com/spyrosmoux/pwm/helpers"
	"github.com/spyrosmoux/pwm/models"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const cipherKey = "thisis32bitlongpassphraseimusing" // TODO(spyrosmoux) make this secret

func CreateSecret(secretName string) string {
	url := helpers.StringInput("Enter a url for your secret: ")
	username := helpers.StringInput("Enter username: ")
	password := helpers.SecretInput("Enter password: ")
	description := helpers.StringInput("Enter a description: ")

	secret := models.Secret{
		Name:        secretName,
		Url:         url,
		Username:    username,
		Password:    password,
		Description: description,
	}

	plaintext := secret.ToString()
	hex := helpers.EncryptAES([]byte(cipherKey), plaintext)

	dstPath, err := storeFile(secret.Name, hex)
	if err != nil {
		log.Fatal(err)
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
func storeFile(secretName string, hex string) (string, error) {
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
func GetSecret(secret string) (string, error) {
	hex, err := os.ReadFile(storageLocation + "/" + secret)
	if err != nil {
		return "", err
	}

	decryptedSecret, err := helpers.DecryptAES([]byte(cipherKey), string(hex))
	if err != nil {
		return "", err
	}

	return decryptedSecret, nil
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
