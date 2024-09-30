package controller

import (
	"database/sql"
	"errors"
	"github.com/alytsin/simplebank/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Api struct {
	Base
	store db.TxStoreInterface
}

func NewApiController(store db.TxStoreInterface) *Api {
	return &Api{
		store: store,
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

	account, err := c.store.CreateAccount(ctx, &arg)
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

	account, err := c.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, ErrorMessage{Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, account)
}
