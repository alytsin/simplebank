package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/alytsin/simplebank/internal/api/security"
	"github.com/alytsin/simplebank/internal/db"
	dbmock "github.com/alytsin/simplebank/internal/db/mock"
	"github.com/alytsin/simplebank/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestValidateAccountForTransfer(t *testing.T) {

	EUR := validator.CurrencyEUR.String()
	USD := validator.CurrencyUSD.String()

	user1 := randomUser()
	user2 := randomUser()
	user3 := randomUser()

	acc1 := randomAccount(user1.Username)
	acc2 := randomAccount(user2.Username)
	acc3 := randomAccount(user3.Username)

	acc1.Currency = USD
	acc2.Currency = USD
	acc3.Currency = EUR

	req := &TransferRequest{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Currency:      USD,
		Amount:        5,
	}

	cases := []struct {
		name          string
		body          *TransferRequest
		stubs         func(s *dbmock.MockTxStoreInterface)
		responseCheck func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "transfer error",
			body: req,
			stubs: func(s *dbmock.MockTxStoreInterface) {

				s.On("ValidAccountIdWithCurrency", mock.Anything, db.ValidAccountIdWithCurrencyParams{
					ID:       acc1.ID,
					Currency: USD,
				}).Return(&db.Account{ID: acc1.ID, Currency: USD}, nil).Once()

				s.On("ValidAccountIdWithCurrency", mock.Anything, db.ValidAccountIdWithCurrencyParams{
					ID:       acc2.ID,
					Currency: USD,
				}).Return(&db.Account{ID: acc2.ID, Currency: USD}, nil).Once()

				s.On("Transfer", mock.Anything, mock.Anything).
					Return(&db.TransferTxResult{}, errors.New("omg")).
					Once()
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				body := recorder.Body.String()
				assert.Contains(t, body, "omg")
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "second account error",
			body: req,
			stubs: func(s *dbmock.MockTxStoreInterface) {
				s.On("ValidAccountIdWithCurrency", mock.Anything, db.ValidAccountIdWithCurrencyParams{
					ID:       acc1.ID,
					Currency: USD,
				}).Return(&db.Account{ID: acc1.ID, Currency: USD}, nil).Once()

				s.On("ValidAccountIdWithCurrency", mock.Anything, db.ValidAccountIdWithCurrencyParams{
					ID:       acc2.ID,
					Currency: USD,
				}).Return(&db.Account{}, errors.New("xyz")).Once()
				s.On("Transfer", mock.Anything, mock.Anything).Times(0)
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				body := recorder.Body.String()
				assert.Contains(t, body, "xyz")
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "first account not found",
			body: req,
			stubs: func(s *dbmock.MockTxStoreInterface) {
				s.On("ValidAccountIdWithCurrency", mock.Anything, db.ValidAccountIdWithCurrencyParams{
					ID:       acc1.ID,
					Currency: USD,
				}).Return(&db.Account{}, sql.ErrNoRows).Once()
				s.On("Transfer", mock.Anything, mock.Anything).Times(0)
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				body := recorder.Body.String()
				assert.Contains(t, body, "does not exist")
				assert.Contains(t, body, strconv.FormatInt(acc1.ID, 10))
				assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "OK",
			body: req,
			stubs: func(s *dbmock.MockTxStoreInterface) {

				s.On("ValidAccountIdWithCurrency", mock.Anything, db.ValidAccountIdWithCurrencyParams{
					ID:       acc1.ID,
					Currency: USD,
				}).Return(&db.Account{ID: acc1.ID, Currency: USD}, nil).Once()

				s.On("ValidAccountIdWithCurrency", mock.Anything, db.ValidAccountIdWithCurrencyParams{
					ID:       acc2.ID,
					Currency: USD,
				}).Return(&db.Account{ID: acc2.ID, Currency: USD}, nil).Once()

				s.On("Transfer", mock.Anything, db.TransferTxParams{
					FromAccountID: acc1.ID,
					ToAccountID:   acc2.ID,
					Amount:        5,
				}).Return(&db.TransferTxResult{FromAccount: &acc1, ToAccount: &acc2}, nil).Once()
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				body := recorder.Body.String()
				assert.Contains(t, body, acc1.Owner)
				assert.Contains(t, body, acc2.Owner)
				assert.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "empty body",
			body: nil,
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "negative amount",
			body: &TransferRequest{
				FromAccountID: 1,
				ToAccountID:   2,
				Currency:      USD,
				Amount:        -100,
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			store := dbmock.MockTxStoreInterface{}
			if c.stubs != nil {
				c.stubs(&store)
			}

			controller := NewApiController(&store, new(security.PasswordPlain))

			rsp := httptest.NewRecorder()
			router := gin.New()
			router.POST("/", controller.CreateTransfer)

			var body io.Reader

			if c.body != nil {
				b, err := json.Marshal(c.body)
				assert.NoError(t, err)
				body = bytes.NewReader(b)
			}

			req, _ := http.NewRequest("POST", "/", body)
			router.ServeHTTP(rsp, req)
			c.responseCheck(rsp)

		})
	}

}
