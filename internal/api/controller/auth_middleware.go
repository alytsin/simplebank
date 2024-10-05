package controller

import (
	"fmt"
	"github.com/alytsin/simplebank/internal/api/security/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader     = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "auth"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader(authorizationHeader)
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorMessage{Error: fmt.Errorf("authorization header is missing")})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorMessage{Error: fmt.Errorf("authorization header is invalid")})
			return
		}

		if authorizationTypeBearer != strings.ToLower(fields[0]) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorMessage{
				Error: fmt.Errorf("authorization type '%s' is not supported", fields[0]),
			})
			return
		}

		payload, err := tokenMaker.VerifyToken(fields[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorMessage{Error: err})
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
