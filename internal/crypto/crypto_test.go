package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptAESDecryptAESRoundTrip(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plaintext := "hello-secret"

	ciphertext, err := EncryptAES(key, plaintext)
	if err != nil {
		t.Fatalf("EncryptAES() error = %v", err)
	}
	if ciphertext == plaintext {
		t.Fatalf("EncryptAES() returned plaintext unchanged")
	}

	decrypted, err := DecryptAES(key, ciphertext)
	if err != nil {
		t.Fatalf("DecryptAES() error = %v", err)
	}
	if decrypted != plaintext {
		t.Fatalf("DecryptAES() = %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptAESInvalidKeyLength(t *testing.T) {
	_, err := EncryptAES([]byte("short"), "abc")
	if err == nil {
		t.Fatalf("EncryptAES() expected error for invalid key length")
	}
}

func TestDecryptAESInvalidHex(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	_, err := DecryptAES(key, "not-hex")
	if err == nil {
		t.Fatalf("DecryptAES() expected error for invalid hex")
	}
}

func TestPKCS7PaddingAddsFullBlockWhenAligned(t *testing.T) {
	input := bytes.Repeat([]byte{'x'}, 16)
	out := PKCS7Padding(input, 16)
	if len(out) != 32 {
		t.Fatalf("PKCS7Padding() len = %d, want 32", len(out))
	}
	for i := 16; i < 32; i++ {
		if out[i] != 16 {
			t.Fatalf("PKCS7Padding() padding byte = %d, want 16", out[i])
		}
	}
}

func TestPKCS7UnPaddingErrors(t *testing.T) {
	t.Run("empty plaintext", func(t *testing.T) {
		_, err := PKCS7UnPadding(nil, 16)
		if err == nil {
			t.Fatalf("PKCS7UnPadding() expected error for empty input")
		}
	})

	t.Run("not multiple of block size", func(t *testing.T) {
		_, err := PKCS7UnPadding([]byte{1, 2, 3}, 16)
		if err == nil {
			t.Fatalf("PKCS7UnPadding() expected error for invalid block multiple")
		}
	})

	t.Run("invalid trailing padding bytes", func(t *testing.T) {
		bad := append(bytes.Repeat([]byte{'a'}, 12), []byte{4, 4, 4, 3}...)
		_, err := PKCS7UnPadding(bad, 16)
		if err == nil {
			t.Fatalf("PKCS7UnPadding() expected error for invalid padding bytes")
		}
	})
}
