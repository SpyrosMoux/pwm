package secrets

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/SpyrosMoux/passwdgen"
	"github.com/SpyrosMoux/pwm/internal/helpers"
	"github.com/SpyrosMoux/pwm/internal/models"
)

func (svc *Service) CreateSecret(secretName string) string {
	// name must not be purely numeric, otherwise later commands that accept
	// either a name or numeric index would misinterpret the argument.
	if _, err := strconv.Atoi(secretName); err == nil {
		helpers.PrintError("secret name cannot be a number")
	}

	url := helpers.StringInput("Enter a url for your secret: ")
	username := helpers.StringInput("Enter username: ")
	password := helpers.SecretInput("Enter password ('a' to autogenerate): ")
	if password == "a" {
		options := passwdgen.NewRandomStringOptions()
		password = passwdgen.RandomStringNumbersSymbols(&options)
		fmt.Printf("Generated password: %s\n", password)
	}
	description := helpers.StringInput("Enter a description: ")

	secret := models.Secret{
		Name:        secretName,
		Url:         url,
		Username:    username,
		Password:    password,
		Description: description,
	}

	err := svc.Crypto.Encrypt(&secret)
	if err != nil {
		helpers.PrintError(err.Error())
	}

	jsonSecret, err := json.Marshal(secret)
	if err != nil {
		helpers.PrintError(err.Error())
	}

	dstPath, err := svc.Store.StoreFile(secret.Name, jsonSecret)
	if err != nil {
		helpers.PrintError(err.Error())
	}

	return "Secret created at " + dstPath
}
