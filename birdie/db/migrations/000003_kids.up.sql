CREATE TYPE "presence" AS ENUM (
    'MORNING',
    'AFTERNOON',
    'EVENING',
    'ABSENT'
    );

CREATE TABLE "kids"
(
    "id"      bigserial PRIMARY KEY NOT NULL,
    "name"    varchar               NOT NULL,
    "surname" varchar               NOT NULL
);

CREATE TABLE "kid_notes"
(
    "id"       bigserial PRIMARY KEY NOT NULL,
    "note"     varchar               NOT NULL,
    "kid_id"   bigint                NOT NULL,
    "presence" presence[]            NOT NULL,
    "has_meal"  bool                  NOT NULL DEFAULT (true),
    "date"     date                  NOT NULL
);

ALTER TABLE "kid_notes" ADD FOREIGN KEY ("kid_id") REFERENCES "kids" ("id");
