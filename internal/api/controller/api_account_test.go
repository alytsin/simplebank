package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alytsin/simplebank/internal/api/security/token"
	"github.com/alytsin/simplebank/internal/db"
	dbmock "github.com/alytsin/simplebank/internal/db/mock"
	"github.com/alytsin/simplebank/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func TestGetAccount(t *testing.T) {

	tokenMaker, _ := token.NewPasetoMaker("")
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name          string
		clientId      int64
		foundAccount  *db.Account
		setupAuth     func(t *testing.T, request *http.Request)
		err           error
		responseCheck func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "zero account id",
			clientId:     0,
			foundAccount: nil,
			err:          nil,
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("user"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:         "not yours account",
			clientId:     1,
			foundAccount: &db.Account{ID: 1, Owner: "owner", Currency: "USD", Balance: 0, CreatedAt: date},
			err:          nil,
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("user"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				body := recorder.Body.String()
				assert.Equal(t, `{"error":"not yours account"}`, body)
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:         "account exists",
			clientId:     1,
			foundAccount: &db.Account{ID: 1, Owner: "owner", Currency: "USD", Balance: 0, CreatedAt: date},
			err:          nil,
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("owner"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:         "common error",
			clientId:     1,
			foundAccount: nil,
			err:          errors.New("error"),
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("user"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				body := recorder.Body.String()
				assert.Equal(t, `{"error":"error"}`, body)
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:         "db no rows error",
			clientId:     1,
			foundAccount: nil,
			err:          sql.ErrNoRows,
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("user"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				body := recorder.Body.String()
				assert.Empty(t, body)
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			store := dbmock.MockTxStoreInterface{}
			store.On("GetAccount", mock.Anything, tc.clientId).
				Return(tc.foundAccount, tc.err).
				Once()

			controller := NewApiController(&store, tokenMaker, nil)

			rsp := httptest.NewRecorder()
			router := gin.New()
			router.Use(controller.AuthMiddleware())
			router.GET("/:id", controller.GetAccount)

			req, _ := http.NewRequest("GET", fmt.Sprintf("/%v", tc.clientId), nil)
			tc.setupAuth(t, req)
			router.ServeHTTP(rsp, req)
			tc.responseCheck(rsp)
		})
	}
}

func TestCreateAccount(t *testing.T) {

	tokenMaker, _ := token.NewPasetoMaker("")
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name    string
		params  *db.CreateAccountParams
		account *db.Account
		//httpStatus    int
		storeError    error
		setupAuth     func(t *testing.T, request *http.Request)
		responseCheck func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "account created",
			params:  &db.CreateAccountParams{Owner: "owner", Currency: "USD"},
			account: &db.Account{ID: 1, Owner: "owner", Currency: "USD", Balance: 0, CreatedAt: date},
			//httpStatus: http.StatusCreated,
			storeError: nil,
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("owner"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				body := recorder.Body.String()
				assert.Contains(t, body, "owner")
				assert.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name:   "account unique violation",
			params: &db.CreateAccountParams{Owner: "owner", Currency: "USD"},
			//httpStatus: http.StatusConflict,
			storeError: &pq.Error{Code: "23505"},
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("owner"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				body := recorder.Body.String()
				assert.Equal(t, `{"error":"pq: "}`, body)
				assert.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name:   "account foreign key violation",
			params: &db.CreateAccountParams{Owner: "owner", Currency: "USD"},
			//httpStatus: http.StatusConflict,
			storeError: &pq.Error{Code: "23503"},
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("owner"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				body := recorder.Body.String()
				assert.Equal(t, `{"error":"pq: "}`, body)
				assert.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name:    "bad currency",
			params:  &db.CreateAccountParams{Owner: "owner", Currency: "EEK"},
			account: nil,
			//httpStatus: http.StatusBadRequest,
			storeError: nil,
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("owner"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				//body := recorder.Body.String()
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:    "empty body",
			params:  nil,
			account: nil,
			//httpStatus: http.StatusBadRequest,
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("owner"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			storeError: nil,
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				//body := recorder.Body.String()
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:    "db error",
			params:  &db.CreateAccountParams{Owner: "owner", Currency: "USD"},
			account: nil,
			//httpStatus: http.StatusInternalServerError,
			storeError: errors.New("error"),
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("owner"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				body := recorder.Body.String()
				assert.Equal(t, `{"error":"error"}`, body)
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			var param interface{}

			store := dbmock.MockTxStoreInterface{}
			if tc.params != nil {
				param = *tc.params
			}

			store.On("CreateAccount", mock.Anything, param).
				Return(tc.account, tc.storeError).
				Once()

			controller := NewApiController(&store, tokenMaker, nil)

			rsp := httptest.NewRecorder()
			router := gin.New()
			router.Use(controller.AuthMiddleware())
			router.POST("/", controller.CreateAccount)

			b, err := json.Marshal(tc.params)
			assert.NoError(t, err)

			req, _ := http.NewRequest("POST", "/", bytes.NewReader(b))
			tc.setupAuth(t, req)
			router.ServeHTTP(rsp, req)
			tc.responseCheck(rsp)

		})
	}

}

func TestListAccounts(t *testing.T) {

	tokenMaker, _ := token.NewPasetoMaker("")
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name          string
		queryString   string
		httpStatus    int
		resultList    []*db.Account
		storeError    error
		setupAuth     func(t *testing.T, request *http.Request)
		responseCheck func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "zero page",
			queryString: "page=0",
			//httpStatus:  http.StatusBadRequest,
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("owner"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				//body := recorder.Body.String()
				//assert.Equal(t, `{"error":"error"}`, body)
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "zero page size",
			queryString: "page=1",
			httpStatus:  http.StatusBadRequest,
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("owner"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				//body := recorder.Body.String()
				//assert.Equal(t, `{"error":"error"}`, body)
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "page size too big",
			queryString: "page=1&page_size=150",
			httpStatus:  http.StatusBadRequest,
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("owner"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				//body := recorder.Body.String()
				//assert.Equal(t, `{"error":"error"}`, body)
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "ok",
			queryString: "page=1&page_size=5",
			httpStatus:  http.StatusOK,
			resultList: []*db.Account{
				{ID: 1, Owner: "owner", Currency: "USD", Balance: 0, CreatedAt: date},
				{ID: 2, Owner: "owner", Currency: "EUR", Balance: 0, CreatedAt: date},
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("owner"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				body := recorder.Body.String()
				assert.Contains(t, body, `"owner":"owner"`)
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:        "not found",
			queryString: "page=1&page_size=5",
			httpStatus:  http.StatusNotFound,
			storeError:  sql.ErrNoRows,
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("owner"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				body := recorder.Body.String()
				assert.Empty(t, body)
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:        "internal server error",
			queryString: "page=1&page_size=5",
			httpStatus:  http.StatusInternalServerError,
			storeError:  errors.New("error"),
			setupAuth: func(t *testing.T, request *http.Request) {
				tk, err := tokenMaker.CreateToken(token.NewPayload("owner"), time.Minute)
				assert.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%s %s", authorizationTypeBearer, tk))
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				body := recorder.Body.String()
				assert.Equal(t, `{"error":"error"}`, body)
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			store := dbmock.MockTxStoreInterface{}
			store.On("ListAccounts", mock.Anything, mock.Anything).
				Return(tc.resultList, tc.storeError).
				Once()

			controller := NewApiController(&store, tokenMaker, nil)

			rsp := httptest.NewRecorder()
			router := gin.New()
			router.Use(controller.AuthMiddleware())
			router.GET("/", controller.ListAccounts)

			req, _ := http.NewRequest("GET", fmt.Sprintf("/?%s", tc.queryString), nil)
			tc.setupAuth(t, req)
			router.ServeHTTP(rsp, req)
			tc.responseCheck(rsp)

		})
	}
}
