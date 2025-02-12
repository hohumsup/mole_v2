package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	converters "mole/data_collection/internal/converters"
	"mole/data_collection/v1/models"
	model "mole/data_collection/v1/models"
	db "mole/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateEntity returns a Gin handler for creating an entity
func CreateEntity(query *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.CreateEntityRequest

		// TODO: Create a reusable function for JSON validation
		if err := c.ShouldBindJSON(&req); err != nil {
			var syntaxError *json.SyntaxError
			var unmarshalTypeError *json.UnmarshalTypeError

			switch {
			case errors.As(err, &syntaxError):
				// Invalid JSON or malformed
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid JSON format",
					"details": err.Error(),
				})
				return

			case errors.As(err, &unmarshalTypeError):
				// Field is present with incorrect type
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":   "Invalid field type",
					"details": err.Error(),
				})
				return

			default:
				// Handle unknown field names or missing required fields
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":   "Invalid request fields",
					"details": err.Error(),
				})
				return
			}
		}

		// Step 1: Check if Entity exists
		entity, err := query.GetEntityByNameAndIntegrationSource(context.Background(), db.GetEntityByNameAndIntegrationSourceParams{
			Name:              req.Name,
			IntegrationSource: req.IntegrationSource,
		})

		var entityID uuid.UUID
		// If the Entity does not exist, create it
		if errors.Is(err, sql.ErrNoRows) {
			arg := db.CreateEntityParams{
				Name:              req.Name,
				Description:       req.Description,
				IntegrationSource: req.IntegrationSource,
				Template:          req.Template,
			}

			_, err := query.CreateEntity(context.Background(), arg)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Fetch the newly created entity
			entity, _ = query.GetEntityByNameAndIntegrationSource(context.Background(), db.GetEntityByNameAndIntegrationSourceParams{
				Name:              req.Name,
				IntegrationSource: req.IntegrationSource,
			})

			entityID = entity.EntityID

		} else if err != nil && req.CreatedAt == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			if req.CreatedAt == nil {
				c.JSON(http.StatusConflict, gin.H{"error": "Entity already exists and no instance creation timestamp was provided"})
				return
			}
			entityID = entity.EntityID
		}

		// Step 2: Check for `created_at` to determine if an instance should be created
		var instanceID *int64
		if req.CreatedAt != nil {
			instanceArg := db.InsertInstanceParams{
				EntityID:  entityID,
				CreatedAt: *req.CreatedAt,
			}

			newInstanceID, err := query.InsertInstance(context.Background(), instanceArg)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			instanceID = &newInstanceID
		}

		// Step 3: Check for `position` and/or `geo_details`
		if instanceID != nil && req.Position != nil {
			positionArg := db.InsertPositionParams{
				InstanceID:        *instanceID,
				LatitudeDegrees:   req.Position.LatitudeDegrees,
				LongitudeDegrees:  req.Position.LongitudeDegrees,
				HeadingDegrees:    converters.Float64ToNullFloat64(req.Position.HeadingDegrees),
				AltitudeHaeMeters: converters.Float64ToNullFloat64(req.Position.AltitudeHaeMeters),
				SpeedMps:          converters.Float64ToNullFloat64(req.Position.SpeedMps),
			}

			err := query.InsertPosition(context.Background(), positionArg)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// Step 4: Return Response With Entity
		response := models.CreateEntityResponse{
			EntityID:          entityID,
			Name:              entity.Name,
			Description:       entity.Description,
			IntegrationSource: entity.IntegrationSource,
			Template:          entity.Template,
		}

		if instanceID != nil {
			response.InstanceID = *instanceID
			response.CreatedAt = *req.CreatedAt
		}

		c.JSON(http.StatusOK, response)
	}
}
