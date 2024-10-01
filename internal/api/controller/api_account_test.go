package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alytsin/simplebank/internal/db"
	dbmock "github.com/alytsin/simplebank/internal/db/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetAccount(t *testing.T) {

	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name         string
		clientId     int64
		foundAccount *db.Account
		httpStatus   int
		err          error
	}{
		{
			name:         "zero account id",
			clientId:     0,
			foundAccount: nil,
			httpStatus:   http.StatusBadRequest,
			err:          nil,
		},
		{
			name:         "account exists",
			clientId:     1,
			foundAccount: &db.Account{ID: 1, Owner: "owner", Currency: "USD", Balance: 0, CreatedAt: date},
			httpStatus:   http.StatusOK,
			err:          nil,
		},
		{
			name:         "common error",
			clientId:     1,
			foundAccount: nil,
			httpStatus:   http.StatusInternalServerError,
			err:          errors.New("error"),
		},
		{
			name:         "db no rows error",
			clientId:     1,
			foundAccount: nil,
			httpStatus:   http.StatusNotFound,
			err:          sql.ErrNoRows,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			store := dbmock.TxStore{}
			store.On("GetAccount", mock.Anything, tc.clientId).
				Return(tc.foundAccount, tc.err).
				Once()

			controller := NewApiController(&store, nil)

			rsp := httptest.NewRecorder()
			router := gin.New()
			router.GET("/:id", controller.GetAccount)

			req, _ := http.NewRequest("GET", fmt.Sprintf("/%v", tc.clientId), nil)
			router.ServeHTTP(rsp, req)

			if tc.foundAccount != nil {
				var result db.Account

				err := json.Unmarshal(rsp.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, &result, tc.foundAccount)
			}

			assert.Equal(t, tc.httpStatus, rsp.Code)

		})
	}
}

func TestCreateAccount(t *testing.T) {

	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name       string
		params     *db.CreateAccountParams
		account    *db.Account
		httpStatus int
		storeError error
	}{
		{
			name:       "account created",
			params:     &db.CreateAccountParams{Owner: "owner", Currency: "USD"},
			account:    &db.Account{ID: 1, Owner: "owner", Currency: "USD", Balance: 0, CreatedAt: date},
			httpStatus: http.StatusCreated,
			storeError: nil,
		},
		{
			name:       "bad currency",
			params:     &db.CreateAccountParams{Owner: "owner", Currency: "EEK"},
			account:    nil,
			httpStatus: http.StatusBadRequest,
			storeError: nil,
		},
		{
			name:       "empty body",
			params:     nil,
			account:    nil,
			httpStatus: http.StatusBadRequest,
			storeError: nil,
		},
		{
			name:       "db error",
			params:     &db.CreateAccountParams{Owner: "owner", Currency: "USD"},
			account:    nil,
			httpStatus: http.StatusInternalServerError,
			storeError: errors.New("error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			var param interface{}

			store := dbmock.TxStore{}
			if tc.params != nil {
				param = *tc.params
			}

			store.On("CreateAccount", mock.Anything, param).
				Return(tc.account, tc.storeError).
				Once()

			controller := NewApiController(&store, nil)

			rsp := httptest.NewRecorder()
			router := gin.New()
			router.POST("/", controller.CreateAccount)

			b, err := json.Marshal(tc.params)
			assert.NoError(t, err)

			req, _ := http.NewRequest("POST", "/", bytes.NewReader(b))
			router.ServeHTTP(rsp, req)

			if tc.account != nil {
				var result db.Account

				err = json.Unmarshal(rsp.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, &result, tc.account)
			}

			assert.Equal(t, tc.httpStatus, rsp.Code)

		})
	}

}
