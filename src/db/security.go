package db

// Security (password, bcrypt) functions for securing the admin password.

import "golang.org/x/crypto/bcrypt"

// The bcrypt cost factor.
const (
	Cost = 14
)

// HashPassword takes a password and generates a Bcrypt hash.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), Cost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

// CheckPassword checks if a password is correct.
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}
