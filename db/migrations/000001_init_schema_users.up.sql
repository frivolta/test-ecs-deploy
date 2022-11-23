CREATE TYPE "role" AS ENUM (
    'TEACHER',
    'ADMIN',
    'PARENT'
    );

CREATE TABLE "users"
(
    "id"         bigserial PRIMARY KEY,
    "full_name"  varchar        NOT NULL,
    "role"       role           NOT NULL DEFAULT ('TEACHER'),
    "email"      varchar UNIQUE NOT NULL,
    "updated_at" timestamptz    NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("role");

CREATE UNIQUE INDEX ON "users" ("email");
