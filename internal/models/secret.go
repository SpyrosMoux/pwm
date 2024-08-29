/*
Copyright © 2024 Spyros Mouchlianitis
*/
package models

import (
	"fmt"
	"github.com/spyrosmoux/pwm/internal/crypto"
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

func (secret Secret) Encrypt(cipherKey []byte) (Secret, error) {
	urlHex, err := crypto.EncryptAES(cipherKey, secret.Url)
	if err != nil {
		return Secret{}, err
	}

	usernameHex, err := crypto.EncryptAES(cipherKey, secret.Username)
	if err != nil {
		return Secret{}, err
	}

	passwordHex, err := crypto.EncryptAES(cipherKey, secret.Password)
	if err != nil {
		return Secret{}, err
	}

	descriptionHex, err := crypto.EncryptAES(cipherKey, secret.Description)
	if err != nil {
		return Secret{}, err
	}

	return Secret{
		Name:        secret.Name,
		Url:         urlHex,
		Username:    usernameHex,
		Password:    passwordHex,
		Description: descriptionHex,
	}, nil
}

func (secret Secret) Decrypt(cipherKey []byte) (Secret, error) {
	urlString, err := crypto.DecryptAES(cipherKey, secret.Url)
	if err != nil {
		return Secret{}, err
	}

	usernameString, err := crypto.DecryptAES(cipherKey, secret.Username)
	if err != nil {
		return Secret{}, err
	}

	passwordString, err := crypto.DecryptAES(cipherKey, secret.Password)
	if err != nil {
		return Secret{}, err
	}

	descriptionString, err := crypto.DecryptAES(cipherKey, secret.Description)
	if err != nil {
		return Secret{}, err
	}

	return Secret{
		Name:        secret.Name,
		Url:         urlString,
		Username:    usernameString,
		Password:    passwordString,
		Description: descriptionString,
	}, nil
}
