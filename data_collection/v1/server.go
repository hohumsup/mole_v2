package v1

import (
	"fmt"
	"mole/data_collection/v1/handlers"
	db "mole/db/sqlc"

	"github.com/gin-gonic/gin"
)

// Server serves CRUD operations for data collection.
type Server struct {
	Query  db.Querier  // Exported field for database queries
	Router *gin.Engine // Exported field for the Gin router
}

// DataCollectionServer initializes a new server.
func DataCollectionServer(query db.Querier, router *gin.Engine) *Server {
	server := &Server{
		Query:  query,
		Router: router,
	}

	// Register routes
	server.RegisterRoutes()

	for _, route := range server.Router.Routes() {
		fmt.Println(route.Method, route.Path)
	}

	return server
}

// RegisterRoutes registers all HTTP routes for the server.
func (server *Server) RegisterRoutes() {
	api := server.Router.Group("/v1/api")
	api.POST("/entity", handlers.CreateEntity(server.Query))
	api.GET("/entity/instances", handlers.GetInstances(server.Query))
}

// Start runs the server on the specified address.
func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}

// ErrorResponse formats error messages for HTTP responses.
func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
