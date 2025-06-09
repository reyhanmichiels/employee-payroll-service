DROP TABLE IF EXISTS "user";
CREATE TABLE IF NOT EXISTS "user" (
    "id" SERIAL PRIMARY KEY,
    "fk_role_id" INT NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "refresh_token" VARCHAR(255),

    -- Utility columns
    "status" SMALLINT NOT NULL DEFAULT 1,
    "flag" INT NOT NULL DEFAULT 0,
    "meta" VARCHAR(255),
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" VARCHAR(255),
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_by" VARCHAR(255),
    "deleted_at" TIMESTAMP,
    "deleted_by" VARCHAR(255)
);

DROP TABLE IF EXISTS "role";
CREATE TABLE IF NOT EXISTS "role" (
    "id" SERIAL PRIMARY KEY,
    "role" VARCHAR(255) NOT NULL,

    -- Utility columns
    "status" SMALLINT NOT NULL DEFAULT 1,
    "flag" INT NOT NULL DEFAULT 0,
    "meta" VARCHAR(255),
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" VARCHAR(255),
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_by" VARCHAR(255),
    "deleted_at" TIMESTAMP,
    "deleted_by" VARCHAR(255)
);

INSERT INTO "role" ("role")
VALUES ('admin'),
       ('user');