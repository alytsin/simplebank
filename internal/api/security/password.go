package security

import "golang.org/x/crypto/bcrypt"

type PasswordInterface interface {
	Hash() (string, error)
	Verify(hash string) bool
}

type Password string

func (p *Password) Hash() (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*p), bcrypt.DefaultCost)
	return string(bytes), err
}

func (p *Password) Verify(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(*p))
	return err == nil
}
