-- Drop foreign keys first
ALTER TABLE IF EXISTS "ontology" DROP CONSTRAINT IF EXISTS "ontology_entity_id_fkey";
ALTER TABLE IF EXISTS "provenance" DROP CONSTRAINT IF EXISTS "provenance_entity_id_fkey";
ALTER TABLE IF EXISTS "geo_detail" DROP CONSTRAINT IF EXISTS "geo_detail_entity_id_fkey";
ALTER TABLE IF EXISTS "position" DROP CONSTRAINT IF EXISTS "position_location_id_fkey";
ALTER TABLE IF EXISTS "location" DROP CONSTRAINT IF EXISTS "location_entity_id_fkey";
ALTER TABLE IF EXISTS "context" DROP CONSTRAINT IF EXISTS "context_entity_id_fkey";

-- Drop indexes on `context`
DROP INDEX IF EXISTS "context_entity_type_specific_type_idx";
DROP INDEX IF EXISTS "context_entity_type_idx";
DROP INDEX IF EXISTS "context_entity_id_idx";

-- Drop `context` table
DROP TABLE IF EXISTS "context" CASCADE;

-- Drop indexes
DROP INDEX IF EXISTS "ontology_entity_type_specific_type_idx";
DROP INDEX IF EXISTS "ontology_entity_type_idx";
DROP INDEX IF EXISTS "ontology_entity_id_idx";

DROP INDEX IF EXISTS "provenance_source_update_time_idx";
DROP INDEX IF EXISTS "provenance_integration_name_idx";
DROP INDEX IF EXISTS "provenance_entity_id_idx";

DROP INDEX IF EXISTS "geo_detail_geo_polygon_idx";
DROP INDEX IF EXISTS "geo_detail_geo_point_idx";
DROP INDEX IF EXISTS "geo_detail_entity_id_idx";

DROP INDEX IF EXISTS "position_latitude_degrees_longitude_degrees_idx";
DROP INDEX IF EXISTS "position_location_id_idx";

DROP INDEX IF EXISTS "location_created_at_modified_at_idx";
DROP INDEX IF EXISTS "location_entity_id_idx";

DROP INDEX IF EXISTS "entity_name_idx";

-- Drop tables (with CASCADE to remove dependencies)
DROP TABLE IF EXISTS "ontology" CASCADE;
DROP TABLE IF EXISTS "provenance" CASCADE;
DROP TABLE IF EXISTS "geo_detail" CASCADE;
DROP TABLE IF EXISTS "position" CASCADE;
DROP TABLE IF EXISTS "location" CASCADE;
DROP TABLE IF EXISTS "entity" CASCADE;

-- Drop sequences (if any exist)
DROP SEQUENCE IF EXISTS ontology_id_seq CASCADE;
DROP SEQUENCE IF EXISTS provenance_id_seq CASCADE;
DROP SEQUENCE IF EXISTS geo_detail_id_seq CASCADE;
DROP SEQUENCE IF EXISTS position_id_seq CASCADE;
DROP SEQUENCE IF EXISTS location_id_seq CASCADE;
DROP SEQUENCE IF EXISTS entity_id_seq CASCADE;

-- Drop extensions
DROP EXTENSION IF EXISTS postgis CASCADE;

-- End of schema teardown
