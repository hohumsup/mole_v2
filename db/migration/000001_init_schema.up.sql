-- Enable necessary extensions
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-------------------------------------------------
-- ENTITY TABLE
-------------------------------------------------
CREATE TABLE "entity" (
  "entity_id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "name" varchar NOT NULL,
  "description" text NOT NULL DEFAULT 'Human-readable description'
);

COMMENT ON COLUMN "entity"."entity_id" IS 'Unique identifier for the entity';
COMMENT ON COLUMN "entity"."name" IS 'Display name for the entity';
COMMENT ON COLUMN "entity"."description" IS 'Human-readable description';

-------------------------------------------------
-- INSTANCE TABLE
-------------------------------------------------
CREATE TABLE "instance" (
  "id" bigserial PRIMARY KEY,
  "entity_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "modified_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX ON "instance" ("entity_id");
CREATE INDEX ON "instance" ("created_at", "modified_at");

ALTER TABLE "instance" 
  ADD FOREIGN KEY ("entity_id") 
    REFERENCES "entity" ("entity_id") 
    ON DELETE CASCADE;

COMMENT ON COLUMN "instance"."id" IS 'Unique ID for the instance record';
COMMENT ON COLUMN "instance"."entity_id" IS 'Reference to the associated entity';
COMMENT ON COLUMN "instance"."created_at" IS 'Client-provided timestamp for when the instance data was created.';
COMMENT ON COLUMN "instance"."modified_at" IS 'Server-generated timestamp for last modification.';

-------------------------------------------------
-- POSITION TABLE
-------------------------------------------------
CREATE TABLE "position" (
  "position_id" bigserial PRIMARY KEY,
  "instance_id" bigint NOT NULL,
  "latitude_degrees" double precision NOT NULL,
  "longitude_degrees" double precision NOT NULL,
  "heading_degrees" double precision,
  "altitude_hae_meters" double precision,
  "speed_mps" double precision
);

CREATE INDEX ON "position" ("instance_id");
CREATE INDEX ON "position" ("latitude_degrees", "longitude_degrees");

ALTER TABLE "position" 
  ADD FOREIGN KEY ("instance_id") 
    REFERENCES "instance" ("id") 
    ON DELETE CASCADE;

COMMENT ON COLUMN "position"."position_id" IS 'Unique ID for the position record';
COMMENT ON COLUMN "position"."instance_id" IS 'Reference to the associated instance';
COMMENT ON COLUMN "position"."latitude_degrees" IS 'WGS84 geodetic latitude in decimal degrees.';
COMMENT ON COLUMN "position"."longitude_degrees" IS 'WGS84 longitude in decimal degrees.';
COMMENT ON COLUMN "position"."heading_degrees" IS 'Heading in degrees.';
COMMENT ON COLUMN "position"."altitude_hae_meters" IS 'Altitude as height above ellipsoid (WGS84), in meters.';
COMMENT ON COLUMN "position"."speed_mps" IS 'Speed as the magnitude of velocity, in meters per second.';

-------------------------------------------------
-- GEO_DETAIL TABLE
-------------------------------------------------
CREATE TABLE "geo_detail" (
  "geo_id" bigserial PRIMARY KEY,
  "instance_id" bigint NOT NULL,
  "geo_point" geometry,
  "geo_line" geometry,
  "geo_polygon" geometry,
  "geo_ellipse" geometry,
  "geo_ellipsoid" geometry
);

CREATE INDEX ON "geo_detail" ("instance_id");
CREATE INDEX ON "geo_detail" ("geo_point");
CREATE INDEX ON "geo_detail" ("geo_polygon");

ALTER TABLE "geo_detail" 
  ADD FOREIGN KEY ("instance_id") 
    REFERENCES "instance" ("id") 
    ON DELETE CASCADE;

COMMENT ON COLUMN "geo_detail"."geo_id" IS 'Unique ID for the geo detail record';
COMMENT ON COLUMN "geo_detail"."instance_id" IS 'Reference to the associated instance';
COMMENT ON COLUMN "geo_detail"."geo_point" IS 'Geospatial point representation of the entity.';
COMMENT ON COLUMN "geo_detail"."geo_line" IS 'Geospatial line representation of the entity.';
COMMENT ON COLUMN "geo_detail"."geo_polygon" IS 'Geospatial polygon representation of the entity.';
COMMENT ON COLUMN "geo_detail"."geo_ellipse" IS 'Geospatial ellipse representation of the entity.';
COMMENT ON COLUMN "geo_detail"."geo_ellipsoid" IS 'Geospatial ellipsoid representation of the entity.';

-------------------------------------------------
-- PROVENANCE TABLE
-------------------------------------------------
CREATE TABLE "provenance" (
  "id" bigserial PRIMARY KEY,
  "entity_id" uuid NOT NULL,
  "data_type" varchar,
  "source_name" varchar, 
  "integration_source" varchar NOT NULL,
  "source_update_time" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE "provenance" 
  ADD FOREIGN KEY ("entity_id") 
    REFERENCES "entity" ("entity_id") 
    ON DELETE CASCADE;
  
CREATE INDEX ON "provenance" ("entity_id");
CREATE INDEX ON "provenance" ("source_update_time");

COMMENT ON COLUMN "provenance"."id" IS 'Unique ID for the provenance record';
COMMENT ON COLUMN "provenance"."entity_id" IS 'Reference to the entity associated with this provenance record';
COMMENT ON COLUMN "provenance"."data_type" IS 'Optional name or identifier for the source system (e.g., ''gps'', ''telemetry'')';
COMMENT ON COLUMN "provenance"."source_name" IS 'Optional name for the entity that generated this entity. Used for events such as detections or tracks';
COMMENT ON COLUMN "provenance"."integration_source" IS 'Integration source used for which system provided the data (e.g., ''TAK'')';
COMMENT ON COLUMN "provenance"."source_update_time" IS 'Last modification time according to the source system';
COMMENT ON COLUMN "provenance"."created_at" IS 'Timestamp for when the provenance record was created';

-------------------------------------------------
-- TEMPLATE TABLE
-------------------------------------------------

CREATE TABLE template (
  template_id serial PRIMARY KEY,
  name varchar NOT NULL UNIQUE
);

INSERT INTO template (template_id, name) VALUES
  (1, 'Event'),
  (2, 'Asset'),
  (3, 'Geo');

-------------------------------------------------
-- CONTEXT TABLE
-------------------------------------------------

CREATE TABLE "context" (
  "context_id" bigserial PRIMARY KEY,
  "entity_id" uuid NOT NULL,
  "template" integer NOT NULL REFERENCES template(template_id),
  "entity_type" varchar,
  "specific_type" varchar,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT unique_entity_context UNIQUE (entity_id)
);

ALTER TABLE "context" 
  ADD FOREIGN KEY ("entity_id") 
    REFERENCES "entity" ("entity_id") 
    ON DELETE CASCADE;

CREATE INDEX ON "context" ("entity_id");
CREATE INDEX ON "context" ("entity_type");


COMMENT ON COLUMN "context"."context_id" IS 'Unique ID for the context record';
COMMENT ON COLUMN "context"."entity_id" IS 'Reference to the associated entity';
COMMENT ON COLUMN "context"."template" IS 'Template type for the context record';
COMMENT ON COLUMN "context"."entity_type" IS 'High-level classification (e.g., ''detection'', ''uav'')';
COMMENT ON COLUMN "context"."specific_type" IS 'A detailed categorization (e.g., ''radar'', ''fixed-wing'')';
COMMENT ON COLUMN "context"."created_at" IS 'Timestamp for when the context record was created';

-------------------------------------------------
-- TRIGGER FUNCTION TO ENFORCE UNIQUE (entity.name, integration_source)
-------------------------------------------------
CREATE OR REPLACE FUNCTION check_unique_entity_name_integration_func()
RETURNS trigger AS $$
DECLARE
    current_name varchar;
    duplicate_count integer;
BEGIN
    -- Get the entity's name for the given entity_id from the new provenance record.
    SELECT name INTO current_name FROM entity WHERE entity_id = NEW.entity_id;
    
    -- Count any existing provenance records (joined with their entities) that share the same entity name
    -- and the same integration_source, but with a different entity_id.
    SELECT count(*) INTO duplicate_count
    FROM entity e
    JOIN provenance p ON e.entity_id = p.entity_id
    WHERE e.name = current_name
      AND p.integration_source = NEW.integration_source
      AND e.entity_id <> NEW.entity_id;
      
    IF duplicate_count > 0 THEN
        RAISE EXCEPTION 'Entity with name "%" and integration_source "%" already exists', current_name, NEW.integration_source;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-------------------------------------------------
-- TRIGGER ON PROVENANCE TABLE TO ENFORCE UNIQUENESS OF (entity.name, integration_source)
-------------------------------------------------
CREATE TRIGGER check_unique_entity_name_integration_trigger
BEFORE INSERT OR UPDATE ON provenance
FOR EACH ROW
EXECUTE FUNCTION check_unique_entity_name_integration_func();