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
	createdEntityIDAsset := uuid.New() // mock entity ID since there's isnt an actual DB to generate one
	createdEntityIDEvent := uuid.New() // mock another ID for a different entity

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

	testCases := []struct {
		name             string
		rawPayload       string
		payload          models.CreateEntityRequest
		setupMocks       func(ctrl *gomock.Controller, mockQuerier *mock_db.MockQuerier)
		expectedHTTPCode int
		expectedEntityID uuid.UUID
		expectedInstID   int64
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

			},
			expectedHTTPCode: http.StatusOK,
			expectedEntityID: createdEntityIDAsset,
			expectedInstID:   1,
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
						Return(int64(2), nil),
				)
			},
			expectedHTTPCode: http.StatusOK,
			expectedEntityID: createdEntityIDEvent,
			expectedInstID:   2,
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
				expectedGetEntityParams := getEntityByNameAndIntegrationSource("detection")

				mockQuerier.EXPECT().
					GetEntityByNameAndIntegrationSource(gomock.Any(), expectedGetEntityParams).
					Return(db.GetEntityByNameAndIntegrationSourceRow{
						EntityID: createdEntityIDEvent,
						Name:     "detection",
					}, nil)
			},
			expectedHTTPCode: http.StatusConflict,
			expectedEntityID: uuid.Nil,
			expectedInstID:   0,
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
					// 1. First check: entity not found.
					mockQuerier.EXPECT().
						GetEntityByNameAndIntegrationSource(gomock.Any(), expectedGetEntityParams).
						Return(db.GetEntityByNameAndIntegrationSourceRow{}, sql.ErrNoRows),

					// 2. Attempt to create the entity with invalid template.
					mockQuerier.EXPECT().
						CreateEntity(gomock.Any(), expectedCreateEntityParams).
						Return(db.CreateEntityRow{}, sql.ErrNoRows),
				)
			},
			expectedHTTPCode: http.StatusBadRequest,
			expectedEntityID: uuid.Nil,
			expectedInstID:   0,
		},
		{
			name:       "malformed JSON payload",
			rawPayload: `{iption": "mole generated entity", "integrati`,
			setupMocks: func(ctrl *gomock.Controller, mockQuerier *mock_db.MockQuerier) {
				// No mock setup needed for validating request fields

			},
			expectedHTTPCode: http.StatusBadRequest,
			expectedEntityID: uuid.Nil,
			expectedInstID:   0,
			expectedError:    `{"details":"invalid character 'i' looking for beginning of object key string", "error":"Invalid JSON format"}`,
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
				// No mock setup needed for validating request fields
			},
			expectedHTTPCode: http.StatusUnprocessableEntity,
			expectedEntityID: uuid.Nil,
			expectedInstID:   0,
		},
	}

	// Run each test case as a subtest.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockQuerier := mock_db.NewMockQuerier(ctrl)

			// Set up mocks only if a raw payload is not provided because when raw payload is provided,
			// we expect binding to fail since it does not match the expected JSON structure.
			if tc.setupMocks != nil {
				tc.setupMocks(ctrl, mockQuerier)
			}

			server := v1.DataCollectionServer(mockQuerier, gin.Default())

			var req *http.Request
			var err error

			// Raw payload is provided for validating JSON errors
			if tc.rawPayload != "" {
				req, err = http.NewRequest(http.MethodPost, "/v1/api/entity", bytes.NewBuffer([]byte(tc.rawPayload)))
			} else {
				payloadBytes, err := json.Marshal(tc.payload)
				require.NoError(t, err)
				req, err = http.NewRequest(http.MethodPost, "/v1/api/entity", bytes.NewBuffer(payloadBytes))
				require.NoError(t, err)
			}

			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			server.Router.ServeHTTP(recorder, req)

			require.Equal(t, tc.expectedHTTPCode, recorder.Code)

			// For success responses, unmarshal and check fields.
			if tc.expectedHTTPCode == http.StatusOK {
				var resp models.CreateEntityResponse
				err = json.Unmarshal(recorder.Body.Bytes(), &resp)
				require.NoError(t, err)
				t.Logf("Response Body: %s", recorder.Body.String())
			} else if tc.rawPayload != "" {
				// For validating JSON errors, compare the expected error message.
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
