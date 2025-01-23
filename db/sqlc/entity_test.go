package db

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntity(t *testing.T) Entity {
	arg := CreateEntityParams{
		EntityID:    uuid.New(),
		Name:        gofakeit.Name(),
		Description: gofakeit.Sentence(5),
	}

	entity, err := testQueries.CreateEntity(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entity)

	require.Equal(t, arg.EntityID, entity.EntityID)
	require.Equal(t, arg.Name, entity.Name)
	require.Equal(t, arg.Description, entity.Description)

	return entity
}

func TestEntityGenerator(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomEntity(t)
	}
}

func TestGetEntity(t *testing.T) {
	entity1 := CreateRandomEntity(t)
	entity2, err := testQueries.GetEntity(context.Background(), entity1.EntityID)
	require.NoError(t, err)
	require.NotEmpty(t, entity2)

	require.Equal(t, entity1.EntityID, entity2.EntityID)
	require.Equal(t, entity1.Name, entity2.Name)
	require.Equal(t, entity1.Description, entity2.Description)
}

func TestGetEntityByName(t *testing.T) {
	entity1 := CreateRandomEntity(t)
	entity2, err := testQueries.GetEntityByName(context.Background(), entity1.Name)
	require.NoError(t, err)
	require.NotEmpty(t, entity2)

	require.Equal(t, entity1.EntityID, entity2.EntityID)
	require.Equal(t, entity1.Name, entity2.Name)
	require.Equal(t, entity1.Description, entity2.Description)
}

func TestGetEntitiesByNames(t *testing.T) {
	var createdEntities []Entity
	var names []string

	for i := 0; i < 5; i++ {
		entity := CreateRandomEntity(t)
		createdEntities = append(createdEntities, entity)
		names = append(names, entity.Name)
	}

	entities, err := testQueries.GetEntitiesByNames(context.Background(), names)
	require.NoError(t, err)
	require.NotEmpty(t, entities)
	require.Len(t, entities, len(createdEntities))

	for _, entity := range entities {
		found := false
		for _, createdEntity := range createdEntities {
			if entity.EntityID == createdEntity.EntityID {
				require.Equal(t, createdEntity.Name, entity.Name)
				require.Equal(t, createdEntity.Description, entity.Description)
				found = true
				break
			}
		}
		require.True(t, found)
	}
}

func TestListEntities(t *testing.T) {
	// Create multiple random entities
	for i := 0; i < 10; i++ {
		CreateRandomEntity(t)
	}

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
	for _, entity := range entities {
		require.NotEmpty(t, entity)
		require.NotZero(t, entity.EntityID)
		require.NotEmpty(t, entity.Name)
		require.NotEmpty(t, entity.Description)
	}
}

func TestUpdateEntity(t *testing.T) {
	entity1 := CreateRandomEntity(t)

	params := UpdateEntityByNameParams{
		Name:        entity1.Name,
		Name_2:      gofakeit.Name(),
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
}

func TestDeleteEntity(t *testing.T) {
	entity1 := CreateRandomEntity(t)

	err := testQueries.DeleteEntity(context.Background(), entity1.EntityID)
	require.NoError(t, err)

	entity2, err := testQueries.GetEntity(context.Background(), entity1.EntityID)
	require.Error(t, err)
	require.Empty(t, entity2)
}
