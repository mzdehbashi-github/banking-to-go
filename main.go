package main

import (
	"database/sql"
	"gopsql/banking/api"
	db "gopsql/banking/db/sqlc"
	"gopsql/banking/token"
	"gopsql/banking/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config := util.LoadConfig()
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can not connect to DB", err)
	}

	store := db.NewStore(conn)
	tokenMaker, err := token.NewJWTMaker(config.PrivateKey, config.PublicKey)
	if err != nil {
		log.Fatal("error in initializing token maker", err)
	}

	server := api.NewServer(store, tokenMaker)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can not start server", err)
	}

}
