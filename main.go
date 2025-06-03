package main

import (
	"database/sql"
	"github.com/EliriaT/dnd-user-service/config"
	"github.com/EliriaT/dnd-user-service/db"
	"github.com/EliriaT/dnd-user-service/server"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	configSet, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("can not open config file")
	}
	conn, err := sql.Open(configSet.DBdriver, configSet.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	queries := db.New(conn)

	srv, _ := server.NewServer(queries, configSet)

	err = srv.Start(configSet.ServerAddress)

	if err != nil {
		log.Fatal("server can not be started. ", err)
	}
}
