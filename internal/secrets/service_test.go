package secrets

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/SpyrosMoux/pwm/internal/models"
)

type fakeStore struct {
	readData       []byte
	readErr        error
	checkExists    bool
	checkExistsErr error
	removeErr      error
}

func (f *fakeStore) StoreFile(fileName string, hex []byte) (string, error) {
	return "", nil
}

func (f *fakeStore) GetFilesSortedByModTime(root string) ([]string, error) {
	return nil, nil
}

func (f *fakeStore) ReadFile(fileName string) ([]byte, error) {
	return f.readData, f.readErr
}

func (f *fakeStore) CheckFileExists(fileName string) (bool, error) {
	return f.checkExists, f.checkExistsErr
}

func (f *fakeStore) RemoveFile(fileName string) error {
	return f.removeErr
}

func (f *fakeStore) GetFilePath(fileName string) (string, error) {
	return "", nil
}

type fakeCrypter struct {
	decryptErr error
}

func (f *fakeCrypter) Encrypt(secret *models.Secret) error {
	return nil
}

func (f *fakeCrypter) Decrypt(secret *models.Secret) error {
	return f.decryptErr
}

func TestReadSecretIntoStruct(t *testing.T) {
	secret := models.Secret{
		Name:        "github",
		Url:         "https://github.com",
		Username:    "user",
		Password:    "pass",
		Description: "desc",
	}
	body, err := json.Marshal(secret)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	t.Run("success", func(t *testing.T) {
		svc := &Service{
			Store:  &fakeStore{readData: body},
			Crypto: &fakeCrypter{},
		}

		got, err := svc.ReadSecretIntoStruct("github")
		if err != nil {
			t.Fatalf("ReadSecretIntoStruct() error = %v", err)
		}
		if got != secret {
			t.Fatalf("ReadSecretIntoStruct() = %#v, want %#v", got, secret)
		}
	})

	t.Run("store read error", func(t *testing.T) {
		readErr := errors.New("read failed")
		svc := &Service{
			Store:  &fakeStore{readErr: readErr},
			Crypto: &fakeCrypter{},
		}

		_, err := svc.ReadSecretIntoStruct("github")
		if !errors.Is(err, readErr) {
			t.Fatalf("ReadSecretIntoStruct() error = %v, want %v", err, readErr)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &Service{
			Store:  &fakeStore{readData: []byte("not-json")},
			Crypto: &fakeCrypter{},
		}

		_, err := svc.ReadSecretIntoStruct("github")
		if err == nil {
			t.Fatalf("ReadSecretIntoStruct() expected JSON unmarshal error")
		}
	})

	t.Run("decrypt error", func(t *testing.T) {
		decryptErr := errors.New("decrypt failed")
		svc := &Service{
			Store:  &fakeStore{readData: body},
			Crypto: &fakeCrypter{decryptErr: decryptErr},
		}

		_, err := svc.ReadSecretIntoStruct("github")
		if !errors.Is(err, decryptErr) {
			t.Fatalf("ReadSecretIntoStruct() error = %v, want %v", err, decryptErr)
		}
	})
}

func TestGetSecret(t *testing.T) {
	secret := models.Secret{
		Name:        "github",
		Url:         "https://github.com",
		Username:    "user",
		Password:    "pass",
		Description: "desc",
	}
	body, err := json.Marshal(secret)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	t.Run("success", func(t *testing.T) {
		svc := &Service{
			Store:  &fakeStore{readData: body},
			Crypto: &fakeCrypter{},
		}

		got, err := svc.GetSecret("github")
		if err != nil {
			t.Fatalf("GetSecret() error = %v", err)
		}
		if got != secret.String() {
			t.Fatalf("GetSecret() = %q, want %q", got, secret.String())
		}
	})

	t.Run("read error", func(t *testing.T) {
		readErr := errors.New("read failed")
		svc := &Service{
			Store:  &fakeStore{readErr: readErr},
			Crypto: &fakeCrypter{},
		}

		_, err := svc.GetSecret("github")
		if !errors.Is(err, readErr) {
			t.Fatalf("GetSecret() error = %v, want %v", err, readErr)
		}
	})

	t.Run("decrypt error", func(t *testing.T) {
		decryptErr := errors.New("decrypt failed")
		svc := &Service{
			Store:  &fakeStore{readData: body},
			Crypto: &fakeCrypter{decryptErr: decryptErr},
		}

		_, err := svc.GetSecret("github")
		if !errors.Is(err, decryptErr) {
			t.Fatalf("GetSecret() error = %v, want %v", err, decryptErr)
		}
	})
}

func TestRemoveSecret(t *testing.T) {
	t.Run("check exists error", func(t *testing.T) {
		checkErr := errors.New("check failed")
		svc := &Service{
			Store:  &fakeStore{checkExistsErr: checkErr},
			Crypto: &fakeCrypter{},
		}

		err := svc.RemoveSecret("github")
		if !errors.Is(err, checkErr) {
			t.Fatalf("RemoveSecret() error = %v, want %v", err, checkErr)
		}
	})

	t.Run("missing file currently returns nil", func(t *testing.T) {
		svc := &Service{
			Store:  &fakeStore{checkExists: false},
			Crypto: &fakeCrypter{},
		}

		err := svc.RemoveSecret("github")
		if err != nil {
			t.Fatalf("RemoveSecret() error = %v, want nil", err)
		}
	})

	t.Run("remove error", func(t *testing.T) {
		removeErr := errors.New("remove failed")
		svc := &Service{
			Store: &fakeStore{
				checkExists: true,
				removeErr:   removeErr,
			},
			Crypto: &fakeCrypter{},
		}

		err := svc.RemoveSecret("github")
		if !errors.Is(err, removeErr) {
			t.Fatalf("RemoveSecret() error = %v, want %v", err, removeErr)
		}
	})

	t.Run("success", func(t *testing.T) {
		svc := &Service{
			Store:  &fakeStore{checkExists: true},
			Crypto: &fakeCrypter{},
		}

		err := svc.RemoveSecret("github")
		if err != nil {
			t.Fatalf("RemoveSecret() error = %v", err)
		}
	})
}
