package helper

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password string) (string, error) {
	bytePassword := []byte(password)
	hashedPassword, passErr := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if passErr != nil {
		return "", errors.New(passErr.Error())
	}
	return string(hashedPassword), nil
}

func ComparePasswordHashes(passwordHash string, password string) bool {
	bytePasswordHash := []byte(passwordHash)
	bytePassword := []byte(password)
	err := bcrypt.CompareHashAndPassword(bytePasswordHash, bytePassword)
	return err == nil
}
