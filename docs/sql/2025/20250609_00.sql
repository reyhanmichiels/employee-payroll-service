DROP TYPE IF EXISTS period_status_enum;
CREATE TYPE period_status_enum AS ENUM ('UPCOMING', 'OPEN', 'CLOSED', 'PROCESSING', 'PROCESSED');

DROP TABLE IF EXISTS "attendance_periods";
CREATE TABLE IF NOT EXISTS "attendance_periods"
(
    "id"            SERIAL PRIMARY KEY,
    "start_date"    DATE               NOT NULL,
    "end_date"      DATE               NOT NULL,
    "period_status" period_status_enum NOT NULL DEFAULT 'UPCOMING',

    -- Utility columns
    "status"        SMALLINT           NOT NULL DEFAULT 1,
    "flag"          INT                NOT NULL DEFAULT 0,
    "meta"          VARCHAR(255),
    "created_at"    TIMESTAMPTZ        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by"    INT,
    "updated_at"    TIMESTAMPTZ        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_by"    INT,
    "deleted_at"    TIMESTAMPTZ,
    "deleted_by"    INT,

    CONSTRAINT payroll_periods_no_overlap
        EXCLUDE USING GIST (
        daterange(start_date, end_date, '[]') WITH &&)
);

DROP TABLE IF EXISTS "attendances";
CREATE TABLE IF NOT EXISTS "attendances"
(
    "id"                      SERIAL PRIMARY KEY,
    "fk_attendance_period_id" INT         NOT NULL,
    "fk_user_id"              INT         NOT NULL,
    "attendance_date"         DATE        NOT NULL,

    -- Utility columns
    "status"                  SMALLINT    NOT NULL DEFAULT 1,
    "flag"                    INT         NOT NULL DEFAULT 0,
    "meta"                    VARCHAR(255),
    "created_at"              TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by"              INT,
    "updated_at"              TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_by"              INT,
    "deleted_at"              TIMESTAMPTZ,
    "deleted_by"              INT,

    CONSTRAINT unique_user_date UNIQUE ("fk_user_id", "attendance_date")
);

DROP TABLE IF EXISTS "overtimes";
CREATE TABLE IF NOT EXISTS "overtimes"
(
    "id"            SERIAL PRIMARY KEY,
    "fk_user_id"    INT           NOT NULL,
    "overtime_date" DATE          NOT NULL,
    "overtime_hour" DECIMAL(5, 2) NOT NULL,
    "approved_date" DATE,
    "approved_by"   INT,

    -- Utility columns
    "status"        SMALLINT      NOT NULL DEFAULT 1,
    "flag"          INT           NOT NULL DEFAULT 0,
    "meta"          VARCHAR(255),
    "created_at"    TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by"    INT,
    "updated_at"    TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_by"    INT,
    "deleted_at"    TIMESTAMPTZ,
    "deleted_by"    INT,

    CONSTRAINT unique_user_overtime_date UNIQUE ("fk_user_id", "overtime_date")
);

DROP TABLE IF EXISTS "reimbursements";
CREATE TABLE IF NOT EXISTS "reimbursements"
(
    "id"            SERIAL PRIMARY KEY,
    "fk_user_id"    INT         NOT NULL,
    "description"   TEXT        NOT NULL,
    "amount"        DECIMAL     NOT NULL,
    "approved_date" DATE,
    "approved_by"   INT,

    -- Utility columns
    "status"        SMALLINT    NOT NULL DEFAULT 1,
    "flag"          INT         NOT NULL DEFAULT 0,
    "meta"          VARCHAR(255),
    "created_at"    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by"    INT,
    "updated_at"    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_by"    INT,
    "deleted_at"    TIMESTAMPTZ,
    "deleted_by"    INT
);
