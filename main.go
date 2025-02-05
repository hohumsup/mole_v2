package main

import (
	"database/sql"
	"log"
	v1 "mole/data_collection/v1"
	db "mole/db/sqlc"

	"github.com/gin-gonic/gin"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://mole_user:secret@localhost:5432/mole?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	query := db.New(conn)

	router := gin.Default()
	server := v1.DataCollectionServer(query, router)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
