package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// GeneratePasswordHash Generates client's password hash
func GeneratePasswordHash(password string) (string, error) {
	pwd := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash Check client's password in login step
func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		//logger.Errorf(err.Error())
		return true
	}
	return false
}
