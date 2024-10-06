package controller

import (
	"fmt"
	"github.com/alytsin/simplebank/internal/api/security/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthMiddleware(t *testing.T) {

	tokenMaker, _ := token.NewPasetoMaker("")

	cases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request)
		responseCheck func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "no header",
			setupAuth: func(t *testing.T, request *http.Request) {
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
				assert.Equal(t, `{"error":"authorization header is missing"}`, recorder.Body.String())
			},
		},
		{
			name: "invalid header",
			setupAuth: func(t *testing.T, request *http.Request) {
				request.Header.Set(authorizationHeader, "nope")
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
				assert.Equal(t, `{"error":"authorization header is invalid"}`, recorder.Body.String())
			},
		},
		{
			name: "unsupported auth type",
			setupAuth: func(t *testing.T, request *http.Request) {
				request.Header.Set(authorizationHeader, "nothing key")
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
				assert.Equal(t, `{"error":"authorization type 'nothing' is not supported"}`, recorder.Body.String())
			},
		},
		{
			name: "expired token",
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("user"), -time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("Bearer %s", tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
				assert.Equal(t, `{"error":"this token has expired"}`, recorder.Body.String())
			},
		},
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("user"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Equal(t, `OK`, recorder.Body.String())
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			controller := NewApiController(nil, tokenMaker, nil)

			rsp := httptest.NewRecorder()
			router := gin.New()
			rg := router.Group("/").Use(controller.AuthMiddleware())
			rg.GET("/auth", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "OK")
			}).Use()

			req, _ := http.NewRequest("GET", "/auth", nil)
			c.setupAuth(t, req)
			router.ServeHTTP(rsp, req)
			c.responseCheck(rsp)
		})
	}
}
