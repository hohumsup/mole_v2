package handlers

import (
	"context"
	db "mole/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateEntityRequest struct {
	Name              string `json:"name" binding:"required"`
	Description       string `json:"description" binding:"required"`
	IntegrationSource string `json:"integration_source" binding:"required"`
}

// CreateEntity returns a Gin handler for creating an entity.
func CreateEntity(query *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateEntityRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		arg := db.CreateEntityParams{
			Name:              req.Name,
			Description:       req.Description,
			IntegrationSource: req.IntegrationSource,
		}

		entity, err := query.CreateEntity(context.Background(), arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, entity)
	}
}
