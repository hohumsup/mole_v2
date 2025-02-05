-- name: CreateEntity :one
WITH new_entity AS (
  INSERT INTO entity (name, description)
  VALUES ($1, $2)
  RETURNING *
),
new_provenance AS (
  INSERT INTO provenance (entity_id, data_type, source_name, integration_source, source_update_time)
  SELECT entity_id, $3, $4, $5, now()
  FROM new_entity
  RETURNING entity_id, integration_source
)
SELECT e.entity_id, e.name, e.description, p.integration_source
FROM new_entity e
JOIN new_provenance p ON e.entity_id = p.entity_id;


-- name: GetEntity :one
SELECT * FROM entity
WHERE entity_id = $1;

-- name: GetEntityByNameAndIntegrationSource :one
SELECT 
    e.entity_id, 
    e.name, 
    e.description, 
    p.integration_source
FROM entity e
JOIN provenance p ON e.entity_id = p.entity_id
WHERE e.name = $1 AND p.integration_source = $2
LIMIT 1;

-- name: GetEntityByNames :many
SELECT e.entity_id, e.name, e.description, p.integration_source
FROM entity e
JOIN provenance p on e.entity_id = p.entity_id
where e.name = $1;

-- name: GetEntitiesByNames :many
SELECT e.entity_id, e.name, e.description, p.integration_source
FROM entity e
JOIN provenance p on e.entity_id = p.entity_id
WHERE name = ANY($1::text[]);

-- name: ListEntities :many
SELECT e.entity_id, e.name, e.description, p.integration_source
FROM entity e
JOIN provenance p on e.entity_id = p.entity_id
ORDER BY e.entity_id
LIMIT $1 OFFSET $2;

-- name: UpdateEntityByName :one
UPDATE entity
SET
  name = $2,
  description = $3
WHERE name = $1
RETURNING entity_id, name, description;

-- name: UpdateEntityIntegrationSourceByNameAndSource :one
UPDATE provenance
SET integration_source = $3
WHERE entity_id = (
  SELECT entity_id FROM entity WHERE name = $1
)
AND integration_source = $2
RETURNING entity_id, integration_source;

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