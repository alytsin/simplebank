package controller

import (
	"github.com/alytsin/simplebank/internal"
	"github.com/alytsin/simplebank/internal/api/security"
	"github.com/alytsin/simplebank/internal/api/security/token"
	"github.com/alytsin/simplebank/internal/db"
)

type Api struct {
	Base
	tokenMaker       token.Maker
	store            db.TxStoreInterface
	passwordVerifier security.PasswordInterface
	config           *internal.Config
}

func NewApiController(
	store db.TxStoreInterface,
	tokenMaker token.Maker,
	passwordVerifier security.PasswordInterface,
	config *internal.Config,
) *Api {
	return &Api{
		store:            store,
		config:           config,
		tokenMaker:       tokenMaker,
		passwordVerifier: passwordVerifier,
	}
}
