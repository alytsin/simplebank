package security

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordInterface interface {
	Hash(password string) (string, error)
	Verify(hash string, password string) bool
}

type Password string

func (p *Password) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (p *Password) Verify(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
