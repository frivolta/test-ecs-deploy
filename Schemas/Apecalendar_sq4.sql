CREATE TYPE "role" AS ENUM (
  'TEACHER',
  'ADMIN',
  'PARENT'
);

CREATE TYPE "presence" AS ENUM (
  'MORNING',
  'AFTERNOON',
  'EVENING',
  'ABSENT'
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "role" role NOT NULL DEFAULT (TEACHER),
  "email" varchar UNIQUE NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "teachers" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "surname" varchar NOT NULL
);

CREATE TABLE "teacher_notes" (
  "id" bigserial PRIMARY KEY,
  "note" varchar NOT NULL,
  "teacher_id" bigint,
  "date" date NOT NULL
);

CREATE TABLE "kids" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "name" varchar NOT NULL,
  "surname" varchar NOT NULL
);

CREATE TABLE "kid_notes" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "note" varchar NOT NULL,
  "kid_id" bigint NOT NULL,
  "presence" presence[] NOT NULL,
  "hasMeal" bool NOT NULL DEFAULT (true),
  "date" date NOT NULL
);

CREATE TABLE "carnets" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "date" date NOT NULL,
  "quantity" int NOT NULL DEFAULT (0),
  "kid_id" bigint NOT NULL
);

CREATE INDEX ON "users" ("role");

CREATE UNIQUE INDEX ON "users" ("email");

CREATE INDEX ON "teacher_notes" ("teacher_id");

CREATE INDEX ON "teacher_notes" ("date");

CREATE INDEX ON "teacher_notes" ("teacher_id", "date");

ALTER TABLE "teacher_notes" ADD FOREIGN KEY ("teacher_id") REFERENCES "teachers" ("id");

ALTER TABLE "kid_notes" ADD FOREIGN KEY ("kid_id") REFERENCES "kids" ("id");

ALTER TABLE "carnets" ADD FOREIGN KEY ("kid_id") REFERENCES "kids" ("id");
