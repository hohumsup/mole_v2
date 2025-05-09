-- name: BulkCreateEntities :exec
-- Description: One shot bulk insert of entities in a single db round trip 
WITH 
  -- Parse the input JSON array into individual rows
  data AS (
    SELECT
      data_item ->> 'entity_name'         AS entity_name,
      data_item ->> 'entity_description'  AS entity_description,
      data_item ->> 'data_type'           AS data_type,
      data_item ->> 'integration_source'  AS integration_source,
      (data_item ->> 'template')::int     AS template,
      data_item ->> 'entity_type'         AS entity_type,
      data_item ->> 'specific_type'       AS specific_type,
      data_item ->> 'instance_produced_by' AS instance_produced_by,
      data_item -> 'instance_metadata'     AS instance_metadata,
      (data_item ->> 'instance_created_at')::timestamptz AS instance_created_at,
      (data_item ->> 'latitude_degrees')::float  AS latitude_degrees,
      (data_item ->> 'longitude_degrees')::float AS longitude_degrees,
      (data_item ->> 'heading_degrees')::float   AS heading_degrees,
      (data_item ->> 'altitude_hae_meters')::float AS altitude_hae_meters,
      (data_item ->> 'speed_mps')::float         AS speed_mps
    FROM jsonb_array_elements($1::jsonb) AS data_item
  ),
  -- Group the payload to ensure that we have a single row per entity_name and integration_source
  grouped_payload AS (
    SELECT 
      entity_name, 
      MAX(entity_description) AS entity_description,
      integration_source,
      MIN(template) AS template,
      MIN(entity_type) AS entity_type, 
      MIN(specific_type) AS specific_type,
      MIN(data_type) AS data_type
    FROM data
    GROUP BY entity_name, integration_source
  ),
  -- Determine which entity_name and integration_source pairs exist in the database
  existing_entities AS (
    SELECT e.entity_id, e.name, p.integration_source
    FROM entity e
    JOIN provenance p ON e.entity_id = p.entity_id
    WHERE (e.name, p.integration_source) IN (
      SELECT DISTINCT entity_name, integration_source FROM data
    )
  ),
  -- Determine which entity_name and integration_source pairs do not yet exist
  to_insert_entities AS (
    SELECT gp.entity_name, gp.entity_description, gp.integration_source, gp.template, gp.entity_type, gp.specific_type, gp.data_type
    FROM grouped_payload gp
    WHERE NOT EXISTS (
      SELECT 1
      FROM entity e
      JOIN provenance p ON e.entity_id = p.entity_id
      WHERE e.name = gp.entity_name
        AND p.integration_source = gp.integration_source
    )
  ),
  -- Insert new entities for missing pairs
  inserted_entities AS (
    INSERT INTO entity (name, description)
    SELECT entity_name, entity_description
    FROM to_insert_entities
    RETURNING entity_id, name
  ),
  -- Associate the newly inserted entities with their integration source
  inserted_entities_with_source AS (
    SELECT ie.entity_id, ie.name, tie.integration_source
    FROM inserted_entities ie
    JOIN to_insert_entities tie ON tie.entity_name = ie.name
  ),
  -- Combine existing entities with the inserted entities
  all_entities_raw AS (
    SELECT * FROM existing_entities
    UNION ALL
    SELECT * FROM inserted_entities_with_source
  ),
  -- Aggregate the entity_id for each unique entity pair 
  all_entities AS (
    SELECT (array_agg(entity_id))[1] AS entity_id, name, integration_source
    FROM all_entities_raw
    GROUP BY name, integration_source
  ),
  -- Insert a single provence row per unique entity pair
  inserted_provenance AS (
    INSERT INTO provenance (entity_id, data_type, integration_source, source_update_time)
    SELECT ae.entity_id, MIN(d.data_type), d.integration_source, now()
    FROM all_entities ae
    JOIN data d ON ae.name = d.entity_name
                AND ae.integration_source = d.integration_source
    LEFT JOIN provenance p ON p.entity_id = ae.entity_id AND p.integration_source = d.integration_source
    WHERE p.entity_id IS NULL
    GROUP BY ae.entity_id, d.integration_source
    RETURNING entity_id, integration_source
  ),
  -- Insert a single context row per unique entity pair
  inserted_context AS (
    INSERT INTO context (entity_id, template, entity_type, specific_type, created_at)
    SELECT t.entity_id, t.template, t.entity_type, t.specific_type, now()
    FROM (
         SELECT ae.entity_id,
                MIN(d.template) AS template,
                MIN(d.entity_type) AS entity_type,
                MIN(d.specific_type) AS specific_type
         FROM all_entities ae
         JOIN data d ON ae.name = d.entity_name
                     AND ae.integration_source = d.integration_source
         GROUP BY ae.entity_id
    ) AS t
    ON CONFLICT (entity_id) DO NOTHING
    RETURNING entity_id, template
  ),
  -- Insert an instance row for every row in the bulk
  inserted_instances AS (
    INSERT INTO instance (entity_id, produced_by, metadata, created_at)
    SELECT ae.entity_id, d.instance_produced_by, d.instance_metadata, d.instance_created_at
    FROM all_entities ae
    JOIN data d ON ae.name = d.entity_name
                AND ae.integration_source = d.integration_source
    RETURNING instance_id, entity_id, created_at
  )
INSERT INTO position (
  instance_id, 
  instance_created_at, 
  latitude_degrees, 
  longitude_degrees, 
  heading_degrees, 
  altitude_hae_meters, 
  speed_mps
)
SELECT 
  ii.instance_id, 
  ii.created_at, 
  d.latitude_degrees, 
  d.longitude_degrees, 
  d.heading_degrees, 
  d.altitude_hae_meters, 
  d.speed_mps
FROM data d
JOIN all_entities ae 
  ON ae.name = d.entity_name
 AND ae.integration_source = d.integration_source
JOIN inserted_instances ii 
  ON ii.entity_id = ae.entity_id
WHERE d.latitude_degrees IS NOT NULL 
  AND d.longitude_degrees IS NOT NULL;
