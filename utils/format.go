package utils

import "net/mail"

func IsValidEmailAddress(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
