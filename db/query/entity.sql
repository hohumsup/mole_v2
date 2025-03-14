------------------------------------------------------
-- Entity Queries 
------------------------------------------------------

-- name: CreateEntity :one
-- Description: Create a new entity
WITH new_entity AS (
  INSERT INTO entity (name, description)
  VALUES ($1, $2)
  RETURNING *
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
JOIN new_context c ON e.entity_id = c.entity_id;

-- name: GetEntity :one
-- Description: Retrieve an entity by ID
SELECT * FROM entity
WHERE entity_id = $1;

-- name: GetEntityByNameAndIntegrationSource :one
-- Description: Retrieve an entity by name and integration source
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
-- Description: Retrieve entities by name
SELECT e.entity_id, e.name, e.description, p.integration_source
FROM entity e
JOIN provenance p on e.entity_id = p.entity_id
WHERE e.name = $1;

-- name: GetEntitiesByNames :many
-- Description: Retrieve entities with an array of names
SELECT e.entity_id, e.name, e.description, p.integration_source
FROM entity e
JOIN provenance p on e.entity_id = p.entity_id
WHERE name = ANY($1::text[]);

-- name: ListEntities :many
-- Description: Retrieve all entities
SELECT e.entity_id, e.name, e.description, p.integration_source
FROM entity e
JOIN provenance p on e.entity_id = p.entity_id
ORDER BY e.entity_id
LIMIT $1 OFFSET $2;

-- name: UpdateEntityByName :one
-- Description: Update an entity by name
UPDATE entity
SET
  name = $2,
  description = $3
WHERE name = $1
RETURNING entity_id, name, description;

-- name: UpdateEntityIntegrationSourceByNameAndSource :one
-- Description: Update an entity's integration source by name and source
-- TODO: Add Template to UpdateEntityIntegrationSourceByNameAndSource
UPDATE provenance
SET integration_source = $3
WHERE entity_id = (
  SELECT entity_id FROM entity WHERE name = $1
)
AND integration_source = $2
RETURNING entity_id, integration_source;

-- name: DeleteEntity :exec
-- Description: Delete an entity by ID
DELETE FROM entity
WHERE entity_id = $1;

------------------------------------------------------
-- Instance / Position Queries
------------------------------------------------------

-- name: InsertInstance :one
INSERT INTO instance (entity_id, produced_by, created_at, metadata)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: InsertPosition :exec
INSERT INTO position (instance_id, latitude_degrees, longitude_degrees, heading_degrees, altitude_hae_meters, speed_mps)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetInstances :many
SELECT 
    e.entity_id,
    e.name AS entity_name,
    p.integration_source,
	  c.template,
    i.*,
    pos.*
FROM entity e
JOIN provenance p ON e.entity_id = p.entity_id
JOIN context c ON e.entity_id = c.entity_id
JOIN instance i ON e.entity_id = i.entity_id
JOIN position pos ON i.id = pos.instance_id
ORDER by i.created_at;