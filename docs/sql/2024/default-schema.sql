DROP TABLE IF EXISTS users;
CREATE TABLE IF NOT EXISTS users
(
    "id"            SERIAL PRIMARY KEY,
    "fk_role_id"    INT          NOT NULL,
    "name"          VARCHAR(255) NOT NULL,
    "email"         VARCHAR(255) NOT NULL,
    "password"      VARCHAR(255) NOT NULL,
    "refresh_token" VARCHAR(255),

    -- Utility columns
    "status"        SMALLINT     NOT NULL DEFAULT 1,
    "flag"          INT          NOT NULL DEFAULT 0,
    "meta"          VARCHAR(255),
    "created_at"    TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by"    VARCHAR(255),
    "updated_at"    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_by"    VARCHAR(255),
    "deleted_at"    TIMESTAMPTZ,
    "deleted_by"    VARCHAR(255)
);

DROP TABLE IF EXISTS roles;
CREATE TABLE IF NOT EXISTS roles
(
    "id"         SERIAL PRIMARY KEY,
    "role"       VARCHAR(255) NOT NULL,

    -- Utility columns
    "status"     SMALLINT     NOT NULL DEFAULT 1,
    "flag"       INT          NOT NULL DEFAULT 0,
    "meta"       VARCHAR(255),
    "created_at" TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" VARCHAR(255),
    "updated_at" TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_by" VARCHAR(255),
    "deleted_at" TIMESTAMPTZ,
    "deleted_by" VARCHAR(255)
);

INSERT INTO users (fk_role_id, name, email, password, status)
VALUES (1, 'admin', 'admin@test.com', '$2a$10$T3jgYJ5EPO6m8G92NoLb4e8NhSBLUkipgfbPlMbmU4nWesRk1g9wi', 1);


INSERT INTO roles ("id", "role")
VALUES (1, 'admin'),
       (2, 'user');