package cmd

import (
	"github.com/spyrosmoux/pwm/helpers"
	"os"
)

const cipherKey = "thisis32bitlongpassphraseimusing"

func CreateUserPassSecret() string {
	username := helpers.StringInput("Enter username: ")
	password := helpers.SecretInput("Enter password: ")

	plaintext := username + ":" + password
	hex := helpers.EncryptAES([]byte(cipherKey), plaintext)

	dstPath := helpers.StringInput("Enter destination filepath: ")
	dstFile, err := os.Create(dstPath)
	if err != nil {
		panic(err)
	}

	_, err = dstFile.Write([]byte(hex))
	if err != nil {
		panic(err)
	}

	defer dstFile.Close()

	return "Secret created at " + dstPath
}

func CreateEmailPassSecret() string {
	// TODO(spyrosmoux) implement email password recipe
	panic("implement me")
}
