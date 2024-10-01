package security

type PasswordPlain string

func (p *PasswordPlain) Hash() (string, error) {
	return string(*p), nil
}

func (p *PasswordPlain) Verify(hash string) bool {
	return string(*p) == hash
}
