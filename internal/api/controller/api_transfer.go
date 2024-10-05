package controller

import (
	"errors"
	"fmt"
	"github.com/alytsin/simplebank/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Api) CreateTransfer(ctx *gin.Context) {

	var err error
	var req TransferRequest

	if !c.validateJsonOrSendBadRequest(ctx, &req) {
		return
	}

	_, valid := c.validateAccountForTransfer(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	_, valid = c.validateAccountForTransfer(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	transfer, err := c.store.Transfer(ctx, db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err})
		return
	}

	ctx.JSON(http.StatusCreated, transfer)

}

func (c *Api) validateAccountForTransfer(ctx *gin.Context, accountId int64, currency string) (*db.Account, bool) {

	account, err := c.store.ValidAccountIdWithCurrency(ctx, db.ValidAccountIdWithCurrencyParams{
		ID:       accountId,
		Currency: currency,
	})
	if err != nil {
		if errors.Is(db.TranslateError(err), db.ErrNoRows) {
			ctx.JSON(http.StatusUnprocessableEntity, ErrorMessage{
				Error: fmt.Errorf(
					"account id '%d' with currency '%s' does not exist",
					accountId,
					currency,
				),
			})
			return nil, false
		}
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err})
		return nil, false
	}

	return account, true
}
