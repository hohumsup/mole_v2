package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	v1 "mole/data_collection/v1"
	"mole/data_collection/v1/models"
	mock_db "mole/db/mock"
	db "mole/db/sqlc"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateEntityAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock_db.NewMockQuerier(ctrl)
	expectedCreateEntityParams := createEntity()
	expectedGetEntityParams := getEntityByNameAndIntegrationSource()
	createdEntityID := uuid.New() // mock entity ID since there's isnt an actual DB to generate one

	// instance
	expectedCreatedAt := time.Now().UTC()

	// position
	expectedLatitude := 37.7749
	expectedLongitude := -122.4194
	expectedHeading := 90.0 // example value
	expectedAltitude := 100.0
	expectedSpeed := 10.0

	// TODO: Need to define test cases for the following:
	// 1. Invalid JSON format
	// 2. Invalid field type
	// 3. Invalid request fields
	// 4. Create entity without instance and position
	// 5. Create entity with instance but without position
	// 6. Create entity with instance and position

	gomock.InOrder(
		// 1. First check: Entity not found.
		mockQuerier.EXPECT().
			GetEntityByNameAndIntegrationSource(gomock.Any(), expectedGetEntityParams).
			Return(db.GetEntityByNameAndIntegrationSourceRow{}, sql.ErrNoRows),

		// 2. Create the entity.
		mockQuerier.EXPECT().
			CreateEntity(gomock.Any(), expectedCreateEntityParams).
			Return(db.CreateEntityRow{
				EntityID:          createdEntityID,
				Name:              expectedCreateEntityParams.Name,
				Description:       expectedCreateEntityParams.Description,
				IntegrationSource: expectedCreateEntityParams.IntegrationSource,
				Template:          expectedCreateEntityParams.Template,
			}, nil),

		// 3. Second check: Entity found.
		mockQuerier.EXPECT().
			GetEntityByNameAndIntegrationSource(gomock.Any(), expectedGetEntityParams).
			Return(db.GetEntityByNameAndIntegrationSourceRow{
				EntityID:          createdEntityID,
				Name:              expectedCreateEntityParams.Name,
				Description:       expectedCreateEntityParams.Description,
				IntegrationSource: expectedCreateEntityParams.IntegrationSource,
				Template:          expectedCreateEntityParams.Template,
			}, nil),

		// 4. Insert instance (requires a timestamp)
		mockQuerier.EXPECT().
			InsertInstance(gomock.Any(), db.InsertInstanceParams{
				EntityID:  createdEntityID,
				CreatedAt: expectedCreatedAt,
			}).
			Return(int64(1), nil), // mock instance ID

		// 5. Insert position
		mockQuerier.EXPECT().
			InsertPosition(gomock.Any(), db.InsertPositionParams{
				InstanceID:        1,
				LatitudeDegrees:   expectedLatitude,
				LongitudeDegrees:  expectedLongitude,
				HeadingDegrees:    sql.NullFloat64{Float64: expectedHeading, Valid: true},
				AltitudeHaeMeters: sql.NullFloat64{Float64: expectedAltitude, Valid: true},
				SpeedMps:          sql.NullFloat64{Float64: expectedSpeed, Valid: true},
			}).
			Return(nil),
	)

	// Initialize Server with mockQuerier
	server := v1.DataCollectionServer(mockQuerier, gin.Default())

	// Build the JSON payload for the POST request.
	reqPayload := models.CreateEntityRequest{
		Name:              expectedCreateEntityParams.Name,
		Description:       expectedCreateEntityParams.Description,
		IntegrationSource: expectedCreateEntityParams.IntegrationSource,
		Template:          expectedCreateEntityParams.Template,
		CreatedAt:         &expectedCreatedAt,
		Position: &models.CreatePosition{
			LatitudeDegrees:   expectedLatitude,
			LongitudeDegrees:  expectedLongitude,
			HeadingDegrees:    &expectedHeading,
			AltitudeHaeMeters: &expectedAltitude,
			SpeedMps:          &expectedSpeed,
		},
	}
	payloadBytes, err := json.Marshal(reqPayload)
	require.NoError(t, err)

	// Build the HTTP POST request.
	url := "/v1/api/entity"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	// Create a recorder to capture the response.
	recorder := httptest.NewRecorder()
	server.Router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)

	var resp models.CreateEntityResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, createdEntityID, resp.EntityID)
	require.Equal(t, int64(1), resp.InstanceID)
	t.Logf("Response Body: %s", recorder.Body.String())
}

func createEntity() db.CreateEntityParams {
	return db.CreateEntityParams{
		Name:              "charlie",
		Description:       "mole generated entity",
		IntegrationSource: "mole",
		Template:          2,
	}
}

func getEntityByNameAndIntegrationSource() db.GetEntityByNameAndIntegrationSourceParams {
	return db.GetEntityByNameAndIntegrationSourceParams{
		Name:              "charlie",
		IntegrationSource: "mole",
	}
}
