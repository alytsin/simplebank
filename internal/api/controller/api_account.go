package controller

import (
	"errors"
	"fmt"
	"github.com/alytsin/simplebank/internal/api/security/token"
	"github.com/alytsin/simplebank/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Api) CreateAccount(ctx *gin.Context) {
	var req CreateAccountRequest

	if !c.validateJsonOrSendBadRequest(ctx, &req) {
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    payload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := c.store.CreateAccount(ctx, arg)
	if err != nil {
		if errors.Is(db.TranslateError(err), db.ErrUniqueViolation) {
			ctx.JSON(http.StatusConflict, ErrorMessage{Error: err})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err})
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
		if errors.Is(db.TranslateError(err), db.ErrNoRows) {
			ctx.String(http.StatusNotFound, "")
			return
		}

		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err})
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != payload.Username {
		ctx.JSON(http.StatusUnauthorized, ErrorMessage{Error: fmt.Errorf("not yours account")})
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (c *Api) ListAccounts(ctx *gin.Context) {
	var req ListAccountsRequest

	if !c.validateQueryOrSendBadRequest(ctx, &req) {
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	accounts, err := c.store.ListAccounts(ctx, db.ListAccountsParams{
		Owner:  payload.Username,
		Limit:  req.PageSize,
		Offset: (req.Page - 1) * req.PageSize,
	})
	if err != nil {
		if errors.Is(db.TranslateError(err), db.ErrNoRows) {
			ctx.String(http.StatusNotFound, "")
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err})
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
