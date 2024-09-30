package controller

import (
	val "github.com/alytsin/simplebank/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("currency", val.CurrencyValidator)
	}

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
