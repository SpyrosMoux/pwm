package cmd

import (
	"github.com/spyrosmoux/pwm/helpers"
	"strings"
)

func Store() string {
	username := helpers.StringInput("Enter username: ")
	password := helpers.SecretInput("Enter password: ")

	return strings.TrimSpace(username) + ":" + strings.TrimSpace(password)
}
