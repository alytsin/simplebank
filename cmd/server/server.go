package main

import (
	"database/sql"
	"errors"
	"github.com/alytsin/simplebank/internal"
	"github.com/alytsin/simplebank/internal/api"
	"github.com/alytsin/simplebank/internal/api/controller"
	"github.com/alytsin/simplebank/internal/api/security"
	"github.com/alytsin/simplebank/internal/api/security/token"
	"github.com/alytsin/simplebank/internal/db"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func main() {

	log.Println("Starting server...")

	config, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	migration, err := migrate.New(config.MigrationUrl, config.DbSource)
	if err != nil {
		log.Fatal("cannot init migration:", err)
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("failed to run migrate up:", err)
	}

	database, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	tokenMaker, err := token.NewPasetoMaker(config.AccessTokenPrivateKey)
	if err != nil {
		log.Fatal("unable to create token maker:", err)
	}

	cntrlr := controller.NewApiController(
		db.NewTxStore(database),
		tokenMaker,
		new(security.Password),
	).SetTokenTTL(config.AccessTokenTTL)

	server := api.NewServer(cntrlr)

	log.Println(server.Run(config.ServerAddress))
}
