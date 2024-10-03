package security

type PasswordPlain string

func (p *PasswordPlain) Hash(password string) (string, error) {
	return password, nil
}

func (p *PasswordPlain) Verify(hash string, password string) bool {
	return password == hash
}
