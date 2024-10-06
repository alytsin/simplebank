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

func (s *Server) Run(listen string) error {

	router := gin.Default()
	//router.Use(gin.Recovery())
	_ = router.SetTrustedProxies(nil)

	router.POST("/users", s.controller.CreateUser)
	router.POST("/users/login", s.controller.LoginUser)

	authGroup := router.Group("/").Use(s.controller.AuthMiddleware())

	authGroup.POST("/transfers", s.controller.CreateTransfer)
	authGroup.GET("/accounts", s.controller.ListAccounts)
	authGroup.GET("/account/:id", s.controller.GetAccount)
	authGroup.POST("/accounts", s.controller.CreateAccount)

	return router.Run(listen)
}
