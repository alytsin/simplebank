package api

import (
	"github.com/alytsin/simplebank/internal/api/controller"
	val "github.com/alytsin/simplebank/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	controller *controller.Api
}

func NewServer(controller *controller.Api) *Server {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("currency", val.CurrencyValidator)
	}

	return &Server{
		controller: controller,
	}
}

func (s *Server) Run() error {

	router := gin.Default()
	router.Use(gin.Recovery())
	_ = router.SetTrustedProxies(nil)

	router.POST("/users", s.controller.CreateUser)

	router.GET("/accounts", s.controller.ListAccounts)
	router.GET("/account/:id", s.controller.GetAccount)
	router.POST("/accounts", s.controller.CreateAccount)

	return router.Run(":8080")
}
