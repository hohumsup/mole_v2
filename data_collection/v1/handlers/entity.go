package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	converters "mole/data_collection/internal/converters"
	"mole/data_collection/v1/models"
	db "mole/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// CreateEntity returns a Gin handler for creating an entity
func CreateEntity(query db.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateEntityRequest

		// TODO: Create a reusable function for JSON validation
		if err := c.ShouldBindJSON(&req); err != nil {
			var syntaxError *json.SyntaxError
			var unmarshalTypeError *json.UnmarshalTypeError

			switch {
			case errors.As(err, &syntaxError):
				// Invalid JSON
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid JSON format",
					"details": fmt.Sprintf("Syntax error at byte offset %d: %v", syntaxError.Offset, syntaxError.Error()),
				})
				return

			case errors.As(err, &unmarshalTypeError):
				// Field type mismatch
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":   "Invalid field type",
					"details": fmt.Sprintf("Field '%s': expected type %s but got value %v", unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Value),
				})
				return

			default:
				// Other errors
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
				// Error due to invalid template
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Fetch the newly created entity
			entity, _ = query.GetEntityByNameAndIntegrationSource(context.Background(), db.GetEntityByNameAndIntegrationSourceParams{
				Name:              req.Name,
				IntegrationSource: req.IntegrationSource,
			})

			entityID = entity.EntityID

			// } else if err != nil && req.CreatedAt == nil {
			// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			// 	return
		} else {
			if req.CreatedAt == nil {
				c.JSON(http.StatusConflict, gin.H{"error": "Entity already exists and no instance creation timestamp was provided"})
				return
			}
			entityID = entity.EntityID
		}

		// Step 2: Check for `created_at` to determine if an instance should be created
		var instanceID *uuid.UUID
		if req.CreatedAt != nil {
			instanceArg := db.InsertInstanceParams{
				EntityID:  entityID,
				CreatedAt: *req.CreatedAt,
			}

			if req.Instance != nil {
				// Convert *string to sql.NullString for ProducedBy.
				if req.Instance.ProducedBy != nil {
					instanceArg.ProducedBy = sql.NullString{
						String: *req.Instance.ProducedBy,
						Valid:  true,
					}
				} else {
					instanceArg.ProducedBy = sql.NullString{Valid: false}
				}

				if req.Instance.Metadata != nil {
					meta, err := converters.ConvertJSONToPQType(*req.Instance.Metadata)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error":   "Invalid JSON in metadata",
							"details": err.Error(),
						})
						return
					}
					instanceArg.Metadata = meta
				} else {
					instanceArg.Metadata = pqtype.NullRawMessage{Valid: false}
				}

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

		if req.Position != nil {
			response.Position = &models.CreatePosition{
				InstanceID:        *instanceID,
				LatitudeDegrees:   req.Position.LatitudeDegrees,
				LongitudeDegrees:  req.Position.LongitudeDegrees,
				HeadingDegrees:    req.Position.HeadingDegrees,
				AltitudeHaeMeters: req.Position.AltitudeHaeMeters,
				SpeedMps:          req.Position.SpeedMps,
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetPositions(query db.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		positions, err := query.GetPositions(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response []models.GetPositions
		for _, pos := range positions {
			response = append(response, models.GetPositions{
				EntityID:          pos.EntityID,
				Name:              pos.EntityName,
				IntegrationSource: pos.IntegrationSource,
				Template:          pos.Template,
				CreatedAt:         pos.CreatedAt,
				ModifiedAt:        pos.ModifiedAt,
				InstanceID:        pos.InstanceID,
				LatitudeDegrees:   pos.LatitudeDegrees,
				LongitudeDegrees:  pos.LongitudeDegrees,
				HeadingDegrees:    converters.NullFloat64ToFloat64(pos.HeadingDegrees),
				AltitudeHaeMeters: converters.NullFloat64ToFloat64(pos.AltitudeHaeMeters),
				SpeedMps:          converters.NullFloat64ToFloat64(pos.SpeedMps),
			})
		}

		c.JSON(http.StatusOK, response)
	}
}
