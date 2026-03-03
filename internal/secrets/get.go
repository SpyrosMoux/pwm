package secrets

import (
	"encoding/json"

	"github.com/SpyrosMoux/pwm/internal/models"
)

// GetSecret reads a given secret, decrypts it and returns it
// Currently does not support subdirectories
// TODO(spyrosmoux) should be able to get secrets in subdirectories
func (svc *Service) GetSecret(secret string) (decryptedSecret string, err error) {
	hex, err := svc.Store.ReadFile(secret)
	if err != nil {
		return
	}

	var jsonSecret models.Secret
	err = json.Unmarshal(hex, &jsonSecret)
	if err != nil {
		return
	}

	err = svc.Crypto.Decrypt(&jsonSecret)
	if err != nil {
		return "", err
	}

	return jsonSecret.String(), nil
}
