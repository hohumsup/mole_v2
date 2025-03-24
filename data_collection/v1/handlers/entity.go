package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	converters "mole/data_collection/internal/converters"
	"mole/data_collection/v1/models"
	db "mole/db/sqlc"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// CreateEntity returns a Gin handler for creating an entity.
func CreateEntity(query db.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateEntityRequest

		// Validate JSON.
		if err := c.ShouldBindJSON(&req); err != nil {
			var syntaxError *json.SyntaxError
			var unmarshalTypeError *json.UnmarshalTypeError
			switch {
			case errors.As(err, &syntaxError):
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid JSON format",
					"details": fmt.Sprintf("Syntax error at byte offset %d: %v", syntaxError.Offset, syntaxError.Error()),
				})
				return
			case errors.As(err, &unmarshalTypeError):
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":   "Invalid field type",
					"details": fmt.Sprintf("Field '%s': expected type %s but got value %v", unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Value),
				})
				return
			default:
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":   "Invalid request fields",
					"details": err.Error(),
				})
				return
			}
		}

		// Step 1: Check if Entity exists.
		entity, err := query.GetEntityByNameAndIntegrationSource(context.Background(), db.GetEntityByNameAndIntegrationSourceParams{
			Name:              req.Name,
			IntegrationSource: req.IntegrationSource,
		})

		var entityID uuid.UUID
		if errors.Is(err, sql.ErrNoRows) {
			arg := db.CreateEntityParams{
				Name:              req.Name,
				Description:       req.Description,
				IntegrationSource: req.IntegrationSource,
				Template:          req.Template,
			}
			_, err := query.CreateEntity(context.Background(), arg)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			// Fetch the newly created entity.
			entity, _ = query.GetEntityByNameAndIntegrationSource(context.Background(), db.GetEntityByNameAndIntegrationSourceParams{
				Name:              req.Name,
				IntegrationSource: req.IntegrationSource,
			})
			entityID = entity.EntityID
		} else {
			if req.CreatedAt == nil {
				c.JSON(http.StatusConflict, gin.H{"error": "Entity already exists and no instance creation timestamp was provided"})
				return
			}
			entityID = entity.EntityID
		}

		// Step 2: Create a new instance (using client-provided created_at).
		var instanceID *uuid.UUID
		var instanceCreatedAt time.Time
		if req.CreatedAt != nil {
			normalizedTime := req.CreatedAt.UTC().Truncate(time.Second)
			log.Printf("Instance created_at before DB %v\n", normalizedTime.Format(time.RFC3339Nano))

			instanceArg := db.InsertInstanceParams{
				EntityID:  entityID,
				CreatedAt: normalizedTime, // Use the client-provided timestamp.
			}
			if req.Instance != nil {
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
			newInstance, err := query.InsertInstance(context.Background(), instanceArg)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			instanceID = &newInstance.InstanceID
			log.Printf("NormalizedTime UnixNano: %d", normalizedTime.UnixNano())
			log.Printf("DB Returned CreatedAt UnixNano: %d", newInstance.CreatedAt.UnixNano())
			instanceCreatedAt = newInstance.CreatedAt // Use the returned value!
		}

		// Step 3: Insert the new position for the new instance.
		if instanceID != nil && req.Position != nil {
			log.Printf("Position %v\n", instanceCreatedAt.Format(time.RFC3339Nano))
			positionArg := db.InsertPositionParams{
				InstanceID:        *instanceID,
				InstanceCreatedAt: instanceCreatedAt, // This should match the instance's created_at.
				LatitudeDegrees:   req.Position.LatitudeDegrees,
				LongitudeDegrees:  req.Position.LongitudeDegrees,
				HeadingDegrees:    converters.Float64ToNullFloat64(req.Position.HeadingDegrees),
				AltitudeHaeMeters: converters.Float64ToNullFloat64(req.Position.AltitudeHaeMeters),
				SpeedMps:          converters.Float64ToNullFloat64(req.Position.SpeedMps),
			}

			log.Printf("Position %v\n", positionArg)
			err := query.InsertPosition(context.Background(), positionArg)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				log.Printf("Error %v\n", err)
				return
			}
		}

		// Step 4: Build and return the response.
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

func GetInstances(query db.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		positions, err := query.GetInstances(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response []models.GetInstances
		for _, pos := range positions {
			response = append(response, models.GetInstances{
				EntityID:          pos.EntityID,
				Name:              pos.EntityName,
				IntegrationSource: pos.IntegrationSource,
				Template:          pos.Template,
				CreatedAt:         pos.InstanceCreatedAt,
				ModifiedAt:        pos.ModifiedAt,
				InstanceID:        pos.InstanceID,
				ProducedBy:        &pos.ProducedBy.String,
				LatitudeDegrees:   pos.LatitudeDegrees,
				LongitudeDegrees:  pos.LongitudeDegrees,
				HeadingDegrees:    &pos.HeadingDegrees.Float64,
				AltitudeHaeMeters: &pos.AltitudeHaeMeters.Float64,
				SpeedMps:          &pos.SpeedMps.Float64,
				Metadata:          &pos.Metadata.RawMessage,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}

// GetLatestInstances returns the latest instance for each entity.
func GetLatestInstances(query db.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		latestInstances, err := query.GetLatestInstances(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Map the database rows to your response model.
		var response []models.GetLatestInstancesResponse
		for _, inst := range latestInstances {
			response = append(response, models.GetLatestInstancesResponse{
				InstanceID:        inst.InstanceID,
				EntityID:          inst.EntityID,
				IntegrationSource: inst.IntegrationSource,
				ProducedBy:        inst.ProducedBy.String, // or use a helper to handle nulls
				CreatedAt:         inst.CreatedAt,
				ModifiedAt:        inst.ModifiedAt,
				Metadata:          string(inst.Metadata.RawMessage),
				Name:              inst.Name,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetHistoricalInstances(query db.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Expect an "interval" query parameter (e.g., "1 hour", "30 minutes", etc.)
		interval := c.Query("interval")
		if interval == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "interval parameter is required"})
			return
		}
		// Convert interval to int64.
		intervalInt, err := strconv.ParseInt(interval, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid interval parameter"})
			return
		}

		// Call the sqlc-generated query.
		rows, err := query.GetHistoricalInstances(c.Request.Context(), intervalInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Map the returned rows to our model.
		var result []models.GetHistoricalInstance
		for _, row := range rows {
			// If ProducedBy is sql.NullString, use row.ProducedBy.String.
			// Similarly for Metadata if needed.
			instance := models.GetHistoricalInstance{
				InstanceID:        row.InstanceID,
				EntityID:          row.EntityID,
				IntegrationSource: row.IntegrationSource,
				ProducedBy:        row.ProducedBy.String, // Adjust if necessary
				CreatedAt:         row.CreatedAt,
				ModifiedAt:        row.ModifiedAt,
				Metadata:          row.Metadata.RawMessage, // Adjust if necessary
				EntityName:        row.EntityName,
			}
			result = append(result, instance)
		}

		c.JSON(http.StatusOK, result)
	}
}
