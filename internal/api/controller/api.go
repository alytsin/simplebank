package controller

import (
	"github.com/alytsin/simplebank/internal/api/security"
	"github.com/alytsin/simplebank/internal/api/security/token"
	"github.com/alytsin/simplebank/internal/db"
	"time"
)

type Api struct {
	Base
	tokenMaker       token.Maker
	store            db.TxStoreInterface
	passwordVerifier security.PasswordInterface
	tokenTTL         time.Duration
}

func NewApiController(
	store db.TxStoreInterface,
	tokenMaker token.Maker,
	passwordVerifier security.PasswordInterface,
) *Api {
	return &Api{
		store:            store,
		tokenMaker:       tokenMaker,
		passwordVerifier: passwordVerifier,
	}
}

func (c *Api) SetTokenTTL(ttl time.Duration) *Api {
	c.tokenTTL = ttl
	return c
}
