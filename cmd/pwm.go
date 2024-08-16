package cmd

import (
	"fmt"
	"github.com/spyrosmoux/pwm/helpers"
	"log"
	"os"
)

const cipherKey = "thisis32bitlongpassphraseimusing" // TODO(spyrosmoux) make this secret

func CreateUserPassSecret() string {
	username := helpers.StringInput("Enter username: ")
	password := helpers.SecretInput("Enter password: ")

	plaintext := fmt.Sprintf("%s\n%s", username, password)
	hex := helpers.EncryptAES([]byte(cipherKey), plaintext)

	dstPath, err := storeFile(hex)
	if err != nil {
		log.Fatal(err)
	}

	return "Secret created at " + dstPath
}

func CreateEmailPassSecret() string {
	// TODO(spyrosmoux) implement email password recipe
	panic("implement me")
}

// storeFile stores a secret in a default or user-defined directory
// provided by the --location flag
func storeFile(hex string) (string, error) {
	_, err := os.Stat(storageLocation)
	if err != nil {
		err := os.Mkdir(storageLocation, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	secretName := helpers.StringInput("Name your secret: ")
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
