package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/alytsin/simplebank/internal/api/security"
	"github.com/alytsin/simplebank/internal/db"
	dbmock "github.com/alytsin/simplebank/internal/db/mock"
	"github.com/alytsin/simplebank/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func randomUser() *db.User {
	return &db.User{
		Username:       util.RandomString(8),
		HashedPassword: util.RandomString(8),
		FullName:       util.RandomString(8),
		Email:          util.RandomEmail(),
	}
}

func TestCreateUser(t *testing.T) {

	user := randomUser()

	cases := []struct {
		name          string
		body          *CreateUserRequest
		createParams  *db.CreateUserParams
		storeError    error
		responseCheck func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "user created",
			body: &CreateUserRequest{
				Username: user.Username,
				Password: user.HashedPassword,
				FullName: user.FullName,
				Email:    user.Email,
			},
			createParams: &db.CreateUserParams{
				Username:       user.Username,
				HashedPassword: user.HashedPassword,
				FullName:       user.FullName,
				Email:          user.Email,
			},
			responseCheck: func(recorder *httptest.ResponseRecorder) {

				data, err := io.ReadAll(recorder.Body)
				assert.NoError(t, err)

				var u CreateUserResponse
				err = json.Unmarshal(data, &u)
				assert.NoError(t, err)

				assert.Equal(t, &CreateUserResponse{
					Username:          user.Username,
					FullName:          user.FullName,
					Email:             user.Email,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
				}, &u)
				assert.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "bad request",
			body: &CreateUserRequest{},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "unique violation",
			body: &CreateUserRequest{
				Username: user.Username,
				Password: user.HashedPassword,
				FullName: user.FullName,
				Email:    user.Email,
			},
			createParams: &db.CreateUserParams{
				Username:       user.Username,
				HashedPassword: user.HashedPassword,
				FullName:       user.FullName,
				Email:          user.Email,
			},
			storeError: &pq.Error{Code: "23505"},
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name: "error",
			body: &CreateUserRequest{
				Username: user.Username,
				Password: user.HashedPassword,
				FullName: user.FullName,
				Email:    user.Email,
			},
			createParams: &db.CreateUserParams{
				Username:       user.Username,
				HashedPassword: user.HashedPassword,
				FullName:       user.FullName,
				Email:          user.Email,
			},
			storeError: errors.New("error"),
			responseCheck: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			store := dbmock.MockTxStoreInterface{}

			var param interface{}
			if tc.createParams != nil {
				param = *tc.createParams
			}

			store.On("CreateUser", mock.Anything, param).
				Return(user, tc.storeError).
				Once()

			controller := NewApiController(&store, nil, new(security.PasswordPlain))

			rsp := httptest.NewRecorder()
			router := gin.New()
			router.POST("/", controller.CreateUser)

			b, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			req, _ := http.NewRequest("POST", "/", bytes.NewReader(b))
			router.ServeHTTP(rsp, req)
			tc.responseCheck(rsp)
		})
	}
}
