-- migrate:up
-- write statements below this line

CREATE TABLE
    IF NOT EXISTS "investigation_results_metadata" (
        "id" BIGSERIAL NOT NULL PRIMARY KEY,
        "investigation_result_id" BIGINT NOT NULL,
        "qc_flag" VARCHAR (50),
        "qc_lot_number" VARCHAR (50) NOT NULL,
        "qc_value" VARCHAR (50) NOT NULL,
        "qc_west_gard_warning" VARCHAR (255) NOT NULL,
        "qc_status" varchar (50) NOT NULL,
        "created_by" BIGINT NOT NULL,
        "updated_by" BIGINT NOT NULL,
        "deleted_by" BIGINT DEFAULT NULL,
        "created_at" TIMESTAMPTZ NOT NULL,
        "updated_at" TIMESTAMPTZ NOT NULL,
        "deleted_at" TIMESTAMPTZ DEFAULT NULL
    );

CREATE INDEX IF NOT EXISTS "idx_investigation_result_id" ON "investigation_results_metadata" ("investigation_result_id");

ALTER TABLE "investigation_results_metadata"
    ADD CONSTRAINT "fk_investigation_results_metadata_investigation_result_id"
    FOREIGN KEY ("investigation_result_id")
    REFERENCES "investigation_results" ("id");

-- migrate:down
-- write rollback statements below this line

DROP TABLE IF EXISTS "investigation_results_metadata";
