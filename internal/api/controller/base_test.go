package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testJsonStruct struct {
	F string `json:"f" binding:"required,len=1"`
}

type testUriStruct struct {
	S int `uri:"clientId" binding:"required,gt=0"`
}

type testQueryStruct struct {
	Q int `form:"q" binding:"required,gt=0"`
}

func TestValidateUriOrSendBadRequest(t *testing.T) {

	testData := []struct {
		name   string
		id     string
		result bool
		code   int
	}{
		{name: "empty id", id: "", result: false, code: http.StatusBadRequest},
		{name: "zero id", id: "0", result: false, code: http.StatusBadRequest},
		{name: "negative id", id: "-1", result: false, code: http.StatusBadRequest},
		{name: "ok", id: "1", result: true, code: http.StatusOK},
	}

	for _, item := range testData {
		t.Run(item.name, func(t *testing.T) {

			c := Base{}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			params := gin.Params{}
			params = append(params, gin.Param{
				Key:   "clientId",
				Value: item.id,
			})

			ctx.Params = params

			var req testUriStruct
			result := c.validateUriOrSendBadRequest(ctx, &req)

			assert.Equal(t, result, item.result)
			assert.Equal(t, w.Code, item.code)
		})
	}

}

func TestValidateQueryOrSendBadRequest(t *testing.T) {
	cases := []struct {
		name    string
		request string
		result  bool
		code    int
	}{
		{name: "zero all", request: "", result: false, code: http.StatusBadRequest},
		{name: "empty value", request: "q=", result: false, code: http.StatusBadRequest},
		{name: "string value", request: "q=a", result: false, code: http.StatusBadRequest},
		{name: "ok", request: "q=1", result: true, code: http.StatusOK},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Request, _ = http.NewRequest("GET", fmt.Sprintf("/?%s", tt.request), nil)

			var req testQueryStruct
			c := Base{}

			result := c.validateQueryOrSendBadRequest(ctx, &req)
			assert.Equal(t, result, tt.result)
			assert.Equal(t, w.Code, tt.code)
		})
	}
}

func TestValidateJsonRequestOrBadRequest(t *testing.T) {

	testData := []struct {
		name    string
		request string
		result  bool
		code    int
	}{
		{name: "empty object", request: `{}`, result: false, code: http.StatusBadRequest},
		{name: "empty value", request: `{"f":""}`, result: false, code: http.StatusBadRequest},
		{name: "ok", request: `{"f":"1"}`, result: true, code: http.StatusOK},
	}

	for _, item := range testData {
		t.Run(item.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Request, _ = http.NewRequest("POST", "/", strings.NewReader(item.request))

			var req testJsonStruct
			c := Base{}

			result := c.validateJsonOrSendBadRequest(ctx, &req)
			assert.Equal(t, result, item.result)
			assert.Equal(t, w.Code, item.code)
		})
	}
}
