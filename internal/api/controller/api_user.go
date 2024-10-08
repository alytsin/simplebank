package controller

import (
	"errors"
	"github.com/alytsin/simplebank/internal/api/security/token"
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

	ctx.JSON(http.StatusCreated, CreateUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	})
}

func (c *Api) LoginUser(ctx *gin.Context) {
	var req LoginUserRequest

	if !c.validateJsonOrSendBadRequest(ctx, &req) {
		return
	}

	user, err := c.store.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(db.TranslateError(err), db.ErrNoRows) {
			ctx.String(http.StatusNotFound, "")
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err})
		return
	}

	if ok := c.passwordVerifier.Verify(user.HashedPassword, req.Password); !ok {
		ctx.JSON(http.StatusUnauthorized, ErrorMessage{Error: errors.New("invalid password")})
		return
	}

	newToken, err := c.tokenMaker.CreateToken(token.NewPayload(user.Username), c.tokenTTL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err})
		return
	}

	ctx.JSON(http.StatusOK, LoginUserResponse{Token: newToken})

}
