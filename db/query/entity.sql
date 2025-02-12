------------------------------------------------------
-- Entity Queries (Appended Below Existing Queries)
------------------------------------------------------

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
),
new_context AS (
  INSERT INTO context (entity_id, template, entity_type, specific_type, created_at)
  SELECT 
    entity_id, 
    $6,
    $7,
    $8,
    now()
  FROM new_entity
  RETURNING entity_id, template
)
SELECT e.entity_id, e.name, e.description, p.integration_source, c.template
FROM new_entity e
JOIN new_provenance p ON e.entity_id = p.entity_id
JOIN new_context c ON e.entity_id = c.entity_id;

-- name: GetEntity :one
SELECT * FROM entity
WHERE entity_id = $1;

-- name: GetEntityByNameAndIntegrationSource :one
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

------------------------------------------------------
-- Instance / Position Queries
------------------------------------------------------

-- name: InsertInstance :one
INSERT INTO instance (entity_id, created_at)
VALUES ($1, $2)
RETURNING id;

-- name: InsertPosition :exec
INSERT INTO position (instance_id, latitude_degrees, longitude_degrees, heading_degrees, altitude_hae_meters, speed_mps)
VALUES ($1, $2, $3, $4, $5, $6);
