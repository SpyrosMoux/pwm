package models

import "testing"

func TestSecretString(t *testing.T) {
	secret := Secret{
		Name:        "github",
		Url:         "https://github.com",
		Username:    "user",
		Password:    "pass",
		Description: "desc",
	}

	got := secret.String()
	want := "Name: github\nUrl: https://github.com\nUsername: user\nPassword: pass\nDescription: desc\n"

	if got != want {
		t.Fatalf("Secret.String() = %q, want %q", got, want)
	}
}
