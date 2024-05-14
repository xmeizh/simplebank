package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // Without this package our server cannot connect to the db
	"github.com/xmeizh/simplebank/api"
	db "github.com/xmeizh/simplebank/db/postgresql"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:mMbzKhVc2DTye79dNfMts@127.0.0.1:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
