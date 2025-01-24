CREATE EXTENSION postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "entity" (
  "entity_id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" text NOT NULL DEFAULT 'Auto-generated description'
);

CREATE TABLE "location" (
  "id" bigserial PRIMARY KEY,
  "entity_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "modified_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "position" (
  "id" bigserial PRIMARY KEY,
  "location_id" bigint NOT NULL,
  "latitude_degrees" double precision NOT NULL,
  "longitude_degrees" double precision NOT NULL,
  "heading_degrees" double precision,
  "altitude_hae_meters" double precision,
  "speed_mps" double precision
);

CREATE TABLE "geo_detail" (
  "id" bigserial PRIMARY KEY,
  "entity_id" uuid NOT NULL,
  "geo_point" geometry,
  "geo_line" geometry,
  "geo_polygon" geometry,
  "geo_ellipse" geometry,
  "geo_ellipsoid" geometry,
  "created_at" timestamptz NOT NULL DEFAULT 'now()',
  "modified_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "provenance" (
  "id" bigserial PRIMARY KEY,
  "entity_id" uuid NOT NULL,
  "integration_name" varchar NOT NULL,
  "data_type" varchar,
  "source_name" varchar,
  "source_update_time" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "context" (
  "id" bigserial PRIMARY KEY,
  "entity_id" uuid NOT NULL,
  "entity_type" varchar NOT NULL,
  "specific_type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE UNIQUE INDEX ON "entity" ("name");

CREATE INDEX ON "location" ("entity_id");

CREATE INDEX ON "location" ("created_at", "modified_at");

CREATE INDEX ON "position" ("location_id");

CREATE INDEX ON "position" ("latitude_degrees", "longitude_degrees");

CREATE INDEX ON "geo_detail" ("entity_id");

CREATE INDEX ON "geo_detail" ("geo_point");

CREATE INDEX ON "geo_detail" ("geo_polygon");

CREATE INDEX ON "provenance" ("entity_id");

CREATE INDEX ON "provenance" ("integration_name");

CREATE INDEX ON "provenance" ("source_update_time");

CREATE INDEX ON "context" ("entity_id");

CREATE INDEX ON "context" ("entity_type");

CREATE UNIQUE INDEX ON "context" ("entity_type", "specific_type");

COMMENT ON COLUMN "entity"."entity_id" IS 'Unique identifier for the entity';

COMMENT ON COLUMN "entity"."name" IS 'Display name for the entity';

COMMENT ON COLUMN "entity"."description" IS 'Human-readable description';

COMMENT ON COLUMN "location"."id" IS 'Unique ID for the location record';

COMMENT ON COLUMN "location"."entity_id" IS 'Reference to the associated entity';

COMMENT ON COLUMN "location"."created_at" IS 'Client-provided timestamp for when the location data was created.';

COMMENT ON COLUMN "location"."modified_at" IS 'Server-generated timestamp for last modification.';

COMMENT ON COLUMN "position"."id" IS 'Unique ID for the position record';

COMMENT ON COLUMN "position"."location_id" IS 'Reference to the associated location';

COMMENT ON COLUMN "position"."latitude_degrees" IS 'WGS84 geodetic latitude in decimal degrees.';

COMMENT ON COLUMN "position"."longitude_degrees" IS 'WGS84 longitude in decimal degrees.';

COMMENT ON COLUMN "position"."altitude_hae_meters" IS 'Altitude as height above ellipsoid (WGS84), in meters.';

COMMENT ON COLUMN "position"."speed_mps" IS 'Speed as the magnitude of velocity, in meters per second.';

COMMENT ON COLUMN "geo_detail"."id" IS 'Unique ID for the geo detail record';

COMMENT ON COLUMN "geo_detail"."entity_id" IS 'Reference to the associated entity';

COMMENT ON COLUMN "geo_detail"."geo_point" IS 'Geospatial point representation of the entity.';

COMMENT ON COLUMN "geo_detail"."geo_line" IS 'Geospatial line representation of the entity.';

COMMENT ON COLUMN "geo_detail"."geo_polygon" IS 'Geospatial polygon representation of the entity.';

COMMENT ON COLUMN "geo_detail"."geo_ellipse" IS 'Geospatial ellipse representation of the entity.';

COMMENT ON COLUMN "geo_detail"."geo_ellipsoid" IS 'Geospatial ellipsoid representation of the entity.';

COMMENT ON COLUMN "geo_detail"."created_at" IS 'Timestamp for when the geo detail was created';

COMMENT ON COLUMN "geo_detail"."modified_at" IS 'Timestamp for when the geo detail was last updated';

COMMENT ON COLUMN "provenance"."id" IS 'Unique ID for the provenance record';

COMMENT ON COLUMN "provenance"."entity_id" IS 'Reference to the entity being tracked';

COMMENT ON COLUMN "provenance"."integration_name" IS 'Name of the system producing this data';

COMMENT ON COLUMN "provenance"."data_type" IS 'Type of the relationship or data (optional)';

COMMENT ON COLUMN "provenance"."source_name" IS 'Optional reference to an `entity_name`';

COMMENT ON COLUMN "provenance"."source_update_time" IS 'Last modification time according to the source system';

COMMENT ON COLUMN "provenance"."created_at" IS 'Timestamp for when the provenance record was created';

COMMENT ON COLUMN "context"."id" IS 'Unique ID for the context record';

COMMENT ON COLUMN "context"."entity_id" IS 'Reference to the associated entity';

COMMENT ON COLUMN "context"."entity_type" IS 'High-level classification (e.g., ''event-type'', ''vehicle'', ''sensor'')';

COMMENT ON COLUMN "context"."specific_type" IS 'A detailed categorization or model within the high-level classification (e.g., ''Detection'', ''Fixed-wing'', ''Weather-station'')';

COMMENT ON COLUMN "context"."created_at" IS 'Timestamp for when the context record was created';

ALTER TABLE "entity" ALTER COLUMN "entity_id" SET DEFAULT gen_random_uuid();

ALTER TABLE "location" ADD FOREIGN KEY ("entity_id") REFERENCES "entity" ("entity_id");

ALTER TABLE "position" ADD FOREIGN KEY ("location_id") REFERENCES "location" ("id");

ALTER TABLE "geo_detail" ADD FOREIGN KEY ("entity_id") REFERENCES "entity" ("entity_id");

ALTER TABLE "provenance" ADD FOREIGN KEY ("entity_id") REFERENCES "entity" ("entity_id");

ALTER TABLE "provenance" ADD COLUMN "location_id" bigint;

ALTER TABLE "provenance" ADD CONSTRAINT "unique_location_id" UNIQUE ("location_id");

ALTER TABLE "provenance" ADD CONSTRAINT "fk_location_id" FOREIGN KEY ("location_id") REFERENCES location (id);

ALTER TABLE "context" ADD FOREIGN KEY ("entity_id") REFERENCES "entity" ("entity_id");
