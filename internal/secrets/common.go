package secrets

import (
	"encoding/json"

	"github.com/SpyrosMoux/pwm/internal/models"
)

func (svc *Service) ReadSecretIntoStruct(secret string) (models.Secret, error) {
	hex, err := svc.Store.ReadFile(secret)
	if err != nil {
		return models.Secret{}, err
	}

	var jsonSecret models.Secret
	err = json.Unmarshal(hex, &jsonSecret)
	if err != nil {
		return models.Secret{}, err
	}

	err = svc.Crypto.Decrypt(&jsonSecret)
	if err != nil {
		return models.Secret{}, err
	}

	return jsonSecret, nil
}
