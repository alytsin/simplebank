package controller

import (
	"github.com/alytsin/simplebank/internal/api/security"
	"github.com/alytsin/simplebank/internal/api/security/token"
	"github.com/alytsin/simplebank/internal/db"
)

type Api struct {
	Base
	tokenMaker       token.Maker
	store            db.TxStoreInterface
	passwordVerifier security.PasswordInterface
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
