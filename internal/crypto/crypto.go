/*
Copyright © 2026 Spyros Mouchlianitis
*/
package crypto

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"errors"

	"github.com/SpyrosMoux/pwm/internal/models"
)

type Crypter interface {
	Encrypt(secret *models.Secret) error
	Decrypt(secret *models.Secret) error
}

type AESGCM struct {
	CipherKey []byte
}

func (aesgcm *AESGCM) Encrypt(secret *models.Secret) error {
	urlHex, err := EncryptAES(aesgcm.CipherKey, secret.Url)
	if err != nil {
		return err
	}

	usernameHex, err := EncryptAES(aesgcm.CipherKey, secret.Username)
	if err != nil {
		return err
	}

	passwordHex, err := EncryptAES(aesgcm.CipherKey, secret.Password)
	if err != nil {
		return err
	}

	descriptionHex, err := EncryptAES(aesgcm.CipherKey, secret.Description)
	if err != nil {
		return err
	}

	secret.Url = urlHex
	secret.Username = usernameHex
	secret.Password = passwordHex
	secret.Description = descriptionHex

	return nil
}

func (aesgcm *AESGCM) Decrypt(secret *models.Secret) error {
	urlString, err := DecryptAES(aesgcm.CipherKey, secret.Url)
	if err != nil {
		return err
	}

	usernameString, err := DecryptAES(aesgcm.CipherKey, secret.Username)
	if err != nil {
		return err
	}

	passwordString, err := DecryptAES(aesgcm.CipherKey, secret.Password)
	if err != nil {
		return err
	}

	descriptionString, err := DecryptAES(aesgcm.CipherKey, secret.Description)
	if err != nil {
		return err
	}

	secret.Url = urlString
	secret.Username = usernameString
	secret.Password = passwordString
	secret.Description = descriptionString

	return nil
}
func EncryptAES(cipherKey []byte, plaintext string) (string, error) {
	// create cipher
	c, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	// Apply PKCS7 padding to the plaintext
	plainBytes := []byte(plaintext)
	paddedPlaintext := PKCS7Padding(plainBytes, aes.BlockSize)

	// allocate space for ciphered data
	out := make([]byte, len(paddedPlaintext))

	// encrypt each block
	for i := 0; i < len(paddedPlaintext); i += aes.BlockSize {
		c.Encrypt(out[i:i+aes.BlockSize], paddedPlaintext[i:i+aes.BlockSize])
	}

	// return hex string
	return hex.EncodeToString(out), nil
}

// DecryptAES decrypts the hex-encoded ciphertext using the AES key and returns the plaintext.
func DecryptAES(cipherKey []byte, hexToDecode string) (string, error) {
	// Decode the hex string to bytes
	cipherBytes, err := hex.DecodeString(hexToDecode)
	if err != nil {
		return "", err
	}

	// create cipher
	c, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	// allocate space for decrypted data
	out := make([]byte, len(cipherBytes))

	// decrypt each block
	for i := 0; i < len(cipherBytes); i += aes.BlockSize {
		c.Decrypt(out[i:i+aes.BlockSize], cipherBytes[i:i+aes.BlockSize])
	}

	// Remove PKCS7 padding
	plainBytes, err := PKCS7UnPadding(out, aes.BlockSize)
	if err != nil {
		return "", err
	}

	return string(plainBytes), nil
}

// PKCS7Padding applies padding to the plaintext.
func PKCS7Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padText...)
}

func PKCS7UnPadding(plaintext []byte, blockSize int) ([]byte, error) {
	length := len(plaintext)
	if length == 0 {
		return nil, errors.New("plaintext is empty")
	}
	if length%blockSize != 0 {
		return nil, errors.New("plaintext is not a multiple of the block size")
	}
	paddingLen := int(plaintext[length-1])
	if paddingLen > blockSize || paddingLen == 0 {
		return nil, errors.New("invalid padding")
	}

	// Validate padding bytes
	for _, padByte := range plaintext[length-paddingLen:] {
		if int(padByte) != paddingLen {
			return nil, errors.New("invalid padding")
		}
	}

	return plaintext[:length-paddingLen], nil
}

func encryptSecretFile() {
	// TODO(spyrosmoux) gpg encrypt a file after creating a secret
	panic("implement me")
}

func decryptSecretFile() {
	// TODO(spyrosmoux) gpg decrypt a file before accessing it's contents
	panic("implement me")
}
