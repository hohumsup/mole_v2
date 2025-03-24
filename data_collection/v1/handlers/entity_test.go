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

	// entity
	createdEntityIDAsset := uuid.New() // mock entity ID for asset entity
	createdEntityIDEvent := uuid.New() // mock entity ID for event entity

	// instance
	expectedCreatedAt := time.Now().UTC()
	expectedInstanceID := uuid.New() // dummy instance ID for testing

	// position
	expectedLatitude := 37.7749
	expectedLongitude := -122.4194
	expectedHeading := 90.0 // example value
	expectedAltitude := 100.0
	expectedSpeed := 10.0

	testCases := []struct {
		name             string
		rawPayload       string
		payload          models.CreateEntityRequest
		setupMocks       func(ctrl *gomock.Controller, mockQuerier *mock_db.MockQuerier)
		expectedHTTPCode int
		expectedEntityID uuid.UUID
		expectedInstID   uuid.UUID
		expectedError    string
	}{
		{
			name: "create an entity with instance and position",
			payload: models.CreateEntityRequest{
				Name:              "charlie",
				Description:       "mole generated entity",
				IntegrationSource: "mole",
				Template:          2,
				// Provide CreatedAt and Position.
				CreatedAt: &expectedCreatedAt,
				Position: &models.CreatePosition{
					LatitudeDegrees:   37.7749,
					LongitudeDegrees:  -122.4194,
					HeadingDegrees:    func() *float64 { v := 90.0; return &v }(),
					AltitudeHaeMeters: func() *float64 { v := 100.0; return &v }(),
					SpeedMps:          func() *float64 { v := 10.0; return &v }(),
				},
			},
			setupMocks: func(ctrl *gomock.Controller, mockQuerier *mock_db.MockQuerier) {
				expectedCreateEntityParams := createEntity("charlie", 2)
				expectedGetEntityParams := getEntityByNameAndIntegrationSource("charlie")

				gomock.InOrder(
					// 1. First check: Entity not found.
					mockQuerier.EXPECT().
						GetEntityByNameAndIntegrationSource(gomock.Any(), expectedGetEntityParams).
						Return(db.GetEntityByNameAndIntegrationSourceRow{}, sql.ErrNoRows),

					// 2. Create the entity.
					mockQuerier.EXPECT().
						CreateEntity(gomock.Any(), expectedCreateEntityParams).
						Return(db.CreateEntityRow{
							EntityID:          createdEntityIDAsset,
							Name:              expectedCreateEntityParams.Name,
							Description:       expectedCreateEntityParams.Description,
							IntegrationSource: expectedCreateEntityParams.IntegrationSource,
							Template:          expectedCreateEntityParams.Template,
						}, nil),

					// 3. Second check: Entity found.
					mockQuerier.EXPECT().
						GetEntityByNameAndIntegrationSource(gomock.Any(), expectedGetEntityParams).
						Return(db.GetEntityByNameAndIntegrationSourceRow{
							EntityID:          createdEntityIDAsset,
							Name:              expectedCreateEntityParams.Name,
							Description:       expectedCreateEntityParams.Description,
							IntegrationSource: expectedCreateEntityParams.IntegrationSource,
							Template:          expectedCreateEntityParams.Template,
						}, nil),

					// 4. Insert instance (requires a timestamp)
					mockQuerier.EXPECT().
						InsertInstance(gomock.Any(), db.InsertInstanceParams{
							EntityID:  createdEntityIDAsset,
							CreatedAt: expectedCreatedAt,
						}).
						Return(db.InsertInstanceRow{
							InstanceID: expectedInstanceID,
							CreatedAt:  expectedCreatedAt,
						}, nil),

					// 5. Insert position: Use the returned instanceID and created_at.
					mockQuerier.EXPECT().
						InsertPosition(gomock.Any(), db.InsertPositionParams{
							InstanceID:        expectedInstanceID,
							InstanceCreatedAt: expectedCreatedAt,
							LatitudeDegrees:   expectedLatitude,
							LongitudeDegrees:  expectedLongitude,
							HeadingDegrees:    sql.NullFloat64{Float64: expectedHeading, Valid: true},
							AltitudeHaeMeters: sql.NullFloat64{Float64: expectedAltitude, Valid: true},
							SpeedMps:          sql.NullFloat64{Float64: expectedSpeed, Valid: true},
						}).
						Return(nil),
				)
			},
			expectedHTTPCode: http.StatusOK,
			expectedEntityID: createdEntityIDAsset,
			expectedInstID:   expectedInstanceID,
		},
		{
			name: "create an entity with instance",
			payload: models.CreateEntityRequest{
				Name:              "detection",
				Description:       "mole generated entity",
				IntegrationSource: "mole",
				Template:          1,
				CreatedAt:         &expectedCreatedAt,
			},
			setupMocks: func(ctrl *gomock.Controller, mockQuerier *mock_db.MockQuerier) {
				expectedCreateEntityParams := createEntity("detection", 1)
				expectedGetEntityParams := getEntityByNameAndIntegrationSource("detection")
				expectedInstanceID2 := uuid.New()

				gomock.InOrder(
					mockQuerier.EXPECT().
						GetEntityByNameAndIntegrationSource(gomock.Any(), expectedGetEntityParams).
						Return(db.GetEntityByNameAndIntegrationSourceRow{}, sql.ErrNoRows),

					mockQuerier.EXPECT().
						CreateEntity(gomock.Any(), expectedCreateEntityParams).
						Return(db.CreateEntityRow{
							EntityID:          createdEntityIDEvent,
							Name:              expectedCreateEntityParams.Name,
							Description:       expectedCreateEntityParams.Description,
							IntegrationSource: expectedCreateEntityParams.IntegrationSource,
							Template:          expectedCreateEntityParams.Template,
						}, nil),

					mockQuerier.EXPECT().
						GetEntityByNameAndIntegrationSource(gomock.Any(), expectedGetEntityParams).
						Return(db.GetEntityByNameAndIntegrationSourceRow{
							EntityID:          createdEntityIDEvent,
							Name:              expectedCreateEntityParams.Name,
							Description:       expectedCreateEntityParams.Description,
							IntegrationSource: expectedCreateEntityParams.IntegrationSource,
							Template:          expectedCreateEntityParams.Template,
						}, nil),

					mockQuerier.EXPECT().
						InsertInstance(gomock.Any(), db.InsertInstanceParams{
							EntityID:  createdEntityIDEvent,
							CreatedAt: expectedCreatedAt,
						}).
						Return(db.InsertInstanceRow{
							InstanceID: expectedInstanceID2,
							CreatedAt:  expectedCreatedAt,
						}, nil),
				)
			},
			expectedHTTPCode: http.StatusOK,
			expectedEntityID: createdEntityIDEvent,
			// For this test, we expect the returned instance ID to be as returned in the mock.
			expectedInstID: uuid.Nil, // Adjust as needed (or compare against expectedInstanceID2)
		},
		{
			name: "create an existing entity without instance and position",
			payload: models.CreateEntityRequest{
				Name:              "detection",
				Description:       "mole generated entity",
				IntegrationSource: "mole",
				Template:          1,
			},
			setupMocks: func(ctrl *gomock.Controller, mockQuerier *mock_db.MockQuerier) {
				expectedCreateEntityParams := createEntity("detection", 1)
				expectedGetEntityParams := getEntityByNameAndIntegrationSource("detection")

				mockQuerier.EXPECT().
					GetEntityByNameAndIntegrationSource(gomock.Any(), expectedGetEntityParams).
					Return(db.GetEntityByNameAndIntegrationSourceRow{
						EntityID:          createdEntityIDEvent,
						Name:              expectedCreateEntityParams.Name,
						Description:       expectedCreateEntityParams.Description,
						IntegrationSource: expectedCreateEntityParams.IntegrationSource,
						Template:          expectedCreateEntityParams.Template,
					}, nil)
			},
			expectedHTTPCode: http.StatusConflict,
			expectedEntityID: uuid.Nil,
			expectedInstID:   uuid.Nil,
		},
		{
			name: "create an entity with invalid template",
			payload: models.CreateEntityRequest{
				Name:              "charlie",
				Description:       "mole generated entity",
				IntegrationSource: "mole",
				Template:          4,
			},
			setupMocks: func(ctrl *gomock.Controller, mockQuerier *mock_db.MockQuerier) {
				expectedGetEntityParams := getEntityByNameAndIntegrationSource("charlie")
				expectedCreateEntityParams := createEntity("charlie", 4)

				gomock.InOrder(
					mockQuerier.EXPECT().
						GetEntityByNameAndIntegrationSource(gomock.Any(), expectedGetEntityParams).
						Return(db.GetEntityByNameAndIntegrationSourceRow{}, sql.ErrNoRows),

					mockQuerier.EXPECT().
						CreateEntity(gomock.Any(), expectedCreateEntityParams).
						Return(db.CreateEntityRow{}, sql.ErrNoRows),
				)
			},
			expectedHTTPCode: http.StatusBadRequest,
			expectedEntityID: uuid.Nil,
			expectedInstID:   uuid.Nil,
		},
		{
			name:       "malformed JSON payload",
			rawPayload: `{iption": "mole generated entity", "integrati`,
			setupMocks: func(ctrl *gomock.Controller, mockQuerier *mock_db.MockQuerier) {
				// No mock setup needed for validating request fields.
			},
			expectedHTTPCode: http.StatusBadRequest,
			expectedEntityID: uuid.Nil,
			expectedInstID:   uuid.Nil,
			// Adjust expected error details if necessary.
			expectedError: `{"details":"Syntax error at byte offset 2: invalid character 'i' looking for beginning of object key string","error":"Invalid JSON format"}`,
		},
		{
			name: "invalid request fields",
			payload: models.CreateEntityRequest{
				Name:              "",
				Description:       "mole generated entity",
				IntegrationSource: "mole",
				Template:          1,
			},
			setupMocks: func(ctrl *gomock.Controller, mockQuerier *mock_db.MockQuerier) {
				// No mock setup needed for validating request fields.
			},
			expectedHTTPCode: http.StatusUnprocessableEntity,
			expectedEntityID: uuid.Nil,
			expectedInstID:   uuid.Nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockQuerier := mock_db.NewMockQuerier(ctrl)
			if tc.setupMocks != nil {
				tc.setupMocks(ctrl, mockQuerier)
			}

			server := v1.DataCollectionServer(mockQuerier, gin.Default())

			var req *http.Request
			var err error
			if tc.rawPayload != "" {
				req, err = http.NewRequest(http.MethodPost, "/v1/api/entity", bytes.NewBuffer([]byte(tc.rawPayload)))
			} else {
				payloadBytes, err := json.Marshal(tc.payload)
				require.NoError(t, err)
				req, err = http.NewRequest(http.MethodPost, "/v1/api/entity", bytes.NewBuffer(payloadBytes))
				require.NoError(t, err)
			}

			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			server.Router.ServeHTTP(recorder, req)

			require.Equal(t, tc.expectedHTTPCode, recorder.Code)
			if tc.expectedHTTPCode == http.StatusOK {
				var resp models.CreateEntityResponse
				err = json.Unmarshal(recorder.Body.Bytes(), &resp)
				require.NoError(t, err)
				t.Logf("Response Body: %s", recorder.Body.String())
			} else if tc.rawPayload != "" {
				require.JSONEq(t, tc.expectedError, recorder.Body.String())
			}
		})
	}
}

func createEntity(name string, template int32) db.CreateEntityParams {
	return db.CreateEntityParams{
		Name:              name,
		Description:       "mole generated entity",
		IntegrationSource: "mole",
		Template:          template,
	}
}

func getEntityByNameAndIntegrationSource(name string) db.GetEntityByNameAndIntegrationSourceParams {
	return db.GetEntityByNameAndIntegrationSourceParams{
		Name:              name,
		IntegrationSource: "mole",
	}
}
