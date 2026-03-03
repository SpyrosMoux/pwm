package secrets

import (
	"github.com/SpyrosMoux/pwm/internal/crypto"
	"github.com/SpyrosMoux/pwm/internal/store"
)

type Service struct {
	Store  store.Storer
	Crypto crypto.Crypter
}
