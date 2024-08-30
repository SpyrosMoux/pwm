/*
Copyright © 2024 Spyros Mouchlianitis
*/
package models

import (
	"fmt"

	"github.com/SpyrosMoux/pwm/internal/crypto"
)

type Secret struct {
	Name        string
	Url         string
	Username    string
	Password    string
	Description string
}

func (secret Secret) String() string {
	return fmt.Sprintf("Name: %s\n"+"Url: %s\n"+"Username: %s\n"+"Password: %s\n"+"Description: %s\n",
		secret.Name, secret.Url, secret.Username, secret.Password,
		secret.Description)
}

func (secret Secret) Encrypt(cipherKey []byte) error {
	urlHex, err := crypto.EncryptAES(cipherKey, secret.Url)
	if err != nil {
		return err
	}

	usernameHex, err := crypto.EncryptAES(cipherKey, secret.Username)
	if err != nil {
		return err
	}

	passwordHex, err := crypto.EncryptAES(cipherKey, secret.Password)
	if err != nil {
		return err
	}

	descriptionHex, err := crypto.EncryptAES(cipherKey, secret.Description)
	if err != nil {
		return err
	}

	secret.Url = urlHex
	secret.Username = usernameHex
	secret.Password = passwordHex
	secret.Description = descriptionHex

	return nil
}

func (secret Secret) Decrypt(cipherKey []byte) error {
	urlString, err := crypto.DecryptAES(cipherKey, secret.Url)
	if err != nil {
		return err
	}

	usernameString, err := crypto.DecryptAES(cipherKey, secret.Username)
	if err != nil {
		return err
	}

	passwordString, err := crypto.DecryptAES(cipherKey, secret.Password)
	if err != nil {
		return err
	}

	descriptionString, err := crypto.DecryptAES(cipherKey, secret.Description)
	if err != nil {
		return err
	}

	secret.Url = urlString
	secret.Username = usernameString
	secret.Password = passwordString
	secret.Description = descriptionString

	return nil
}
