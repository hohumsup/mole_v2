package db

import (
	"context"
	"database/sql"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

// Should put this in a helper or utility file
func FormatName(name string) string {
	lower := strings.ToLower(name)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	formatted := re.ReplaceAllString(lower, "_")
	return strings.Trim(formatted, "_")
}

func CreateEntity(t *testing.T, random bool) CreateEntityRow {
	var name string
	if random {
		name = gofakeit.Name()
	} else {
		name = "charlie"
	}

	arg := CreateEntityParams{
		Name:              FormatName(name),
		Description:       gofakeit.Sentence(5),
		IntegrationSource: gofakeit.RandomString([]string{"telemetry", "gps"}),
		Template:          int32(rand.Intn(3) + 1),
	}

	entityRow, err := testQueries.CreateEntity(context.Background(), arg)
	t.Log(arg)
	if err != nil {
		// If the error is due to a duplicate entity, fetch the existing one
		t.Log(err)
		if strings.Contains(err.Error(), "already exists") {
			t.Log("Duplicate entity detected, fetching existing entity:")
			t.Log(err)

			// Fetch the existing entity
			existingEntity, getErr := testQueries.GetEntityByNameAndIntegrationSource(
				context.Background(),
				GetEntityByNameAndIntegrationSourceParams{
					Name:              arg.Name,
					IntegrationSource: arg.IntegrationSource,
				},
			)
			require.NoError(t, getErr)
			require.NotEmpty(t, existingEntity)
			return CreateEntityRow(existingEntity)

		} else {
			t.Log(entityRow)
			// For any other error, fail the test
			require.NoError(t, err)
		}
	}

	t.Log(entityRow)
	require.NoError(t, err)
	require.NotEmpty(t, entityRow)

	// Validate fields match input
	require.Equal(t, arg.Name, entityRow.Name)
	require.Equal(t, arg.Description, entityRow.Description)
	require.Equal(t, arg.IntegrationSource, entityRow.IntegrationSource)

	return entityRow
}

func TestEntityGenerator(t *testing.T) {
	for i := 0; i < 3; i++ {
		_ = CreateEntity(t, false)
	}
}

func TestGetEntity(t *testing.T) {
	entity1 := CreateEntity(t, true)
	entity2, err := testQueries.GetEntity(context.Background(), entity1.EntityID)
	require.NoError(t, err)
	require.NotEmpty(t, entity2)

	require.Equal(t, entity1.EntityID, entity2.EntityID)
	require.Equal(t, entity1.Name, entity2.Name)
	require.Equal(t, entity1.Description, entity2.Description)

	t.Log("Found:", entity1.Name)
}

func TestGetEntityByNameAndIntegrationSource(t *testing.T) {
	entity1 := CreateEntity(t, false) // Creates "charlie" with a random integration source

	// Query for the entity using name + integration_source
	entity2, err := testQueries.GetEntityByNameAndIntegrationSource(
		context.Background(), GetEntityByNameAndIntegrationSourceParams{
			Name:              entity1.Name,
			IntegrationSource: entity1.IntegrationSource,
		},
	)
	require.NoError(t, err)
	require.NotEmpty(t, entity2)

	// Ensure the fetched entity matches the one we just created
	require.Equal(t, entity1.EntityID, entity2.EntityID)
	require.Equal(t, entity1.Name, entity2.Name)
	require.Equal(t, entity1.Description, entity2.Description)
	require.Equal(t, entity1.IntegrationSource, entity2.IntegrationSource)

	t.Logf("Found entity: %s with integration source: %s", entity1.Name, entity1.IntegrationSource)
}

func TestGetEntitiesByNames(t *testing.T) {
	var names []string

	for i := 0; i < 5; i++ {
		entity := CreateEntity(t, true)
		names = append(names, entity.Name)
	}

	perm := rand.Perm(len(names))
	// Select two random names from the created list.
	selectedNames := []string{names[perm[0]], names[perm[1]]}

	entities, err := testQueries.GetEntitiesByNames(context.Background(), selectedNames)
	require.NoError(t, err)
	require.Len(t, entities, 2)

	for _, entity := range entities {
		require.Contains(t, selectedNames, entity.Name)
		t.Logf("Found entity: %s with integration source: %s", entity.Name, entity.IntegrationSource)
	}
}

func TestListEntities(t *testing.T) {
	// List entities with a limit of 5
	limit := int32(5)
	offset := int32(0)
	params := ListEntitiesParams{
		Limit:  limit,
		Offset: offset,
	}

	entities, err := testQueries.ListEntities(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, entities)
	require.Len(t, entities, int(limit))

	// Verify that the listed entities are not empty and have valid data
	entityList := []string{}
	for _, entity := range entities {
		require.NotEmpty(t, entity)
		require.NotZero(t, entity.EntityID)
		require.NotEmpty(t, entity.Name)
		require.NotEmpty(t, entity.Description)
		entityList = append(entityList, entity.Name)
		t.Logf("Found entity: %s with integration source: %s", entity.Name, entity.IntegrationSource)
	}

	t.Log("Entity list:", entityList)
}

func TestUpdateEntity(t *testing.T) {
	entity1 := CreateEntity(t, true)

	params := UpdateEntityByNameParams{
		Name:        entity1.Name,
		Name_2:      FormatName(gofakeit.Name()),
		Description: gofakeit.Sentence(5),
	}

	// Call UpdateEntityByName
	entity2, err := testQueries.UpdateEntityByName(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, entity2)

	// Verify the updated entity
	require.Equal(t, entity1.EntityID, entity2.EntityID)
	require.Equal(t, params.Name_2, entity2.Name)             // Updated name
	require.Equal(t, params.Description, entity2.Description) // Updated description

	t.Logf("Update entity %s to %s", params.Name, params.Name_2)
}

func TestUpdateEntityIntegrationSourceByNameAndSource(t *testing.T) {
	entity := CreateEntity(t, true)
	var newIntegrationSource string

	switch entity.IntegrationSource {
	case "telemetry", "gps":
		newIntegrationSource = "self_reported"
	}

	params := UpdateEntityIntegrationSourceByNameAndSourceParams{
		Name:                entity.Name,
		IntegrationSource:   entity.IntegrationSource,
		IntegrationSource_2: newIntegrationSource,
	}

	updatedEntity, err := testQueries.UpdateEntityIntegrationSourceByNameAndSource(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, updatedEntity)

	require.Equal(t, entity.EntityID, updatedEntity.EntityID)
	require.Equal(t, newIntegrationSource, updatedEntity.IntegrationSource)

	t.Logf("Updated entity %s integration source from %s to %s", entity.Name, entity.IntegrationSource, updatedEntity.IntegrationSource)
}

func TestDeleteRandomEntity(t *testing.T) {
	limit := int32(5)
	offset := int32(0)
	params := ListEntitiesParams{
		Limit:  limit,
		Offset: offset,
	}

	entities, err := testQueries.ListEntities(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, entities, "No entities available for deletion")

	randomIndex := rand.Intn(len(entities))
	randomEntity := entities[randomIndex]

	// Delete the selected entity.
	err = testQueries.DeleteEntity(context.Background(), randomEntity.EntityID)
	require.NoError(t, err)

	// Try to retrieve the deleted entity.
	deletedEntity, err := testQueries.GetEntity(context.Background(), randomEntity.EntityID)
	require.Error(t, err)
	require.Empty(t, deletedEntity)

	t.Log("Entity deleted:", randomEntity.Name)
}

func TestCreateAndUpdateEntityWithInstanceAndPosition(t *testing.T) {
	// Add position data to an existing entity and create an entity with position
	for _, toggle := range []bool{true, false} {
		getEntity := CreateEntity(t, toggle)

		entity, err := testQueries.GetEntityByNameAndIntegrationSource(context.Background(), GetEntityByNameAndIntegrationSourceParams{
			Name:              getEntity.Name,
			IntegrationSource: getEntity.IntegrationSource,
		})
		if err != nil && err != sql.ErrNoRows {
			t.Fatalf("Failed to get entity by name and integration source: %v", err)
		}

		var entityID uuid.UUID
		if err == sql.ErrNoRows {
			// Create a new entity if it doesn’t exist
			newEntity, err := testQueries.CreateEntity(context.Background(), CreateEntityParams{
				Name:              getEntity.Name,
				Description:       getEntity.Description,
				DataType:          sql.NullString{Valid: false},
				SourceName:        sql.NullString{Valid: false},
				IntegrationSource: getEntity.IntegrationSource,
			})
			if err != nil {
				t.Fatalf("Failed to create new entity: %v", err)
			}
			entityID = newEntity.EntityID
		} else {
			// If entity exists, use its ID
			entityID = entity.EntityID
		}

		// Insert location
		locationID, err := testQueries.InsertInstance(context.Background(), InsertInstanceParams{
			EntityID:  entityID,
			CreatedAt: time.Now().UTC(),
		})
		if err != nil {
			t.Fatalf("Failed to insert location: %v", err)
		}

		// Hardcoded position values
		position := InsertPositionParams{
			InstanceID:        locationID,
			LatitudeDegrees:   36.7749,
			LongitudeDegrees:  -123.4194,
			HeadingDegrees:    sql.NullFloat64{Float64: 90.0, Valid: true},
			AltitudeHaeMeters: sql.NullFloat64{Float64: 100.0, Valid: true},
			SpeedMps:          sql.NullFloat64{Float64: 22.0, Valid: true},
		}

		// Insert position
		err = testQueries.InsertPosition(context.Background(), position)
		if err != nil {
			t.Fatalf("Failed to insert position: %v", err)
		}

		t.Log("Instance and position data created for", entity.Name)
	}
}
