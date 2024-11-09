package utils

import (
	"beetle-quest/pkg/models"
	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
)

func GenerateRandomSalt(byte_size int) ([]byte, error) {
	bytes := make([]byte, byte_size)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}

func CompareHashPassword(password, hash []byte) error {
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		return models.ErrInvalidUsernameOrPass
	}
	return nil
}

func GenerateHashFromPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, models.ErrCheckingPassword
	}
	return hashedPassword, nil
}
