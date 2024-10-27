package controller

import (
	"errors"
	"github.com/alytsin/simplebank/internal/api/security/token"
	"github.com/alytsin/simplebank/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (c *Api) RenewAccessToken(ctx *gin.Context) {
	var req RenewAccessTokenRequest

	if !c.validateJsonOrSendBadRequest(ctx, &req) {
		return
	}

	refreshPayload, err := c.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorMessage{Error: err})
		return
	}

	session, err := c.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if errors.Is(db.TranslateError(err), db.ErrNoRows) {
			ctx.String(http.StatusNotFound, "")
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err})
		return
	}

	if session.Username != refreshPayload.Username {
		ctx.JSON(http.StatusUnauthorized, ErrorMessage{Error: errors.New("session username mismatch")})
		return
	}

	accessToken, err := c.tokenMaker.CreateToken(token.NewPayload(session.Username), c.config.AccessTokenTTL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err})
		return
	}

	ctx.JSON(http.StatusOK, RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiredAt: time.Now().Add(c.config.AccessTokenTTL),
	})

}

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

	accessDeadline := time.Now().Add(c.config.AccessTokenTTL)
	accessToken, err := c.tokenMaker.CreateToken(token.NewPayload(user.Username), c.config.AccessTokenTTL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err})
		return
	}

	refreshDeadline := time.Now().Add(c.config.RefreshTokenTTL)
	refreshTokenPayload := token.NewPayload(user.Username)
	refreshToken, err := c.tokenMaker.CreateToken(refreshTokenPayload, c.config.RefreshTokenTTL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorMessage{Error: err})
		return
	}

	_, err = c.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshDeadline,
	})

	ctx.JSON(http.StatusOK, LoginUserResponse{
		AccessToken:           accessToken,
		AccessTokenExpiredAt:  accessDeadline,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: refreshDeadline,
	})

}
