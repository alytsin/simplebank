package controller

import (
	"database/sql"
	"errors"
	"github.com/alytsin/simplebank/internal/api/security"
	"github.com/alytsin/simplebank/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
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

func (c *Api) CreateAccount(ctx *gin.Context) {
	var req CreateAccountRequest

	if !c.validateJsonOrSendBadRequest(ctx, &req) {
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := c.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

func (c *Api) GetAccount(ctx *gin.Context) {
	var req GetAccountRequest

	if !c.validateUriOrSendBadRequest(ctx, &req) {
		return
	}

	account, err := c.store.GetAccount(ctx, req.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, "")
			return
		}

		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (c *Api) ListAccounts(ctx *gin.Context) {
	var req ListAccountsRequest

	if !c.validateQueryOrSendBadRequest(ctx, &req) {
		return
	}

	accounts, err := c.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.Page - 1) * req.PageSize,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, "")
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

//func (c *Api) CreateTransfer(ctx *gin.Context) {
//
//}
