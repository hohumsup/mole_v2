package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CreateEntityRequest struct {
	Name              string          `json:"name" binding:"required"`
	Description       string          `json:"description" binding:"required"`
	DataType          sql.NullString  `json:"data_type"`   // Nullable
	SourceName        sql.NullString  `json:"source_name"` // Nullable
	IntegrationSource string          `json:"integration_source"`
	Template          int32           `json:"template" binding:"required"`
	CreatedAt         *time.Time      `json:"created_at,omitempty"` // Optional instance timestamp
	Position          *CreatePosition `json:"position,omitempty"`   // Optional position data
	Instance          *CreateInstance `json:"instance,omitempty"`   // Optional instance data
}

type CreateInstance struct {
	EntityID   uuid.UUID        `json:"entity_id"`
	CreatedAt  time.Time        `json:"created_at"` // Required to track an entity's instance
	ProducedBy *string          `json:"produced_by"`
	Metadata   *json.RawMessage `json:"metadata"`
}

type CreatePosition struct {
	InstanceID        uuid.UUID `json:"instance_id"`
	LatitudeDegrees   float64   `json:"latitude_degrees" binding:"required"`
	LongitudeDegrees  float64   `json:"longitude_degrees" binding:"required"`
	HeadingDegrees    *float64  `json:"heading_degrees"`     // Nullable
	AltitudeHaeMeters *float64  `json:"altitude_hae_meters"` // Nullable
	SpeedMps          *float64  `json:"speed_mps"`           // Nullable
}

type CreateEntityResponse struct {
	EntityID          uuid.UUID       `json:"entity_id"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	IntegrationSource string          `json:"integration_source"`
	Template          int32           `json:"template"`
	InstanceID        uuid.UUID       `json:"instance_id"`
	Position          *CreatePosition `json:"position,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
}

type GetInstances struct {
	EntityID          uuid.UUID        `json:"entity_id"`
	Name              string           `json:"name"`
	IntegrationSource string           `json:"integration_source"`
	Template          int32            `json:"template"`
	CreatedAt         time.Time        `json:"created_at"`
	ModifiedAt        time.Time        `json:"modified_at"`
	InstanceID        uuid.UUID        `json:"instance_id"`
	ProducedBy        *string          `json:"produced_by"`
	LatitudeDegrees   float64          `json:"latitude_degrees"`
	LongitudeDegrees  float64          `json:"longitude_degrees"`
	HeadingDegrees    *float64         `json:"heading_degrees"`
	AltitudeHaeMeters *float64         `json:"altitude_hae_meters"`
	SpeedMps          *float64         `json:"speed_mps"`
	Metadata          *json.RawMessage `json:"metadata"`
}

// GetLatestInstancesResponse defines the JSON response for the latest instance endpoint.
type GetLatestInstancesResponse struct {
	InstanceID        uuid.UUID `json:"instance_id"`
	EntityID          uuid.UUID `json:"entity_id"`
	ProducedBy        string    `json:"produced_by"` // Converts sql.NullString to string (empty if null)
	IntegrationSource string    `json:"integration_source"`
	CreatedAt         time.Time `json:"created_at"`
	ModifiedAt        time.Time `json:"modified_at"`
	Metadata          string    `json:"metadata"` // You can unmarshal the JSON as a string; adjust if you prefer a structured type
	Name              string    `json:"name"`
}

// GetHistoricalInstancesResponse defines the JSON response for the historical (aggregated) endpoint.
type GetHistoricalInstance struct {
	InstanceID        uuid.UUID       `json:"instance_id"`
	EntityID          uuid.UUID       `json:"entity_id"`
	ProducedBy        string          `json:"produced_by"` // Adjust if using sql.NullString
	IntegrationSource string          `json:"integration_source"`
	CreatedAt         time.Time       `json:"created_at"`
	ModifiedAt        time.Time       `json:"modified_at"`
	Metadata          json.RawMessage `json:"metadata"` // Adjust if using a different JSON type
	EntityName        string          `json:"entity_name"`
}
