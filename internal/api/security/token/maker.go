package token

import (
	"aidanwoods.dev/go-paseto"
	"time"
)

const (
	tokenDataKey = "data"
)

type Maker interface {
	CreateToken(payload *Payload, duration time.Duration) (string, error)
	VerifyToken(string) (*Payload, error)
}

type PasetoMaker struct {
	privateKey paseto.V4AsymmetricSecretKey
	publicKey  paseto.V4AsymmetricPublicKey
	dataKey    string
}

func NewPasetoMaker(hexPrivateKey string) (*PasetoMaker, error) {

	var err error
	var privateKey paseto.V4AsymmetricSecretKey

	if hexPrivateKey != "" {
		if privateKey, err = paseto.NewV4AsymmetricSecretKeyFromHex(hexPrivateKey); err != nil {
			return nil, err
		}
	} else {
		privateKey = paseto.NewV4AsymmetricSecretKey()
	}

	return &PasetoMaker{
		dataKey:    tokenDataKey,
		privateKey: privateKey,
		publicKey:  privateKey.Public(),
	}, nil
}

func (p *PasetoMaker) CreateToken(payload *Payload, duration time.Duration) (string, error) {

	now := time.Now()

	token := paseto.NewToken()
	token.SetIssuedAt(now)
	token.SetExpiration(now.Add(duration))
	token.SetNotBefore(now)

	if err := token.Set(p.dataKey, payload); err != nil {
		return "", err
	}

	signed := token.V4Sign(p.privateKey, nil)
	return signed, nil
}

func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {

	parser := paseto.NewParser()
	parser.AddRule(paseto.ValidAt(time.Now()))
	parser.AddRule(paseto.NotBeforeNbf())
	parser.AddRule(paseto.NotExpired())

	parsedToken, err := parser.ParseV4Public(p.publicKey, token, nil)
	if err != nil {
		return nil, err
	}

	var payload *Payload
	if err := parsedToken.Get(p.dataKey, &payload); err != nil {
		return nil, err
	}

	return payload, nil
}
