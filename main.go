package main

import (
	"database/sql"
	"gopsql/banking/api"
	db "gopsql/banking/db/sqlc"
	"gopsql/banking/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("error in loading configuration ", err)
	}
	log.Println("config: ", config)

	conn, err := sql.Open(config.DBDrive, config.DBSource)
	if err != nil {
		log.Fatal("can not connect to DB", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can not start server", err)
	}

}
