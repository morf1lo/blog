package auth

import "golang.org/x/crypto/bcrypt"

func VerifyPassword(hashed []byte, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashed, password)
	return err == nil
}
