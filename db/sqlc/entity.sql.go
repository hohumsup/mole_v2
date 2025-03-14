// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: entity.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sqlc-dev/pqtype"
)

const createEntity = `-- name: CreateEntity :one

WITH new_entity AS (
  INSERT INTO entity (name, description)
  VALUES ($1, $2)
  RETURNING entity_id, name, description
),
new_provenance AS (
  INSERT INTO provenance (entity_id, data_type, integration_source, source_update_time)
  SELECT entity_id, $3, $4, now()
  FROM new_entity
  RETURNING entity_id, integration_source
),
new_context AS (
  INSERT INTO context (entity_id, template, entity_type, specific_type, created_at)
  SELECT 
    entity_id, 
    $5,
    $6,
    $7,
    now()
  FROM new_entity
  RETURNING entity_id, template
)
SELECT e.entity_id, e.name, e.description, p.integration_source, c.template
FROM new_entity e
JOIN new_provenance p ON e.entity_id = p.entity_id
JOIN new_context c ON e.entity_id = c.entity_id
`

type CreateEntityParams struct {
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	DataType          sql.NullString `json:"data_type"`
	IntegrationSource string         `json:"integration_source"`
	Template          int32          `json:"template"`
	EntityType        sql.NullString `json:"entity_type"`
	SpecificType      sql.NullString `json:"specific_type"`
}

type CreateEntityRow struct {
	EntityID          uuid.UUID `json:"entity_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	IntegrationSource string    `json:"integration_source"`
	Template          int32     `json:"template"`
}

// ----------------------------------------------------
// Entity Queries
// ----------------------------------------------------
// Description: Create a new entity
func (q *Queries) CreateEntity(ctx context.Context, arg CreateEntityParams) (CreateEntityRow, error) {
	row := q.db.QueryRowContext(ctx, createEntity,
		arg.Name,
		arg.Description,
		arg.DataType,
		arg.IntegrationSource,
		arg.Template,
		arg.EntityType,
		arg.SpecificType,
	)
	var i CreateEntityRow
	err := row.Scan(
		&i.EntityID,
		&i.Name,
		&i.Description,
		&i.IntegrationSource,
		&i.Template,
	)
	return i, err
}

const deleteEntity = `-- name: DeleteEntity :exec
DELETE FROM entity
WHERE entity_id = $1
`

// Description: Delete an entity by ID
func (q *Queries) DeleteEntity(ctx context.Context, entityID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteEntity, entityID)
	return err
}

const getEntitiesByNames = `-- name: GetEntitiesByNames :many
SELECT e.entity_id, e.name, e.description, p.integration_source
FROM entity e
JOIN provenance p on e.entity_id = p.entity_id
WHERE name = ANY($1::text[])
`

type GetEntitiesByNamesRow struct {
	EntityID          uuid.UUID `json:"entity_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	IntegrationSource string    `json:"integration_source"`
}

// Description: Retrieve entities with an array of names
func (q *Queries) GetEntitiesByNames(ctx context.Context, dollar_1 []string) ([]GetEntitiesByNamesRow, error) {
	rows, err := q.db.QueryContext(ctx, getEntitiesByNames, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetEntitiesByNamesRow{}
	for rows.Next() {
		var i GetEntitiesByNamesRow
		if err := rows.Scan(
			&i.EntityID,
			&i.Name,
			&i.Description,
			&i.IntegrationSource,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEntity = `-- name: GetEntity :one
SELECT entity_id, name, description FROM entity
WHERE entity_id = $1
`

// Description: Retrieve an entity by ID
func (q *Queries) GetEntity(ctx context.Context, entityID uuid.UUID) (Entity, error) {
	row := q.db.QueryRowContext(ctx, getEntity, entityID)
	var i Entity
	err := row.Scan(&i.EntityID, &i.Name, &i.Description)
	return i, err
}

const getEntityByNameAndIntegrationSource = `-- name: GetEntityByNameAndIntegrationSource :one
SELECT 
    e.entity_id, 
    e.name, 
    e.description, 
    p.integration_source,
    c.template
FROM entity e
JOIN provenance p ON e.entity_id = p.entity_id
JOIN context c ON e.entity_id = c.entity_id
WHERE e.name = $1 AND p.integration_source = $2
LIMIT 1
`

type GetEntityByNameAndIntegrationSourceParams struct {
	Name              string `json:"name"`
	IntegrationSource string `json:"integration_source"`
}

type GetEntityByNameAndIntegrationSourceRow struct {
	EntityID          uuid.UUID `json:"entity_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	IntegrationSource string    `json:"integration_source"`
	Template          int32     `json:"template"`
}

// Description: Retrieve an entity by name and integration source
func (q *Queries) GetEntityByNameAndIntegrationSource(ctx context.Context, arg GetEntityByNameAndIntegrationSourceParams) (GetEntityByNameAndIntegrationSourceRow, error) {
	row := q.db.QueryRowContext(ctx, getEntityByNameAndIntegrationSource, arg.Name, arg.IntegrationSource)
	var i GetEntityByNameAndIntegrationSourceRow
	err := row.Scan(
		&i.EntityID,
		&i.Name,
		&i.Description,
		&i.IntegrationSource,
		&i.Template,
	)
	return i, err
}

const getEntityByNames = `-- name: GetEntityByNames :many
SELECT e.entity_id, e.name, e.description, p.integration_source
FROM entity e
JOIN provenance p on e.entity_id = p.entity_id
WHERE e.name = $1
`

type GetEntityByNamesRow struct {
	EntityID          uuid.UUID `json:"entity_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	IntegrationSource string    `json:"integration_source"`
}

// Description: Retrieve entities by name
func (q *Queries) GetEntityByNames(ctx context.Context, name string) ([]GetEntityByNamesRow, error) {
	rows, err := q.db.QueryContext(ctx, getEntityByNames, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetEntityByNamesRow{}
	for rows.Next() {
		var i GetEntityByNamesRow
		if err := rows.Scan(
			&i.EntityID,
			&i.Name,
			&i.Description,
			&i.IntegrationSource,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getInstances = `-- name: GetInstances :many
SELECT 
    e.entity_id,
    e.name AS entity_name,
    p.integration_source,
	  c.template,
    i.id, i.entity_id, i.produced_by, i.created_at, i.modified_at, i.metadata,
    pos.instance_id, pos.latitude_degrees, pos.longitude_degrees, pos.heading_degrees, pos.altitude_hae_meters, pos.speed_mps
FROM entity e
JOIN provenance p ON e.entity_id = p.entity_id
JOIN context c ON e.entity_id = c.entity_id
JOIN instance i ON e.entity_id = i.entity_id
JOIN position pos ON i.id = pos.instance_id
ORDER by i.created_at
`

type GetInstancesRow struct {
	EntityID          uuid.UUID             `json:"entity_id"`
	EntityName        string                `json:"entity_name"`
	IntegrationSource string                `json:"integration_source"`
	Template          int32                 `json:"template"`
	ID                uuid.UUID             `json:"id"`
	EntityID_2        uuid.UUID             `json:"entity_id_2"`
	ProducedBy        sql.NullString        `json:"produced_by"`
	CreatedAt         time.Time             `json:"created_at"`
	ModifiedAt        time.Time             `json:"modified_at"`
	Metadata          pqtype.NullRawMessage `json:"metadata"`
	InstanceID        uuid.UUID             `json:"instance_id"`
	LatitudeDegrees   float64               `json:"latitude_degrees"`
	LongitudeDegrees  float64               `json:"longitude_degrees"`
	HeadingDegrees    sql.NullFloat64       `json:"heading_degrees"`
	AltitudeHaeMeters sql.NullFloat64       `json:"altitude_hae_meters"`
	SpeedMps          sql.NullFloat64       `json:"speed_mps"`
}

func (q *Queries) GetInstances(ctx context.Context) ([]GetInstancesRow, error) {
	rows, err := q.db.QueryContext(ctx, getInstances)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetInstancesRow{}
	for rows.Next() {
		var i GetInstancesRow
		if err := rows.Scan(
			&i.EntityID,
			&i.EntityName,
			&i.IntegrationSource,
			&i.Template,
			&i.ID,
			&i.EntityID_2,
			&i.ProducedBy,
			&i.CreatedAt,
			&i.ModifiedAt,
			&i.Metadata,
			&i.InstanceID,
			&i.LatitudeDegrees,
			&i.LongitudeDegrees,
			&i.HeadingDegrees,
			&i.AltitudeHaeMeters,
			&i.SpeedMps,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertInstance = `-- name: InsertInstance :one

INSERT INTO instance (entity_id, produced_by, created_at, metadata)
VALUES ($1, $2, $3, $4)
RETURNING id
`

type InsertInstanceParams struct {
	EntityID   uuid.UUID             `json:"entity_id"`
	ProducedBy sql.NullString        `json:"produced_by"`
	CreatedAt  time.Time             `json:"created_at"`
	Metadata   pqtype.NullRawMessage `json:"metadata"`
}

// ----------------------------------------------------
// Instance / Position Queries
// ----------------------------------------------------
func (q *Queries) InsertInstance(ctx context.Context, arg InsertInstanceParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, insertInstance,
		arg.EntityID,
		arg.ProducedBy,
		arg.CreatedAt,
		arg.Metadata,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const insertPosition = `-- name: InsertPosition :exec
INSERT INTO position (instance_id, latitude_degrees, longitude_degrees, heading_degrees, altitude_hae_meters, speed_mps)
VALUES ($1, $2, $3, $4, $5, $6)
`

type InsertPositionParams struct {
	InstanceID        uuid.UUID       `json:"instance_id"`
	LatitudeDegrees   float64         `json:"latitude_degrees"`
	LongitudeDegrees  float64         `json:"longitude_degrees"`
	HeadingDegrees    sql.NullFloat64 `json:"heading_degrees"`
	AltitudeHaeMeters sql.NullFloat64 `json:"altitude_hae_meters"`
	SpeedMps          sql.NullFloat64 `json:"speed_mps"`
}

func (q *Queries) InsertPosition(ctx context.Context, arg InsertPositionParams) error {
	_, err := q.db.ExecContext(ctx, insertPosition,
		arg.InstanceID,
		arg.LatitudeDegrees,
		arg.LongitudeDegrees,
		arg.HeadingDegrees,
		arg.AltitudeHaeMeters,
		arg.SpeedMps,
	)
	return err
}

const listEntities = `-- name: ListEntities :many
SELECT e.entity_id, e.name, e.description, p.integration_source
FROM entity e
JOIN provenance p on e.entity_id = p.entity_id
ORDER BY e.entity_id
LIMIT $1 OFFSET $2
`

type ListEntitiesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListEntitiesRow struct {
	EntityID          uuid.UUID `json:"entity_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	IntegrationSource string    `json:"integration_source"`
}

// Description: Retrieve all entities
func (q *Queries) ListEntities(ctx context.Context, arg ListEntitiesParams) ([]ListEntitiesRow, error) {
	rows, err := q.db.QueryContext(ctx, listEntities, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListEntitiesRow{}
	for rows.Next() {
		var i ListEntitiesRow
		if err := rows.Scan(
			&i.EntityID,
			&i.Name,
			&i.Description,
			&i.IntegrationSource,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateEntityByName = `-- name: UpdateEntityByName :one
UPDATE entity
SET
  name = $2,
  description = $3
WHERE name = $1
RETURNING entity_id, name, description
`

type UpdateEntityByNameParams struct {
	Name        string `json:"name"`
	Name_2      string `json:"name_2"`
	Description string `json:"description"`
}

// Description: Update an entity by name
func (q *Queries) UpdateEntityByName(ctx context.Context, arg UpdateEntityByNameParams) (Entity, error) {
	row := q.db.QueryRowContext(ctx, updateEntityByName, arg.Name, arg.Name_2, arg.Description)
	var i Entity
	err := row.Scan(&i.EntityID, &i.Name, &i.Description)
	return i, err
}

const updateEntityIntegrationSourceByNameAndSource = `-- name: UpdateEntityIntegrationSourceByNameAndSource :one
UPDATE provenance
SET integration_source = $3
WHERE entity_id = (
  SELECT entity_id FROM entity WHERE name = $1
)
AND integration_source = $2
RETURNING entity_id, integration_source
`

type UpdateEntityIntegrationSourceByNameAndSourceParams struct {
	Name                string `json:"name"`
	IntegrationSource   string `json:"integration_source"`
	IntegrationSource_2 string `json:"integration_source_2"`
}

type UpdateEntityIntegrationSourceByNameAndSourceRow struct {
	EntityID          uuid.UUID `json:"entity_id"`
	IntegrationSource string    `json:"integration_source"`
}

// Description: Update an entity's integration source by name and source
// TODO: Add Template to UpdateEntityIntegrationSourceByNameAndSource
func (q *Queries) UpdateEntityIntegrationSourceByNameAndSource(ctx context.Context, arg UpdateEntityIntegrationSourceByNameAndSourceParams) (UpdateEntityIntegrationSourceByNameAndSourceRow, error) {
	row := q.db.QueryRowContext(ctx, updateEntityIntegrationSourceByNameAndSource, arg.Name, arg.IntegrationSource, arg.IntegrationSource_2)
	var i UpdateEntityIntegrationSourceByNameAndSourceRow
	err := row.Scan(&i.EntityID, &i.IntegrationSource)
	return i, err
}
