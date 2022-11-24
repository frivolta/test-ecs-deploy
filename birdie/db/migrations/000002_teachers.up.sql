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


CREATE INDEX ON "teacher_notes" ("teacher_id");

CREATE INDEX ON "teacher_notes" ("date");

CREATE INDEX ON "teacher_notes" ("teacher_id", "date");

ALTER TABLE "teacher_notes" ADD FOREIGN KEY ("teacher_id") REFERENCES "teachers" ("id");
