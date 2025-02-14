package models

import (
	"database/sql"
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
}

type CreateInstance struct {
	EntityID  uuid.UUID `json:"entity_id" binding:"required"`
	CreatedAt time.Time `json:"created_at"` // Required to track an entity's instance
}

type CreatePosition struct {
	InstanceID        int64    `json:"instance_id"`
	LatitudeDegrees   float64  `json:"latitude_degrees" binding:"required"`
	LongitudeDegrees  float64  `json:"longitude_degrees" binding:"required"`
	HeadingDegrees    *float64 `json:"heading_degrees"`     // Nullable
	AltitudeHaeMeters *float64 `json:"altitude_hae_meters"` // Nullable
	SpeedMps          *float64 `json:"speed_mps"`           // Nullable
}

type CreateEntityResponse struct {
	EntityID          uuid.UUID       `json:"entity_id"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	IntegrationSource string          `json:"integration_source"`
	Template          int32           `json:"template"`
	InstanceID        int64           `json:"instance_id"`
	Position          *CreatePosition `json:"position,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
}
