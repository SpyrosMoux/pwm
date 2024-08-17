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

func DeleteSecret(secret string) error {
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
