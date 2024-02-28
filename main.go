package main

import (
	"database/sql"
	"log"

	"github.com/adolfo-cia/go-simple-bank-course/api"
	db "github.com/adolfo-cia/go-simple-bank-course/db/sqlc"
	"github.com/adolfo-cia/go-simple-bank-course/utils"
	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalln("Cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalln("cannot start the server: ", err)
	}
}
