package secrets

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SpyrosMoux/passwdgen"
	"github.com/SpyrosMoux/pwm/internal/helpers"
)

// UpdateSecret reads a secret, prompts for field updates, and rewrites it atomically.
// Preserves the original modification time so the secret maintains its index position.
func (svc *Service) UpdateSecret(secretName string) error {
	// Capture original mtime before any changes
	dstPath, err := svc.Store.GetFilePath(secretName)
	if err != nil {
		return err
	}

	stat, err := os.Stat(dstPath)
	if err != nil {
		return err
	}
	originalModTime := stat.ModTime()

	// Read and decrypt the existing secret
	secret, err := svc.ReadSecretIntoStruct(secretName)
	if err != nil {
		return err
	}

	fmt.Println("Updating secret: " + secretName)
	fmt.Println("(Press Enter to skip a field)")
	fmt.Println()

	// Prompt for optional updates
	fmt.Printf("URL [%s]: ", secret.Url)
	newUrl := helpers.StringInput("")
	if newUrl != "" {
		secret.Url = newUrl
	}

	fmt.Printf("Username [%s]: ", secret.Username)
	newUsername := helpers.StringInput("")
	if newUsername != "" {
		secret.Username = newUsername
	}

	fmt.Printf("Password [%s]: ", secret.Password)
	newPassword := helpers.OptionalSecretInput("('a' to autogenerate, Enter to skip): ")
	if newPassword == "a" {
		options := passwdgen.NewRandomStringOptions()
		generatedPassword := passwdgen.RandomStringNumbersSymbols(&options)
		fmt.Printf("Generated password: %s\n", generatedPassword)
		secret.Password = generatedPassword
	} else if newPassword != "" {
		secret.Password = newPassword
	}

	fmt.Printf("Description [%s]: ", secret.Description)
	newDescription := helpers.StringInput("")
	if newDescription != "" {
		secret.Description = newDescription
	}

	// Re-encrypt the updated secret
	err = svc.Crypto.Encrypt(&secret)
	if err != nil {
		return err
	}

	// Marshal to JSON
	jsonSecret, err := json.Marshal(secret)
	if err != nil {
		return err
	}

	// Atomically write: write to temp file, then replace
	tempPath := string(dstPath) + ".tmp"

	tempFile, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	_, err = tempFile.Write(jsonSecret)
	if err != nil {
		os.Remove(tempPath)
		return err
	}

	tempFile.Close()

	// Atomic replace
	err = os.Rename(tempPath, string(dstPath))
	if err != nil {
		os.Remove(tempPath)
		return err
	}

	// Restore original modification time to preserve index position
	err = os.Chtimes(string(dstPath), originalModTime, originalModTime)
	if err != nil {
		// Log but don't fail: mtime restoration is non-critical
		fmt.Fprintf(os.Stderr, "Warning: could not restore modification time: %v\n", err)
	}

	return nil
}
