package main

import (
	"database/sql"
	"github.com/alytsin/simplebank/internal"
	"github.com/alytsin/simplebank/internal/api"
	"github.com/alytsin/simplebank/internal/api/controller"
	"github.com/alytsin/simplebank/internal/api/security"
	"github.com/alytsin/simplebank/internal/db"
	"log"
)

func main() {

	config, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	database, err := sql.Open("postgres", config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	server := api.NewServer(controller.NewApiController(
		db.NewTxStore(database),
		new(security.Password),
	))
	log.Println(server.Run())

}
