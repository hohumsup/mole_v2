package main

import (
	"context"
	"log"
	v1 "mole/data_collection/v1"
	db "mole/db/sqlc"
	"mole/util"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	query := db.New(connPool)

	router := gin.Default()
	server := v1.DataCollectionServer(query, router)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
