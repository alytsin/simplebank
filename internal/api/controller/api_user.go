package controller

import (
	"errors"
	"github.com/alytsin/simplebank/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Api) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest

	if !c.validateJsonOrSendBadRequest(ctx, &req) {
		return
	}

	hash, err := c.passwordVerifier.Hash(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err})
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hash,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := c.store.CreateUser(ctx, arg)
	if err != nil {
		if errors.Is(db.TranslateError(err), db.ErrUniqueViolation) {
			ctx.JSON(http.StatusConflict, ErrorMessage{Error: err})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}
