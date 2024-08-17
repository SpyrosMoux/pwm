package models

import "fmt"

type Secret struct {
	Name        string
	Url         string
	Username    string
	Password    string
	Description string
}

func (secret Secret) ToString() string {
	return fmt.Sprintf("Name: %s\n"+"Url: %s\n"+"Username: %s\n"+"Password: %s\n"+"Description: %s\n",
		secret.Name, secret.Url, secret.Username, secret.Password,
		secret.Description)
}
