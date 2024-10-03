package controller

import (
	"github.com/alytsin/simplebank/internal/api/security"
	"github.com/alytsin/simplebank/internal/db"
)

type Api struct {
	Base
	store            db.TxStoreInterface
	passwordVerifier security.PasswordInterface
}

func NewApiController(
	store db.TxStoreInterface,
	passwordVerifier security.PasswordInterface,
) *Api {
	return &Api{
		store:            store,
		passwordVerifier: passwordVerifier,
	}
}

//func (c *Api) CreateTransfer(ctx *gin.Context) {
//
//}
