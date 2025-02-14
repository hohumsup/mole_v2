package tests

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

		// mockQuerier.EXPECT().
		// 	InsertInstance(gomock.Any(), db.InsertInstanceParams{
		// 		EntityID:  createdEntityID,
		// 		CreatedAt: time.Now().UTC(),
		// 	}).
		// 	Return(db.InsertInstance, nil),
	)

	// Initialize Server with mockQuerier
	server := v1.DataCollectionServer(mockQuerier, gin.Default())

	// Build the JSON payload for the POST request.
	reqPayload := models.CreateEntityRequest{
		Name:              expectedCreateEntityParams.Name,
		Description:       expectedCreateEntityParams.Description,
		IntegrationSource: expectedCreateEntityParams.IntegrationSource,
		Template:          expectedCreateEntityParams.Template,
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
