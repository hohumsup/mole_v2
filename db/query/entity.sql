-- name: CreateEntity :one
INSERT INTO entity (name, description)
VALUES ($1, $2) RETURNING *;

-- name: GetEntity :one
SELECT * FROM entity
WHERE entity_id = $1;

-- name: GetEntityByName :one
SELECT * FROM entity
WHERE name = $1;

-- name: GetEntitiesByNames :many
SELECT * FROM entity
WHERE name = ANY($1::text[]); -- ensures that the input parameter is explicitly cast as a PostgreSQL array of text

-- name: ListEntities :many
SELECT * FROM entity
ORDER BY entity_id
LIMIT $1 OFFSET $2;

-- name: UpdateEntityByName :one
UPDATE entity
SET
  name = $2,
  description = $3
WHERE name = $1
RETURNING entity_id, name, description;

-- name: DeleteEntity :exec
DELETE FROM entity
WHERE entity_id = $1;

-- name: CreateEntityWithPosition :one
WITH entity_inserted AS (
  -- Insert the entity if it doesn't already exist
  INSERT INTO entity (entity_id, name, description)
  VALUES (uuid_generate_v4(), $1, $2)
  ON CONFLICT (name) DO NOTHING -- Entity exists
  RETURNING entity_id
),
entity_selected AS (
  -- Retrieve the entity_id, whether newly inserted or already existing
  SELECT entity_id
  FROM entity
  WHERE name = $1
),
location_inserted AS (
  -- Insert the location for the entity
  INSERT INTO location (entity_id, created_at, modified_at)
  SELECT entity_id, now(), now()
  FROM entity_selected
  RETURNING id AS location_id
),
position_inserted AS (
  -- Insert the position data for the location
  INSERT INTO position (location_id, latitude_degrees, longitude_degrees, heading_degrees, altitude_hae_meters, speed_mps)
  SELECT location_id, $3, $4, $5, $6, $7
  FROM location_inserted
  RETURNING id AS position_id
)
-- Return the IDs of the created/selected entity, location, and position
SELECT
  es.entity_id,
  li.location_id,
  pi.position_id
FROM entity_selected es
LEFT JOIN location_inserted li ON TRUE
LEFT JOIN position_inserted pi ON TRUE;