-- name: CreateEntity :one
INSERT INTO entity (entity_id, name, description)
VALUES ($1, $2, $3) RETURNING *;

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