package db

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type BulkEntity struct {
	EntityName         string                 `json:"entity_name"`
	EntityDescription  string                 `json:"entity_description"`
	DataType           string                 `json:"data_type"`
	IntegrationSource  string                 `json:"integration_source"`
	Template           int                    `json:"template"`
	EntityType         string                 `json:"entity_type"`
	SpecificType       string                 `json:"specific_type"`
	InstanceProducedBy string                 `json:"instance_produced_by"`
	InstanceMetadata   map[string]interface{} `json:"instance_metadata"`
	InstanceCreatedAt  string                 `json:"instance_created_at"`
	LatitudeDegrees    float64                `json:"latitude_degrees"`
	LongitudeDegrees   float64                `json:"longitude_degrees"`
	HeadingDegrees     float64                `json:"heading_degrees"`
	AltitudeHaeMeters  float64                `json:"altitude_hae_meters"`
	SpeedMps           float64                `json:"speed_mps"`
}

func TestBulkCreateEntities(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	now2 := now.Add(time.Minute)

	payload := fmt.Sprintf(`[
		{
			"entity_name": "entity1",
			"entity_description": "desc1",
			"integration_source": "telemetry",
			"template": 1,
			"instance_produced_by": "system",
			"instance_metadata": {"key": "value1"},
			"instance_created_at": "%s",
			"latitude_degrees": 36.7749,
			"longitude_degrees": -123.4194,
			"heading_degrees": 90,
			"altitude_hae_meters": 100,
			"speed_mps": 22
		},
		{
			"entity_name": "entity2",
			"entity_description": "desc2",
			"integration_source": "gps",
			"template": 2,
			"instance_produced_by": "user",
			"instance_metadata": {"key": "value2"},
			"instance_created_at": "%s",
			"latitude_degrees": null,
			"longitude_degrees": null,
			"heading_degrees": null,
			"altitude_hae_meters": null,
			"speed_mps": null
		}
	]`, now.Format(time.RFC3339), now2.Format(time.RFC3339))

	err := testQueries.BulkCreateEntities(context.Background(), []byte(payload))
	require.NoError(t, err, "Bulk create should not return an error")
}

func TestCustomBulkCreateEntities(t *testing.T) {

	payload := make([]BulkEntity, 1000)
	for i := 0; i < 1000; i++ {
		now := time.Now().UTC()

		entityNames := []string{"entity", "entity2", "entity3"}
		producedBy := fmt.Sprintf("%s%d", entityNames[i%3], i)
		integrationSources := []string{"mole", "tak", "blue"}
		payload[i] = BulkEntity{
			EntityName:         producedBy,
			EntityDescription:  "desc",
			IntegrationSource:  integrationSources[i%3],
			EntityType:         "type",
			SpecificType:       "specific_type",
			DataType:           "data_type",
			Template:           1,
			InstanceProducedBy: producedBy,
			InstanceMetadata:   map[string]interface{}{"key": "value"},
			InstanceCreatedAt:  now.Format(time.RFC3339),
			LatitudeDegrees:    36.2313,
			LongitudeDegrees:   -112.2345,
			HeadingDegrees:     55,
			AltitudeHaeMeters:  55,
			SpeedMps:           22,
		}
	}

	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	err = testQueries.BulkCreateEntities(context.Background(), []byte(payloadBytes))
	require.NoError(t, err, "Bulk create should not return an error")
}
