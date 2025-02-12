package main

import (
	"database/sql"
	"log"
	v1 "mole/data_collection/v1"
	db "mole/db/sqlc"
	"mole/util"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	query := db.New(conn)

	router := gin.Default()
	server := v1.DataCollectionServer(query, router)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
