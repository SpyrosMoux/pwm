/*
Copyright © 2024 Spyros Mouchlianitis
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/SpyrosMoux/pwm/internal/helpers"
	"github.com/SpyrosMoux/pwm/internal/models"
	"golang.design/x/clipboard"
)

type Secreter interface {
	Encrypt([]byte) error
	Decrypt([]byte) error
}

// TODO(spyrosmoux) make this secret
// perhaps use an init command to set it, along with the storage location ?
// could also be useful in case of master password!
const cipherKey = "thisis32bitlongpassphraseimusing"

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

	err := Secreter.Encrypt(&secret, []byte(cipherKey))
	if err != nil {
		log.Fatal(err)
	}

	jsonSecret, err := json.Marshal(secret)
	if err != nil {
		log.Fatal(err)
	}

	dstPath, err := storeFile(secret.Name, jsonSecret)
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

func CopySecret(secretName string) error {
	fmt.Println("Copying secret " + secretName)

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
