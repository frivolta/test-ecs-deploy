CREATE TABLE "carnets"
(
    "id"       bigserial PRIMARY KEY NOT NULL,
    "date"     date                  NOT NULL,
    "quantity" int                   NOT NULL DEFAULT (0),
    "kid_id"   bigint                NOT NULL
);
ALTER TABLE "carnets"
    ADD FOREIGN KEY ("kid_id") REFERENCES "kids" ("id");
