package secrets

import (
	"github.com/SpyrosMoux/pwm/internal/helpers"
	"golang.design/x/clipboard"
)

func (svc *Service) CopySecret(secretName string) error {
	helpers.PrintInfo("Copying secret " + secretName)

	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		return err
	}

	secret, err := svc.ReadSecretIntoStruct(secretName)
	if err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtText, []byte(secret.Password))

	return nil
}
