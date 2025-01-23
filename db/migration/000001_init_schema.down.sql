-- Drop foreign keys first
ALTER TABLE "ontology" DROP CONSTRAINT "ontology_entity_id_fkey";
ALTER TABLE "provenance" DROP CONSTRAINT "provenance_entity_id_fkey";
ALTER TABLE "geo_detail" DROP CONSTRAINT "geo_detail_entity_id_fkey";
ALTER TABLE "position" DROP CONSTRAINT "position_location_id_fkey";
ALTER TABLE "location" DROP CONSTRAINT "location_entity_id_fkey";

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

-- Drop tables
DROP TABLE IF EXISTS "ontology";
DROP TABLE IF EXISTS "provenance";
DROP TABLE IF EXISTS "geo_detail";
DROP TABLE IF EXISTS "position";
DROP TABLE IF EXISTS "location";
DROP TABLE IF EXISTS "entity";

-- Drop extensions
DROP EXTENSION IF EXISTS postgis CASCADE;

-- End of schema teardown
