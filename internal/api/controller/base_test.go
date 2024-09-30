package controller

import (
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

func TestValidateUriOrSendBadRequest(t *testing.T) {

	c := Base{}

	testData := []struct {
		id     string
		result bool
		code   int
	}{
		{id: "", result: false, code: http.StatusBadRequest},
		{id: "0", result: false, code: http.StatusBadRequest},
		{id: "-1", result: false, code: http.StatusBadRequest},
		{id: "1", result: true, code: http.StatusOK},
	}

	for _, item := range testData {
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
	}

}

func TestValidateJsonRequestOrBadRequest(t *testing.T) {

	testData := []struct {
		request string
		result  bool
		code    int
	}{
		{request: `{}`, result: false, code: http.StatusBadRequest},
		{request: `{"f":""}`, result: false, code: http.StatusBadRequest},
		{request: `{"f":"1"}`, result: true, code: http.StatusOK},
	}

	for _, item := range testData {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request, _ = http.NewRequest("GET", "/", strings.NewReader(item.request))

		var req testJsonStruct
		c := Base{}

		result := c.validateJsonOrSendBadRequest(ctx, &req)
		assert.Equal(t, result, item.result)
		assert.Equal(t, w.Code, item.code)
	}
}
