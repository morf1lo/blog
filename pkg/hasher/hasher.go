package hasher

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func NewHash(length uint8) (string, error) {
	h := make([]byte, length)
	_, err := rand.Read(h)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h), nil
}

func HashPassword(bytes []byte) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(bytes, 10)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func Hash(id int, length uint8) (string, error) {
	idBytes := []byte(strconv.Itoa(id))
	b := make([]byte, length)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	combined := append(idBytes, b...)

	return hex.EncodeToString(combined), nil
}
