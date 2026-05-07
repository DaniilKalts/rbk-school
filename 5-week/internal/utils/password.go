package utils

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

const saltSize = 16

func GenerateSalt() (string, error) {
	buf := make([]byte, saltSize)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}

func HashPassword(plainText string, salt string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText+salt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func VerifyPassword(plainText string, salt string, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plainText+salt))
	return err == nil
}
