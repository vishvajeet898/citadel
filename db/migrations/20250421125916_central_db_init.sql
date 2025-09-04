-- migrate:up
-- write statements below this line

CREATE TABLE IF NOT EXISTS "attachments" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "task_id" BIGINT NULL,
    "investigation_result_id" BIGINT NULL,
    "reference" VARCHAR(100) NOT NULL,
    "attachment_url" VARCHAR(255) NOT NULL,
    "thumbnail_url" VARCHAR(255) NULL,
    "thumbnail_reference" VARCHAR(100) NULL,
    "attachment_type" VARCHAR(50) NOT NULL,
    "attachment_label" VARCHAR(50) NULL,
    "is_reportable" BOOLEAN NOT NULL,
    "extension" VARCHAR(10) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "co_authorized_pathologists" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "task_id" BIGINT DEFAULT NULL,
    "investigation_result_id" BIGINT DEFAULT NULL,
    "co_authorized_by" BIGINT NOT NULL,
    "co_authorized_to" BIGINT NOT NULL,
    "co_authorized_at" TIMESTAMPTZ DEFAULT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "external_investigation_results" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "contact_id" BIGINT NOT NULL,
    "master_investigation_id" BIGINT DEFAULT NULL,
    "master_investigation_method_mapping_id" BIGINT DEFAULT NULL,
    "system_external_investigation_result_id" BIGINT NOT NULL,
    "system_external_report_id" BIGINT NOT NULL,
    "loinc_code" VARCHAR(100) NOT NULL,
    "investigation_name" VARCHAR(100) NOT NULL,
    "investigation_value" VARCHAR(100) NOT NULL,
    "uom" VARCHAR(100) DEFAULT NULL,
    "reference_range_text" VARCHAR(500) DEFAULT NULL,
    "is_abnormal" BOOLEAN NOT NULL,
    "lab_name" VARCHAR(100) NULL,
    "reported_at" TIMESTAMPTZ NOT NULL,
    "abnormality" VARCHAR(100) DEFAULT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "investigation_data" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "investigation_result_id" BIGINT NOT NULL,
    "data" text NOT NULL,
    "data_type" VARCHAR(50) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "investigation_results" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "test_details_id" BIGINT DEFAULT NULL,
    "oms_test_id" VARCHAR(30) DEFAULT NULL,
    "oms_order_id" VARCHAR(30) DEFAULT NULL,
    "master_investigation_id" BIGINT NOT NULL,
    "master_investigation_method_mapping_id" BIGINT NOT NULL,
    "investigation_name" VARCHAR(255) NOT NULL,
    "investigation_value" VARCHAR(255) NOT NULL,
    "result_representation_type" VARCHAR(100) DEFAULT NULL,
    "department" VARCHAR(100) DEFAULT NULL,
    "uom" VARCHAR(100) DEFAULT NULL,
    "method" VARCHAR(100) DEFAULT NULL,
    "method_type" VARCHAR(50) DEFAULT NULL,
    "reference_range_text" TEXT DEFAULT NULL,
    "lis_code" VARCHAR(50) NOT NULL,
    "abnormality" VARCHAR(50) DEFAULT NULL,
    "is_abnormal" BOOLEAN NOT NULL,
    "approved_by" BIGINT DEFAULT NULL,
    "approved_at" TIMESTAMPTZ DEFAULT NULL,
    "entered_by" BIGINT DEFAULT NULL,
    "entered_at" TIMESTAMPTZ DEFAULT NULL,
    "investigation_status" VARCHAR(100) DEFAULT NULL,
    "is_auto_approved" BOOLEAN NOT NULL,
    "is_non_reportable" BOOLEAN NOT NULL,
    "auto_verified" BOOLEAN NOT NULL,
    "is_nabl_approved" BOOLEAN NOT NULL,
    "source" VARCHAR(50) NOT NULL,
    "is_critical" BOOLEAN NOT NULL,
    "device_value" VARCHAR(255) DEFAULT NULL,
    "approval_source" VARCHAR(20) DEFAULT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "patient_details" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "dob" DATE DEFAULT NULL,
    "expected_dob" DATE DEFAULT NULL,
    "gender" VARCHAR(10) NOT NULL,
    "number" VARCHAR(20) NOT NULL,
    "system_patient_id" VARCHAR(50) DEFAULT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "remarks" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "investigation_result_id" BIGINT NOT NULL,
    "description" TEXT NOT NULL,
    "remark_type" VARCHAR(50) NOT NULL,
    "remark_by" BIGINT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "rerun_investigation_results" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "test_details_id" BIGINT NOT NULL,
    "master_investigation_id" BIGINT NOT NULL,
    "investigation_name" VARCHAR(255) NOT NULL,
    "investigation_value" VARCHAR(255) NOT NULL,
    "result_representation_type" VARCHAR(50) DEFAULT NULL,
    "lis_code" VARCHAR(50) NOT NULL,
    "rerun_triggered_by" BIGINT NOT NULL,
    "rerun_triggered_at" TIMESTAMPTZ NOT NULL,
    "rerun_reason" VARCHAR(255) NOT NULL,
    "rerun_remarks" TEXT DEFAULT NULL,
    "device_value" VARCHAR(255) DEFAULT NULL,
    "entered_by" BIGINT DEFAULT NULL,
    "entered_at" TIMESTAMPTZ DEFAULT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "tasks" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "order_id" BIGINT NOT NULL,
    "request_id" BIGINT NOT NULL,
    "oms_order_id" VARCHAR(30) DEFAULT NULL,
    "oms_request_id" VARCHAR(30) DEFAULT NULL,
    "lab_id" BIGINT NOT NULL,
    "city_code" VARCHAR(10) NOT NULL,
    "status" VARCHAR(50) NOT NULL,
    "previous_status" VARCHAR(50) DEFAULT NULL,
    "order_type" VARCHAR(50) NOT NULL,
    "patient_details_id" BIGINT NOT NULL,
    "doctor_tat" TIMESTAMPTZ DEFAULT NULL,
    "is_active" BOOLEAN NOT NULL,
    "completed_at" TIMESTAMPTZ DEFAULT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "task_metadata" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "task_id" BIGINT NOT NULL,
    "contains_package" BOOLEAN NOT NULL,
    "contains_morphle" BOOLEAN NOT NULL,
    "is_critical" BOOLEAN NOT NULL,
    "doctor_name" VARCHAR(255) DEFAULT NULL,
    "doctor_number" VARCHAR(20) DEFAULT NULL,
    "doctor_notes" TEXT DEFAULT NULL,
    "partner_name" VARCHAR(255) DEFAULT NULL,
    "last_event_sent_at" TIMESTAMPTZ DEFAULT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "task_pathologist_mapping" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "pathologist_id" BIGINT NOT NULL,
    "task_id" BIGINT DEFAULT NULL,
    "is_active" BOOLEAN NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "task_visit_mapping" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "task_id" BIGINT NOT NULL,
    "visit_id" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "templates" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "title" VARCHAR(255) DEFAULT NULL,
    "template_type" VARCHAR(50) NOT NULL,
    "description" TEXT DEFAULT NULL,
    "display_order" INTEGER DEFAULT 0,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "test_details" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "task_id" BIGINT NOT NULL,
    "oms_order_id" VARCHAR(30) DEFAULT NULL,
    "city_code" VARCHAR(10) DEFAULT NULL,
    "oms_test_id" BIGINT NOT NULL,
    "central_oms_test_id" VARCHAR(30) DEFAULT NULL,
    "test_name" VARCHAR(255) NOT NULL,
    "lis_code" VARCHAR(50) DEFAULT NULL,
    "master_test_id" BIGINT NOT NULL,
    "master_package_id" BIGINT DEFAULT NULL,
    "test_type" VARCHAR(10) DEFAULT NULL,
    "department" VARCHAR(100) DEFAULT NULL,
    "status" VARCHAR(50) NOT NULL,
    "doctor_tat" TIMESTAMPTZ DEFAULT NULL,
    "is_auto_approved" BOOLEAN NOT NULL,
    "report_sent_at" TIMESTAMPTZ DEFAULT NULL,
    "is_manual_report_upload" BOOLEAN DEFAULT FALSE,
    "lab_id" BIGINT DEFAULT NULL,
    "processing_lab_id" BIGINT DEFAULT NULL,
    "lab_eta" TIMESTAMPTZ DEFAULT NULL,
    "lab_tat" BIGINT DEFAULT NULL,
    "report_eta" TIMESTAMPTZ DEFAULT NULL,
    "oms_status" VARCHAR(25) DEFAULT NULL,
    "report_status" VARCHAR(25) DEFAULT NULL,
    "is_duplicate" BOOLEAN DEFAULT FALSE,
    "approval_source" VARCHAR(20) DEFAULT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "test_details_metadata" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "test_details_id" BIGINT NOT NULL,
    "barcodes" VARCHAR(255) NOT NULL,
    "is_critical" BOOLEAN DEFAULT FALSE,
    "is_completed_in_oms" BOOLEAN DEFAULT FALSE,
    "picked_at" TIMESTAMPTZ DEFAULT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "users" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "user_name" VARCHAR(255) NOT NULL,
    "system_user_id" VARCHAR(50) DEFAULT NULL,
    "user_type" VARCHAR(50) NOT NULL,
    "email" VARCHAR(255) DEFAULT NULL,
    "attune_user_id" VARCHAR(50) DEFAULT NULL,
    "agent_id" VARCHAR(50) DEFAULT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "order_details" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "oms_order_id" VARCHAR(30) NOT NULL,
    "oms_request_id" VARCHAR(30) NOT NULL,
    "uuid" UUID DEFAULT NULL,
    "city_code" VARCHAR(10) NOT NULL,
    "patient_details_id" BIGINT NOT NULL,
    "order_status" VARCHAR(50) NOT NULL,
    "partner_id" BIGINT DEFAULT NULL,
    "doctor_id" BIGINT DEFAULT NULL,
    "trf_id" VARCHAR(255) DEFAULT NULL,
    "request_source" VARCHAR(255) DEFAULT NULL,
    "bulk_order_id" BIGINT DEFAULT NULL,
    "servicing_lab_id" BIGINT NOT NULL,
    "collection_type" INT NOT NULL,
    "camp_id" BIGINT DEFAULT NULL,
    "referred_by" VARCHAR(255) DEFAULT NULL,
    "collected_on" TIMESTAMPTZ DEFAULT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "samples" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "oms_city_code" VARCHAR(10) NOT NULL,
    "oms_order_id" VARCHAR(30) NOT NULL,
    "oms_request_id" VARCHAR(30) NOT NULL,
    "visit_id" VARCHAR(100) DEFAULT NULL,
    "vial_type_id" BIGINT NOT NULL,
    "barcode" VARCHAR(100) DEFAULT NULL,
    "status" VARCHAR(50) NOT NULL,
    "lab_id" BIGINT DEFAULT NULL,
    "destination_lab_id" BIGINT DEFAULT NULL,
    "rejection_reason" VARCHAR(255) DEFAULT NULL,
    "parent_sample_id" BIGINT DEFAULT NULL,
    "sample_number" BIGINT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE  IF NOT EXISTS "samples_audit" (
    "log_action" VARCHAR(3) DEFAULT 'INS',
    "log_id" BIGSERIAL NOT NULL PRIMARY KEY,
    "log_timestamp" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "id" BIGINT NOT NULL,
    "oms_city_code" VARCHAR(10) NOT NULL,
    "oms_order_id" VARCHAR(30) NOT NULL,
    "oms_request_id" VARCHAR(30) NOT NULL,
    "visit_id" VARCHAR(100) DEFAULT NULL,
    "vial_type_id" BIGINT NOT NULL,
    "barcode" VARCHAR(100) DEFAULT NULL,
    "status" VARCHAR(50) NOT NULL,
    "lab_id" BIGINT DEFAULT NULL,
    "destination_lab_id" BIGINT DEFAULT NULL,
    "rejection_reason" VARCHAR(255) DEFAULT NULL,
    "parent_sample_id" BIGINT DEFAULT NULL,
    "sample_number" BIGINT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE OR REPLACE FUNCTION samples_audit_insert_function()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO samples_audit(
        log_action, log_timestamp,
        id, oms_city_code, oms_order_id, oms_request_id, visit_id, vial_type_id, barcode, status,
        lab_id, destination_lab_id, rejection_reason, parent_sample_id, sample_number,
        created_at, updated_at, deleted_at,
        created_by, updated_by, deleted_by
    )
    VALUES (
        'INS', NEW.updated_at,
        NEW.id, NEW.oms_city_code, NEW.oms_order_id, NEW.oms_request_id, NEW.visit_id,
        NEW.vial_type_id, NEW.barcode, NEW.status,
        NEW.lab_id, NEW.destination_lab_id, NEW.rejection_reason, NEW.parent_sample_id,
        NEW.sample_number, NEW.created_at, NEW.updated_at, NEW.deleted_at,
        NEW.created_by, NEW.updated_by, NEW.deleted_by
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER "samples_audit_insert_trigger"
AFTER INSERT ON "samples"
FOR EACH ROW
EXECUTE FUNCTION samples_audit_insert_function();


CREATE OR REPLACE FUNCTION samples_audit_update_function()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO samples_audit(
        log_action, log_timestamp,
        id, oms_city_code, oms_order_id, oms_request_id, visit_id, vial_type_id, barcode, status,
        lab_id, destination_lab_id, rejection_reason, parent_sample_id, sample_number,
        created_at, updated_at, deleted_at,
        created_by, updated_by, deleted_by
    )
    VALUES (
        'UPD', NEW.updated_at,
        NEW.id, NEW.oms_city_code, NEW.oms_order_id, NEW.oms_request_id, NEW.visit_id,
        NEW.vial_type_id, NEW.barcode, NEW.status,
        NEW.lab_id, NEW.destination_lab_id, NEW.rejection_reason, NEW.parent_sample_id,
        NEW.sample_number, NEW.created_at, NEW.updated_at, NEW.deleted_at,
        NEW.created_by, NEW.updated_by, NEW.deleted_by
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER "samples_audit_update_trigger"
AFTER UPDATE ON "samples"
FOR EACH ROW
EXECUTE FUNCTION samples_audit_update_function();


CREATE OR REPLACE FUNCTION samples_audit_delete_function()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO samples_audit(
        log_action, log_timestamp,
        id, oms_city_code, oms_order_id, oms_request_id, visit_id, vial_type_id, barcode, status,
        lab_id, destination_lab_id, rejection_reason, parent_sample_id, sample_number,
        created_at, updated_at, deleted_at,
        created_by, updated_by, deleted_by
    )
    VALUES (
        'DEL', OLD.deleted_at,
        OLD.id, OLD.oms_city_code, OLD.oms_order_id, OLD.oms_request_id, OLD.visit_id,
        OLD.vial_type_id, OLD.barcode, OLD.status,
        OLD.lab_id, OLD.destination_lab_id, OLD.rejection_reason, OLD.parent_sample_id,
        OLD.sample_number, OLD.created_at, OLD.updated_at, OLD.deleted_at,
        OLD.created_by, OLD.updated_by, OLD.deleted_by
    );
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER "samples_audit_delete_trigger"
AFTER DELETE ON "samples"
FOR EACH ROW
EXECUTE FUNCTION samples_audit_delete_function();


CREATE TABLE IF NOT EXISTS "sample_metadata" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "oms_city_code" VARCHAR(10) NOT NULL,
    "oms_order_id" VARCHAR(30) NOT NULL,
    "sample_id" BIGINT NOT NULL,
    "collection_sequence_number" BIGINT NOT NULL,
    "last_updated_at" TIMESTAMPTZ DEFAULT NULL,
    "transferred_at" TIMESTAMPTZ DEFAULT NULL,
    "outsourced_at" TIMESTAMPTZ DEFAULT NULL,
    "collected_at" TIMESTAMPTZ DEFAULT NULL,
    "received_at" TIMESTAMPTZ DEFAULT NULL,
    "accessioned_at" TIMESTAMPTZ DEFAULT NULL,
    "rejected_at" TIMESTAMPTZ DEFAULT NULL,
    "not_received_at" TIMESTAMPTZ DEFAULT NULL,
    "lis_sync_at" TIMESTAMPTZ DEFAULT NULL,
    "barcode_scanned_at" TIMESTAMPTZ DEFAULT NULL,
    "not_received_reason" VARCHAR(255) DEFAULT NULL,
    "barcode_image_url" VARCHAR(255) DEFAULT NULL,
    "task_sequence" BIGINT DEFAULT NULL,
    "collect_later_reason" VARCHAR(255) DEFAULT NULL,
    "rejecting_lab" BIGINT DEFAULT NULL,
    "collected_volume" BIGINT DEFAULT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "sample_metadata_audit" (
    "log_action" VARCHAR(3) DEFAULT 'INS',
    "log_id" BIGSERIAL NOT NULL PRIMARY KEY,
    "log_timestamp" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "id" BIGINT NOT NULL,
    "oms_city_code" VARCHAR(10) NOT NULL,
    "oms_order_id" VARCHAR(30) NOT NULL,
    "sample_id" BIGINT NOT NULL,
    "collection_sequence_number" BIGINT NOT NULL,
    "last_updated_at" TIMESTAMPTZ DEFAULT NULL,
    "transferred_at" TIMESTAMPTZ DEFAULT NULL,
    "outsourced_at" TIMESTAMPTZ DEFAULT NULL,
    "collected_at" TIMESTAMPTZ DEFAULT NULL,
    "received_at" TIMESTAMPTZ DEFAULT NULL,
    "accessioned_at" TIMESTAMPTZ DEFAULT NULL,
    "rejected_at" TIMESTAMPTZ DEFAULT NULL,
    "not_received_at" TIMESTAMPTZ DEFAULT NULL,
    "lis_sync_at" TIMESTAMPTZ DEFAULT NULL,
    "barcode_scanned_at" TIMESTAMPTZ DEFAULT NULL,
    "not_received_reason" VARCHAR(255) DEFAULT NULL,
    "barcode_image_url" VARCHAR(255) DEFAULT NULL,
    "task_sequence" BIGINT DEFAULT NULL,
    "collect_later_reason" VARCHAR(255) DEFAULT NULL,
    "rejecting_lab" BIGINT DEFAULT NULL,
    "collected_volume" BIGINT DEFAULT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE OR REPLACE FUNCTION sample_metadata_audit_insert_function()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO "sample_metadata_audit"(
        log_action, log_timestamp, id, oms_city_code, oms_order_id, sample_id,
        collection_sequence_number, last_updated_at, transferred_at,
        outsourced_at, collected_at, received_at, accessioned_at, rejected_at, not_received_at,
        lis_sync_at, barcode_scanned_at, not_received_reason, barcode_image_url,
        task_sequence, collect_later_reason, rejecting_lab, collected_volume,
        created_at, updated_at, deleted_at, created_by, updated_by, deleted_by
    )
    VALUES (
        'INS', NEW.updated_at,
        NEW.id, NEW.oms_city_code, NEW.oms_order_id, NEW.sample_id,
        NEW.collection_sequence_number, NEW.last_updated_at, NEW.transferred_at,
        NEW.outsourced_at, NEW.collected_at, NEW.received_at, NEW.accessioned_at, NEW.rejected_at, NEW.not_received_at,
        NEW.lis_sync_at, NEW.barcode_scanned_at, NEW.not_received_reason, NEW.barcode_image_url,
        NEW.task_sequence, NEW.collect_later_reason, NEW.rejecting_lab, NEW.collected_volume,
        NEW.created_at, NEW.updated_at, NEW.deleted_at, NEW.created_by, NEW.updated_by, NEW.deleted_by
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER "sample_metadata_audit_insert_trigger"
AFTER INSERT ON "sample_metadata"
FOR EACH ROW
EXECUTE FUNCTION sample_metadata_audit_insert_function();


CREATE OR REPLACE FUNCTION sample_metadata_audit_update_function()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO "sample_metadata_audit"(
        log_action, log_timestamp, id, oms_city_code, oms_order_id, sample_id,
        collection_sequence_number, last_updated_at, transferred_at,
        outsourced_at, collected_at, received_at, accessioned_at, rejected_at, not_received_at,
        lis_sync_at, barcode_scanned_at, not_received_reason, barcode_image_url,
        task_sequence, collect_later_reason, rejecting_lab, collected_volume,
        created_at, updated_at, deleted_at, created_by, updated_by, deleted_by
    )
    VALUES (
        'UPD', NEW.updated_at,
        NEW.id, NEW.oms_city_code, NEW.oms_order_id, NEW.sample_id,
        NEW.collection_sequence_number, NEW.last_updated_at, NEW.transferred_at,
        NEW.outsourced_at, NEW.collected_at, NEW.received_at, NEW.accessioned_at, NEW.rejected_at, NEW.not_received_at,
        NEW.lis_sync_at, NEW.barcode_scanned_at, NEW.not_received_reason, NEW.barcode_image_url,
        NEW.task_sequence, NEW.collect_later_reason, NEW.rejecting_lab, NEW.collected_volume,
        NEW.created_at, NEW.updated_at, NEW.deleted_at, NEW.created_by, NEW.updated_by, NEW.deleted_by
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER "sample_metadata_audit_update_trigger"
AFTER UPDATE ON "sample_metadata"
FOR EACH ROW
EXECUTE FUNCTION sample_metadata_audit_update_function();


CREATE OR REPLACE FUNCTION sample_metadata_audit_delete_function()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO "sample_metadata_audit"(
        log_action, log_timestamp, id, oms_city_code, oms_order_id, sample_id,
        collection_sequence_number, last_updated_at, transferred_at,
        outsourced_at, collected_at, received_at, accessioned_at, rejected_at, not_received_at,
        lis_sync_at, barcode_scanned_at, not_received_reason, barcode_image_url,
        task_sequence, collect_later_reason, rejecting_lab, collected_volume,
        created_at, updated_at, deleted_at, created_by, updated_by, deleted_by
    )
    VALUES (
        'DEL', NEW.updated_at,
        OLD.id, OLD.oms_city_code, OLD.oms_order_id, OLD.sample_id,
        OLD.collection_sequence_number, OLD.last_updated_at, OLD.transferred_at,
        OLD.outsourced_at, OLD.collected_at, OLD.received_at, OLD.accessioned_at, OLD.rejected_at, OLD.not_received_at,
        OLD.lis_sync_at, OLD.barcode_scanned_at, OLD.not_received_reason, OLD.barcode_image_url,
        OLD.task_sequence, OLD.collect_later_reason, OLD.rejecting_lab, OLD.collected_volume,
        OLD.created_at, OLD.updated_at, OLD.deleted_at, OLD.created_by, OLD.updated_by, OLD.deleted_by
    );
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER "sample_metadata_audit_delete_trigger"
AFTER DELETE ON "sample_metadata"
FOR EACH ROW
EXECUTE FUNCTION sample_metadata_audit_delete_function();


CREATE TABLE IF NOT EXISTS "test_sample_mapping" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "oms_city_code" VARCHAR(10) NOT NULL,
    "oms_test_id" VARCHAR(30) NOT NULL,
    "vial_type_id" BIGINT NOT NULL,
    "sample_id" BIGINT NOT NULL,
    "sample_number" BIGINT NOT NULL,
    "oms_order_id" VARCHAR(30) NOT NULL,
    "recollection_pending" BOOLEAN DEFAULT FALSE,
    "is_rejected" BOOLEAN DEFAULT FALSE,
    "rejection_reason" VARCHAR(255) DEFAULT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "created_by" BIGINT NOT NULL,
    "updated_by" BIGINT NOT NULL,
    "deleted_by" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "ets_events" (
  "id" BIGSERIAL NOT NULL PRIMARY KEY,
  "test_id" VARCHAR(30) DEFAULT NULL,
  "is_active" BOOLEAN DEFAULT TRUE,
  "created_at" TIMESTAMPTZ NOT NULL,
  "updated_at" TIMESTAMPTZ NOT NULL,
  "deleted_at" TIMESTAMPTZ DEFAULT NULL,
  "created_by" BIGINT NOT NULL,
  "updated_by" BIGINT NOT NULL,
  "deleted_by" BIGINT DEFAULT NULL
);

INSERT INTO "users" ("user_name","system_user_id","user_type","email","created_at","updated_at","deleted_at","created_by","updated_by","deleted_by","attune_user_id","agent_id") VALUES
('Priyanka Tiwari','961','pathologist','priyanka.tiwari@orangehealth.in','2024-08-09 10:29:14.91983+05:30','2024-10-17 17:11:03.660547+05:30',NULL,0,25,NULL,'1242838',NULL),
('Shubham Naik','961','pathologist','shubham.naik@orangehealth.in','2024-08-08 16:00:00+05:30','2024-10-17 17:11:03.691332+05:30',NULL,0,25,NULL,'1242838',NULL),
('Dr. Joan Maria Vynetta D','961','pathologist','vynetta@orangehealth.in','2024-08-08 18:02:00+05:30','2024-08-08 18:02:00+05:30',NULL,0,0,NULL,'1242838','vynetta'),
('Amit Kumar','961','pathologist','amit.k@orangehealth.in','2024-08-08 18:10:00+05:30','2024-10-17 17:11:03.718926+05:30',NULL,25,25,NULL,NULL,NULL),
('Dr. Adil Khanday','961','pathologist','adil@orangehealth.in','2024-08-20 22:29:23.8207+05:30','2025-04-19 15:24:10.018772+05:30','2025-04-19 15:24:10.018772+05:30',25,25,25,NULL,NULL),
('Dr. Madhutej TH','300607','pathologist','','2024-08-22 01:10:15.64198+05:30','2024-08-22 01:10:15.64198+05:30','2024-09-17 16:53:20.23676+05:30',0,0,25,'2204429',NULL),
('Dr. Aekta','349887','pathologist','aekta@orangehealth.in','2024-08-22 01:10:15.64198+05:30','2024-08-22 01:10:15.64198+05:30',NULL,0,0,NULL,'18160979','aekta'),
('Dr. Sanchit Singhal','9641','pathologist','sanchit@orangehealth.in','2024-08-22 01:10:15.64198+05:30','2024-08-22 01:10:15.64198+05:30',NULL,0,0,NULL,'1248689','sanchit'),
('Dr. Alekya Seerapu','39927','pathologist','alekya@orangehealth.in','2024-08-22 01:10:15.64198+05:30','2024-08-22 01:10:15.64198+05:30',NULL,0,0,NULL,'2214500','Alekya'),
('Dr. Astha Srivastava','591936','pathologist','','2024-08-22 01:10:15.64198+05:30','2024-08-22 01:10:15.64198+05:30','2024-09-17 16:53:36.525569+05:30',0,0,25,'18162506',NULL),
('Dr. Jushmita Pathak','488919','pathologist','jushmita.pathak@orangehealth.in','2024-08-22 01:10:15.64198+05:30','2024-08-22 01:10:15.64198+05:30',NULL,0,0,NULL,'18161236','Jushmita'),
('Dr. Noopur Srivastava','759342','pathologist','','2024-08-22 01:10:15.64198+05:30','2024-08-22 01:10:15.64198+05:30',NULL,0,0,NULL,'18161257',NULL),
('Dr. Abhijay Dharmadhikari','705255','pathologist','abhijay.d@orangehealth.in','2024-08-22 01:10:15.64198+05:30','2024-08-22 01:10:15.64198+05:30',NULL,0,0,NULL,'18165165','Abhijay'),
('Dr. Neha Jagdishrajji Gajbi','197883','pathologist','','2024-08-22 01:10:15.64198+05:30','2024-08-22 01:10:15.64198+05:30','2024-10-29 13:34:04.750524+05:30',25,25,25,'2186837',NULL),
('Dr. Radhika Puri','759343','pathologist','','2024-08-22 01:10:15.64198+05:30','2024-08-22 01:10:15.64198+05:30',NULL,0,0,NULL,'18162769',NULL),
('Dr. Chitra Chauhan','759341','pathologist','','2024-08-22 01:10:15.64198+05:30','2025-03-22 15:14:03.730027+05:30','2025-03-22 15:14:03.730027+05:30',25,25,25,'18162780',NULL),
('Dr. Ankita Nayak','238349','pathologist','ankita.nayak@orangehealth.in','2024-08-22 01:10:15.64198+05:30','2024-10-09 16:14:19.677955+05:30',NULL,25,25,NULL,'18160621',NULL),
('Dr. Ankita Singh','759396','pathologist','','2024-08-22 01:10:15.64198+05:30','2024-09-24 18:45:47.781525+05:30',NULL,0,25,NULL,'18165952',NULL),
('Dr. Akshay Prashantkumar Vadavadgi','786744','pathologist','akshay.v@orangehealth.in','2024-08-22 01:10:15.64198+05:30','2024-09-09 12:39:43.32779+05:30',NULL,0,0,NULL,'18166742','Dr.Akshay'),
('Dr. Ritika Doshi','867554','pathologist','ritika@orangehealth.in','2024-08-22 01:10:15.64198+05:30','2024-08-22 01:10:15.64198+05:30',NULL,0,0,NULL,'18168912','Ritika'),
('Dr. Srishti Miraj','1021582','pathologist','srishti.miraj@orangehealth.in','2024-08-22 01:10:15.64198+05:30','2024-08-22 01:10:15.64198+05:30',NULL,0,0,NULL,'18169652','SrishtiMiraj'),
('Dr. Anushree R','681881','pathologist','anushree@orangehealth.in','2024-08-22 01:10:15.64198+05:30','2024-08-22 01:10:15.64198+05:30',NULL,0,0,NULL,'18169975','Anushree'),
('Citadel',NULL,'system',NULL,'2024-08-24 21:24:09.499193+05:30','2024-08-24 21:24:09.499193+05:30',NULL,0,0,NULL,NULL,NULL),
('LIS',NULL,'system',NULL,'2024-08-24 21:24:49.531871+05:30','2024-08-24 21:24:49.531871+05:30',NULL,0,0,NULL,NULL,NULL),
('Dr. Shrish Sharma','961','pathologist','shrish.chandra@orangehealth.in','2024-08-27 12:48:57.762586+05:30','2024-10-17 17:11:03.832606+05:30',NULL,0,25,NULL,'18165952','Adil'),
('Dr. Sushma S','1130538','pathologist','sushma.s@orangehealth.in','2024-09-12 15:17:49.460837+05:30','2024-09-12 15:17:49.460837+05:30',NULL,0,0,NULL,'18170217','SushmaS'),
('Dr. Drisya Jagadees','1186492','pathologist','drisya.jagadees@orangehealth.in','2024-10-14 15:36:16.601171+05:30','2024-10-15 12:40:00.250907+05:30',NULL,25,25,NULL,'18170949',NULL),
('Vandan Rogheliya',NULL,'pathologist','vandan@orangehealth.in','2024-12-23 22:44:24.739358+05:30','2024-12-23 22:44:24.739358+05:30',NULL,25,25,NULL,NULL,NULL),
('Praveen Kumar',NULL,'pathologist','praveenkumar.k@orangehealth.in','2024-12-23 22:47:39.696615+05:30','2024-12-23 22:47:39.696615+05:30',NULL,25,25,NULL,NULL,NULL),
('Dr. Vinay Yadav','1327113','pathologist','vinay.yadav@orangehealth.in','2025-02-11 15:40:29.433309+05:30','2025-02-11 20:35:44.582215+05:30',NULL,25,25,NULL,'18172689',NULL),
('Dr. Sayantani Sarkar','1365385','pathologist','sayantani.sarkar@orangehealth.in','2025-02-25 15:12:41.441392+05:30','2025-02-25 15:12:41.441392+05:30',NULL,25,25,NULL,'18172767',NULL),
('Dr. Noora M Ahmed','1390022','pathologist','noora.ahmed@orangehealth.in','2025-03-10 11:33:50.2436+05:30','2025-03-10 11:33:50.2436+05:30',NULL,25,25,NULL,'2528750',NULL),
('Dr. Shivani Tiwari','1414615','pathologist','shivani.tiwari@orangehealth.in','2025-03-22 14:31:25.336762+05:30','2025-03-22 14:31:25.336762+05:30',NULL,25,25,NULL,'2528951',NULL);

INSERT INTO "users" ("email", "user_name", "user_type", "created_at", "updated_at", "created_by", "updated_by") VALUES
('chilumula.karthik@orangehealth.in','Chilumula','lab-technician', NOW(), NOW(), 25, 25),
('rajendra@orangehealth.in','Rajendra','super-admin', NOW(), NOW(), 25, 25),
('aman.singh@orangehealth.in','Aman Kumar ','crm', NOW(), NOW(), 25, 25),
('bipasa.sarkar@orangehealth.in','Bipasa','crm', NOW(), NOW(), 25, 25),
('ananth.kp@orangehealth.in','Ananth Kumar','lab-technician', NOW(), NOW(), 25, 25),
('ramya.m@orangehealth.in','Ramya ','crm', NOW(), NOW(), 25, 25),
('shaik.shahbaz@orangehealth.in','Shaik ','crm', NOW(), NOW(), 25, 25),
('famsheera.ma@orangehealth.in','Famsheera ','lab-technician', NOW(), NOW(), 25, 25),
('jagadish.n@orangehealth.in','Jagadish','lab-technician', NOW(), NOW(), 25, 25),
('jitendra.jha@orangehealth.in','Jitendra','lab-technician', NOW(), NOW(), 25, 25),
('daniel.dispatch@orangehealth.in','Daniel','crm', NOW(), NOW(), 25, 25),
('swetha.a@orangehealth.in','Swetha','lab-technician', NOW(), NOW(), 25, 25),
('shabana.salim@orangehealth.in','Shabana','lab-technician', NOW(), NOW(), 25, 25),
('sekar.subramani@orangehealth.in','Sekar ','lab-technician', NOW(), NOW(), 25, 25),
('amreen.shaikh@orangehealth.in','Amreen ','lab-technician', NOW(), NOW(), 25, 25),
('ranjith.rk@orangehealth.in','R Ranjith ','crm', NOW(), NOW(), 25, 25),
('siddharaj@orangehealth.in','Siddharaj','super-admin', NOW(), NOW(), 25, 25),
('gangaraju.kr@orangehealth.in','Gangaraju ','lab-technician', NOW(), NOW(), 25, 25),
('akhil.ahmed@orangehealth.in','Akhil Ahmed ','super-admin', NOW(), NOW(), 25, 25),
('sarjil.a@orangehealth.in','Sarjil ','lab-technician', NOW(), NOW(), 25, 25),
('vivek.tiwari@orangehealth.in','Vivek','admin', NOW(), NOW(), 25, 25),
('prenay.p@orangehealth.in','Prenay','admin', NOW(), NOW(), 25, 25),
('saqib.ahmad@orangehealth.in','Saqib Ahmad ','crm', NOW(), NOW(), 25, 25),
('tejaswini.chachiya@orangehealth.in','Tejaswini','lab-technician', NOW(), NOW(), 25, 25),
('ateeq.mohammed@orangehealth.in','Ateeq','lab-technician', NOW(), NOW(), 25, 25),
('kanak@orangehealth.in','Kanak','super-admin', NOW(), NOW(), 25, 25),
('gopal.k@orangehealth.in','gopal','lab-technician', NOW(), NOW(), 25, 25),
('subrata.biswas@orangehealth.in','Subrata Kumar ','super-admin', NOW(), NOW(), 25, 25),
('aishwarya.m@orangehealth.in','Aishwarya ','crm', NOW(), NOW(), 25, 25),
('neeraj.tyagi@orangehealth.in','Neeraj','lab-technician', NOW(), NOW(), 25, 25),
('kavyakumari.nb@orangehealth.in','Kavyakumari ','lab-technician', NOW(), NOW(), 25, 25),
('gireesh.bs@orangehealth.in','Gireesh ','lab-technician', NOW(), NOW(), 25, 25),
('ameersohil.nayak@orangehealth.in','Ameersohil B ','lab-technician', NOW(), NOW(), 25, 25),
('darshana.nishar@orangehealth.in','Darshana Ashok ','admin', NOW(), NOW(), 25, 25),
('chandrashekar.hugar@orangehealth.in','Chandrashekar','lab-technician', NOW(), NOW(), 25, 25),
('parthugari.rakesh@orangehealth.in','Parthugari ','lab-technician', NOW(), NOW(), 25, 25),
('chinmayi.m@orangehealth.in','Chinmayi ','lab-technician', NOW(), NOW(), 25, 25),
('chaithra.ys@orangehealth.in','Chaithra','lab-technician', NOW(), NOW(), 25, 25),
('anand.k@orangehealth.in','Anand ','crm', NOW(), NOW(), 25, 25),
('shaikh.amjad@orangehealth.in','Shaikh','lab-technician', NOW(), NOW(), 25, 25),
('dushyant.raghav@orangehealth.in','Dushyant','admin', NOW(), NOW(), 25, 25),
('yogesh.ak@orangehealth.in','Yogesh Amar Kamable','lab-technician', NOW(), NOW(), 25, 25),
('araveetikumar.reddy@orangehealth.in','Araveeti ','lab-technician', NOW(), NOW(), 25, 25),
('maskani.v@orangehealth.in','Maskani ','lab-technician', NOW(), NOW(), 25, 25),
('mohd.tausif@orangehealth.in','Mohd Tausif ','lab-technician', NOW(), NOW(), 25, 25),
('justine.nirosha@orangehealth.in','Justine','crm', NOW(), NOW(), 25, 25),
('pavan.kishore@orangehealth.in','Pavan','lab-technician', NOW(), NOW(), 25, 25),
('abbasali.r@orangehealth.in','Abbas Ali ','crm', NOW(), NOW(), 25, 25),
('amol.jadhav@orangehealth.in','Amol','lab-technician', NOW(), NOW(), 25, 25),
('ankit.kumar@orangehealth.in','Ankit ','lab-technician', NOW(), NOW(), 25, 25),
('darveena@orangehealth.in','Darveena','pathologist', NOW(), NOW(), 25, 25),
('savitri@orangehealth.in','Savitri ','lab-technician', NOW(), NOW(), 25, 25),
('jay.pasupuleti@orangehealth.in','Pasupuleti ','lab-technician', NOW(), NOW(), 25, 25),
('shaik.khalid@orangehealth.in','Shaik Mahammad ','lab-technician', NOW(), NOW(), 25, 25),
('mehak.khannum@orangehealth.in','Mehak ','crm', NOW(), NOW(), 25, 25),
('ranjini.m@orangehealth.in','Ranjini','crm', NOW(), NOW(), 25, 25),
('sanjay.k@orangehealth.in','Sanjay','lab-technician', NOW(), NOW(), 25, 25),
('ajaykumar@orangehealth.in','Ajay','lab-technician', NOW(), NOW(), 25, 25),
('kalvala.sushma@orangehealth.in','Kalvala ','lab-technician', NOW(), NOW(), 25, 25),
('erum.khan@orangehealth.in','Erum','lab-technician', NOW(), NOW(), 25, 25),
('madhushree.b@orangehealth.in','Madhushree','super-admin', NOW(), NOW(), 25, 25),
('rajini@orangehealth.in','Rajini ','admin', NOW(), NOW(), 25, 25),
('harsh.sharma@orangehealth.in','Harsh','lab-technician', NOW(), NOW(), 25, 25),
('anand.viswan@orangehealth.in','Anand ','crm', NOW(), NOW(), 25, 25),
('souravkar@orangehealth.in','Sourav ','admin', NOW(), NOW(), 25, 25),
('sudhanshu@orangehealth.in','Sudhanshu ','lab-technician', NOW(), NOW(), 25, 25),
('madhusudhan.a@orangehealth.in','Madhusudhan ','lab-technician', NOW(), NOW(), 25, 25),
('nagati.venkat@orangehealth.in','Nagati Venkat ','lab-technician', NOW(), NOW(), 25, 25),
('abdulmoeed.chaudhary@orangehealth.in','Abdul Moeed ','lab-technician', NOW(), NOW(), 25, 25),
('mohd.samar@orangehealth.in','Mohd ','lab-technician', NOW(), NOW(), 25, 25),
('nuzhath.kouser@orangehealth.in','Nuzhath ','lab-technician', NOW(), NOW(), 25, 25),
('khushboo.bharti@orangehealth.in','Khushboo ','crm', NOW(), NOW(), 25, 25),
('abhishek.k@orangehealth.in','Abhishek ','lab-technician', NOW(), NOW(), 25, 25),
('george.c@orangehealth.in','George ','lab-technician', NOW(), NOW(), 25, 25),
('abdul.shaikh@orangehealth.in','Abdul Amulal ','lab-technician', NOW(), NOW(), 25, 25),
('lokesh.ch@orangehealth.in','Lokesh Roshan ','crm', NOW(), NOW(), 25, 25),
('ruman@orangehealth.in','Ruman ','lab-technician', NOW(), NOW(), 25, 25),
('ngen-lab@orangehealth.in','Noida General ','lab-technician', NOW(), NOW(), 25, 25),
('sandeep.kagiyavar@orangehealth.in','Sandeep ','lab-technician', NOW(), NOW(), 25, 25),
('harshavardhana.reddy@orangehealth.in','Harshavardhana','crm', NOW(), NOW(), 25, 25),
('manjunath.b@orangehealth.in','Manjunath Basavaraj ','lab-technician', NOW(), NOW(), 25, 25),
('mohammed.kalaam@orangehealth.in','Mohammed','lab-technician', NOW(), NOW(), 25, 25),
('manisha.ds@orangehealth.in','Manisha Dattatray ','lab-technician', NOW(), NOW(), 25, 25),
('vinay.kumar@orangehealth.in','Vinay','lab-technician', NOW(), NOW(), 25, 25),
('santhosha.h@orangehealth.in','Santhosha','lab-technician', NOW(), NOW(), 25, 25),
('amarjeet@orangehealth.in','Amarjeet','lab-technician', NOW(), NOW(), 25, 25),
('somashekar@orangehealth.in','Somashekar','super-admin', NOW(), NOW(), 25, 25),
('zaffar.khan@orangehealth.in','Zaffar','lab-technician', NOW(), NOW(), 25, 25),
('shalini.d@orangehealth.in','Mary Shalini ','crm', NOW(), NOW(), 25, 25),
('sowmiya.anandan@orangehealth.in','Sowmiya','lab-technician', NOW(), NOW(), 25, 25),
('mercy.merlin@orangehealth.in','Mercy Merlin S ','lab-technician', NOW(), NOW(), 25, 25),
('dikshitha.yb@orangehealth.in','YB ','lab-technician', NOW(), NOW(), 25, 25),
('akhil.anand@orangehealth.in','Akhil ','crm', NOW(), NOW(), 25, 25),
('gopal.n@orangehealth.in','Gopal ','lab-technician', NOW(), NOW(), 25, 25),
('krutika.satam@orangehealth.in','Krutika ','lab-technician', NOW(), NOW(), 25, 25),
('vijay.kumar@orangehealth.in','Vijay ','lab-technician', NOW(), NOW(), 25, 25),
('harshavardhan@orangehealth.in','Harshavardhan ','crm', NOW(), NOW(), 25, 25),
('stella.l@orangehealth.in','Stella ','lab-technician', NOW(), NOW(), 25, 25),
('sundeep.ponkam@orangehealth.in','Sundeep','super-admin', NOW(), NOW(), 25, 25),
('mohan@orangehealth.in','Mohan Reddy','lab-technician', NOW(), NOW(), 25, 25),
('rajashekhar.a@orangehealth.in','Rajashekhar ','crm', NOW(), NOW(), 25, 25),
('rajeshwari@orangehealth.in','Rajeshwari','lab-technician', NOW(), NOW(), 25, 25),
('ramkumar@orangehealth.in','Ram','lab-technician', NOW(), NOW(), 25, 25),
('veerababu.m@orangehealth.in','Mareedu ','lab-technician', NOW(), NOW(), 25, 25),
('balaji.mg@orangehealth.in','Balaji Naika ','lab-technician', NOW(), NOW(), 25, 25),
('krishan.kumar@orangehealth.in','Krishan ','lab-technician', NOW(), NOW(), 25, 25),
('ajith.kumar@orangehealth.in','Ajith ','crm', NOW(), NOW(), 25, 25),
('sanchit.gupta@orangehealth.in','Sanchit ','crm', NOW(), NOW(), 25, 25),
('ajoy.sen@orangehealth.in','Ajoy ','lab-technician', NOW(), NOW(), 25, 25),
('hemanth.r@orangehealth.in','Hemanth','super-admin', NOW(), NOW(), 25, 25);

INSERT INTO "templates" ("title","template_type","description","created_at","updated_at","deleted_at","created_by","updated_by","deleted_by","display_order") VALUES
('','rerun_reason','Incomplete data provided','2024-08-03 16:24:43.762861+05:30','2024-09-25 10:48:46.095988+05:30',NULL,25,25,NULL,1),
('','rerun_reason','Rerun in dilution','2024-08-03 16:26:51.507399+05:30','2024-09-25 10:48:46.121698+05:30',NULL,25,25,NULL,2),
('','rerun_reason','Unsatisfactory slide','2024-08-03 16:26:51.507399+05:30','2024-09-25 10:48:46.145704+05:30',NULL,25,25,NULL,3),
('','rerun_reason','Sample rejection required','2024-08-03 16:26:51.507399+05:30','2024-09-25 10:48:46.169438+05:30',NULL,25,25,NULL,4),
('','rerun_reason','Confirmation needed for demographics','2024-08-03 16:26:51.507399+05:30','2024-09-25 10:48:46.22469+05:30',NULL,25,25,NULL,5),
('','withheld_reason','Patient communication required','2024-08-03 16:27:30.344021+05:30','2024-09-25 10:49:45.861163+05:30',NULL,25,25,NULL,1),
('','medical_remark','Kindly correlate clinically and with iron studies for further evaluation.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:01.491313+05:30',NULL,25,25,NULL,3),
('','medical_remark','Kindly correlate clinically and with IgE levels for further evaluation.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:01.516612+05:30',NULL,25,25,NULL,4),
('','medical_remark','Kindly correlate clinically and refer to the interpretation provided below. Recommended fasting status sample for confirmation if clinically indicated.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:01.882438+05:30',NULL,25,25,NULL,12),
('','medical_remark','Low Iron not correlating with Hb - Kindly correlate clinically and rule out inflammation / subclinical iron deficiency.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:01.921176+05:30',NULL,25,25,NULL,13),
('','medical_remark','Advice to rule out viral etiology in view of leukopenia.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:01.611897+05:30',NULL,25,25,NULL,5),
('','medical_remark','Kindly correlate clinically and advise further evaluation.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:28:57.17437+05:30','2024-09-25 11:28:57.17437+05:30',25,25,25,0),
('','medical_remark','Manually checked. No clumps noted. Kindly correlate clinically and advise follow up.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:01.702193+05:30',NULL,25,25,NULL,8),
('','medical_remark','VLDL values cannot be estimated due to high triglycerides.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:03.197154+05:30',NULL,25,25,NULL,48),
('','medical_remark','After 1 hour of 75 grams of oral glucose.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:01.983129+05:30',NULL,25,25,NULL,14),
('','medical_remark','After 2 hours of 75 grams of oral glucose.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.095179+05:30',NULL,25,25,NULL,15),
('','medical_remark','After 1 hour of 50 grams of oral glucose.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.13255+05:30',NULL,25,25,NULL,16),
('','medical_remark','Kindly correlate clinically and suggest follow up.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:01.43887+05:30',NULL,25,25,NULL,1),
('','medical_remark','Large platelets noted. Kindly correlate clinically and advise follow up.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:01.639068+05:30',NULL,25,25,NULL,6),
('','medical_remark','Manually checked. No clumps noted. Large platelets and few giant platelets noted. Kindly correlate clinically and advise follow up.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:01.740114+05:30',NULL,25,25,NULL,9),
('','medical_remark','In view of the presence of Haemoglobin variants suggested Hb Electrophoresis for further evaluation.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.176774+05:30',NULL,25,25,NULL,17),
('','medical_remark','Suggest Syphilis antibodies for confirmation.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.228217+05:30',NULL,25,25,NULL,18),
('','medical_remark','Rapid card shows indication of malarial antigen, however smear does not show parasite because of low density. Advice to repeat with a repeat sample at the height of fever.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.257703+05:30',NULL,25,25,NULL,19),
('','medical_remark','Kindly correlate clinically and with Vitamin B12 and folic acid levels for further evaluation.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.309867+05:30',NULL,25,25,NULL,20),
('','medical_remark','Dimorphic anaemia. Kindly correlate clinically and advise iron studies, Vitamin B12 and folic acid levels for further evaluation.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.363335+05:30',NULL,25,25,NULL,21),
('','medical_remark','Kindly correlate clinically and advise iron studies and haemoglobin electrophoresis as Mentzer index <14 and Sehgal index <972 for further evaluation.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:03.229646+05:30',NULL,25,25,NULL,49),
('','medical_remark','Kindly correlate clinically and advise free T3, free T4 and Anti TPO levels for further evaluation. Recommended fasting status sample for confirmation if clinically indicated.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:03.258147+05:30',NULL,25,25,NULL,50),
('','medical_remark','Kindly correlate clinically and advise to rule out cross reaction with other infections as IgM can be non-specific.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.394099+05:30',NULL,25,25,NULL,22),
('','medical_remark','In view of Typhidot being positive and Widal being negative suggest Blood Culture to rule out other Salmonella infection.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.42185+05:30',NULL,25,25,NULL,23),
('','medical_remark','Kindly correlate clinically and advise to rule out pre-analytical errors with a repeat sample.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.46404+05:30',NULL,25,25,NULL,24),
('','medical_remark','Rechecked with a repeat sample. If not correlating clinically, suggested retesting on a fresh sample. ','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.49837+05:30',NULL,25,25,NULL,25),
('','medical_remark','RBC agglutination noted. Sample run after 30 min of incubation. Kindly correlate clinically and advise to rule out cold agglutinins causing infections.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.5249+05:30',NULL,25,25,NULL,26),
('','medical_remark','Drug testing: This is a screening test. Kindly confirm with GC/MS if indicated.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.55002+05:30',NULL,25,25,NULL,27),
('','medical_remark','Sample lipemic. Rechecked in dilution. Kindly correlate clinically and advise repeat sample if clinically indicated.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.576384+05:30',NULL,25,25,NULL,28),
('','medical_remark','Rechecked with repeat sample and ECLIA/rapid test and suggested viral load by PCR for further evaluation.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.603883+05:30',NULL,25,25,NULL,29),
('','medical_remark','Kindly correlate clinically with diet and therapeutic history and advise follow up.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:03.012016+05:30',NULL,25,25,NULL,42),
('','medical_remark','Kindly correlate clinically with therapeutic history and advise follow up.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:03.103818+05:30',NULL,25,25,NULL,45),
('','medical_remark','Approximate manual count. Few small clumps noted. Advise follow up and if not correlating clinically, suggested retesting on a fresh sample.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:01.79053+05:30',NULL,25,25,NULL,10),
('','medical_remark','Manual review done. Clumps present/Fibrin strands noted. Platelets appear adequate, approximate count given. Kindly correlate clinically and suggest a repeat sample to rule out pre-analytical errors if indicated.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:01.832184+05:30',NULL,25,25,NULL,11),
('','medical_remark','TPHA positive: Kindly correlate clinically and advise dark field microscopy for further confirmation if indicated.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.630408+05:30',NULL,25,25,NULL,30),
('','medical_remark','Dengue ELISA: NS1 or IgM suspicious: Kindly correlate clinically and advise to repeat after a few days for confirmation','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.656212+05:30',NULL,25,25,NULL,31),
('','medical_remark','Values lower than fasting glucose is a normal phenomenon and suggests a surge of insulin. Kindly correlate clinically with diet or therapeutic history and advise follow up.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.978161+05:30',NULL,25,25,NULL,41),
('','medical_remark','Kindly correlate clinically and recommended sample collection after 10-12 hours of overnight fasting.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:03.130574+05:30',NULL,25,25,NULL,46),
('','medical_remark','Pooled Prolactin - 3 samples collected at 20 min intervals and serum is pooled and values obtained','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.687922+05:30',NULL,25,25,NULL,32),
('','medical_remark','Suggested pooled prolactin for further evaluation if clinically indicated.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.719517+05:30',NULL,25,25,NULL,33),
('','medical_remark','Critical values. Tried calling but not responding. Please contact your doctor immediately.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.747267+05:30',NULL,25,25,NULL,34),
('','medical_remark','HbA1c <3.4 % - Values are beyond linearity range and cannot be reported. Recommended fructosamine levels and variant study if clinically indicated.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.789254+05:30',NULL,25,25,NULL,35),
('','medical_remark','Kindly correlate clinically and with therapeutic history and advise to rule out other skin lesions/infections','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.827186+05:30',NULL,25,25,NULL,36),
('','medical_remark','In view of grade 3 reaction on reverse grouping suggested cross-matching in case of transfusion.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.859945+05:30',NULL,25,25,NULL,37),
('','medical_remark','Kindly repeat after 6 months to correlate with reverse blood grouping','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.887695+05:30',NULL,25,25,NULL,38),
('','medical_remark','It is presumed that patient counselling is done at the referring centre for received samples','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.915104+05:30',NULL,25,25,NULL,39),
('','medical_remark','Positivity in IgG ELISA indicates old infection. Dengue card can show negativity due to higher threshold value of a card.','2024-08-03 16:29:36.354806+05:30','2024-09-25 11:30:02.947543+05:30',NULL,25,25,NULL,40),
('','peripheral_smear','<b>NNBP</b><br/><br><b>RBCs</b>:  Normocytic Normochromic  <br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen  <br/><br><b>Platelets</b>:  Normal in number with normal morphology  <br/><br><b>Haemoparasites</b>:  Not seen  <br/><br><b>Impression</b>:  NORMOCYTIC NORMOCHROMIC BLOOD PICTURE  <br/><br>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>MHA</b><br/><br><b>RBCs</b>:  Show mild to moderate anisopoikilocytosis with normocytic hypochromic and microcytic hypochromic RBCs. Few elliptocytes seen. <br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Normal in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  MICROCYTIC HYPOCHROMIC ANAEMIA <br/><br><b>Comments</b>:  Kindly correlate clinically and with iron studies for further evaluation <br/><br><br>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>Macrocytic anaemia</b><br/><br><b>RBCs</b>:  Show mild anisopoikilocytosis with many macrocytes and normocytes. Few hypochromic RBCs and elliptocytes seen. <br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. Few hypersegmented neutrophils seen. No abnormal or immature cells seen. <br/><br><b>Platelets</b>:  Normal in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  MACROCYTIC ANAEMIA <br/><br><b>Comments</b>:  Kindly correlate clinically and with Vitamin B12 and folic acid levels for further evaluation. <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Normocytic Normochromic </b><br/><br><b>WBCs</b>:  Total count normal in number with increased the number of neutrophils. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Decreased in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC NORMOCHROMIC BLOOD PICTURE WITH RELATIVE NEUTROPHILIA AND THROMBOCYTOPENIA. <br/><br>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>DIMORPHIC ANAEMIA:</b><br/><br><b>RBCs</b>:  Show mild anisopoikilocytosis with a mixed population of macrocytes and microcytes. Few normocytes seen. Few hypochromic RBCs and elliptocytes seen. <br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Normal in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  DIMORPHIC ANAEMIA <br/><br><b>Comments</b>:  Kindly correlate clinically and with iron studies, Vitamin B12 and folic acid levels for further evaluation. <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>NNA:</b><br/><br><b>RBCs</b>:  Predominantly normocytic normochromic RBCs. Few hypochromic RBCs and elliptocytes seen. <br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Normal in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC NORMOCHROMIC ANAEMIA. <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>NHA:</b><br/><br><b>RBCs</b>:  Predominantly normocytic hypochromic RBCs. Few elliptocytes seen. <br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Normal in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC HYPOCHROMIC ANAEMIA. <br/><br><b>Comment</b>:  Kindly correlate clinically. <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>NNBP WITH </b><br/><br><b>RBCs</b>:  Normocytic Normochromic ,Immature Rouleaux formation noted. <br/><br><b>WBCs</b>:  Total count normal in number with increase in the number of lymphocytes. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Normal in number with large platelets present. <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC NORMOCHROMIC BLOOD PICTURE WITH LEUCOPENIA RELATIVE  <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>MHBP</b><br/><br><b>RBCs</b>:  Show mild to moderate anisopoikilocytosis with normocytic hypochromic and microcytic hypochromic RBCs. Few elliptocytes seen. <br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Decreased in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  MICROCYTIC HYPOCHROMIC BLOOD PICTURE WITH THROMBOCYTOPENIA. <br/><br><b>Comments</b>:  Kindly correlate clinically and advise iron studies for further evaluationLYMPHOCYTOSIS AND THROMBOCYTOPENIA. <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Show mild to moderate anisopoikilocytosis with normocytic hypochromic and microcytic hypochromic RBCs. Few elliptocytes seen.</b><br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Decreased in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  MICROCYTIC HYPOCHROMIC ANAEMIA WITH THROMBOCYTOPENIA <br/><br><b>Comments</b>:  Kindly correlate clinically and advise iron studies for further evaluation. <br/><br>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Show mild to moderate anisopoikilocytosis with normocytic hypochromic and microcytic hypochromic RBCs. Few elliptocytes seen.</b><br/><br><b>WBCs</b>:  Total count severely decreased in number with normal distribution. Neutrophils show shift to left. <br/><br><b>Platelets</b>:  Decreased in number with few giant platelets <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  MICROCYTIC HYPOCHROMIC ANAEMIA WITH SEVERE LEUCOPENIA AND THROMBOCYTOPENIA. <br/><br><b>Comments</b>:  Kindly correlate clinically and advise iron studies for further evaluation and follow up <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Predominantly normocytic normochromic RBCs. Few microcytic hypochromic RBCs and elliptocytes seen.</b><br/><br><b>WBCs</b>:  Total count increased in number with increase in neutrophils and have  normal morphology.  No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Normal in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC HYPOCHROMIC ANAEMIA WITH NEUTROPHILIC LEUCOCYTOSIS. <br/><br><b>Comment</b>:  Kindly correlate clinically. <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Normocytic Normochromic </b><br/><br><b>WBCs</b>:  Total count decreased in number with normal morphology and distribution. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Normal in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC NOR <br/><b>Impression</b>:  MICROCYTIC HYPOCHROMIC ANAEMIA WITH EOSINOPHILIA AND LYMPHOCYTOSIS. MOCHROMIC BLOOD PICTURE WITH MILD LEUCOPENIA <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Predominantly normocytic normochromic RBCs. Few hypochromic RBCs and elliptocytes seen.</b><br/><br><b>WBCs</b>:  Total count decreased  in number with normal morphology and distribution. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Decreased  in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC NORMOCHROMIC ANAEMIA WITH LEUCOPENIA AND  THROMBOCYTOPENIA. <br/><br><b>Comment</b>:  Kindly correlate clinically. <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Predominantly normocytic normochromic RBCs. Few hypochromic RBCs and elliptocytes seen.</b><br/><br><b>WBCs</b>:  Total count decreased  in number with normal morphology and distribution. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Decreased  in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC NORMOCHROMIC ANAEMIA WITH LEUCOPENIA AND  THROMBOCYTOPENIA. <br/><br><b>Comment</b>:  Kindly correlate clinically. <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Normocytic Normochromic </b><br/><br><b>WBCs</b>:  Total count normal in number with relative increase in  the number of lymphocytes. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Normal  in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC NORMOCHROMIC BLOOD PICTURE WITH RELATIVE LYMPHOCYTOSIS <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Show mild to moderate anisopoikilocytosis with normocytic hypochromic and microcytic hypochromic RBCs. Few elliptocytes seen.</b><br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Normal in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  MICROCYTIC HYPOCHROMIC BLOOD PICTURE. <br/><br><b>Comments</b>:  Kindly correlate clinically and advise iron studies and haemoglobin electrophoresis as Mentzer index <14 and Sehgal index <972 for further evaluation <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Normocytic Normochromic</b><br/><br><b>WBCs</b>:  Total count increased in number with normal morphology and distribution. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Normal in number with normal morphology  <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC NORMOCHROMIC BLOOD PICTURE WITH LEUCOCYTOSIS <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Normocytic Normochromic</b><br/><br><b>WBCs</b>:  Total count normal in number with mildly increase in the number of eosinophils. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Decreased in number with normal morphology  <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC NORMOCHROMIC BLOOD PICTURE WITH MILD EOSINOPHILIA AND THROMBOCYTOPENIA <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Normocytic hypochromic</b><br/><br><b>WBCs</b>:  Total count decrease in number with relative increase in  the number of lymphocytes. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Decreased in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC NORMOCHROMIC ANAEMIA WITH LEUCOPENIA AND RELATIVE LYMPHOCYTOSIS AND THROMBOCYTOPENIA. <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Normocytic Normochromic</b><br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Decreased in number with few macroplatelets. <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC NORMOCHROMIC BLOOD PICTURE WITH THROMBOCYTOPENIA. <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Normocytic Normochromic</b><br/><br><b>WBCs</b>:  Total count normal in number with increased the number of neutrophils. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Normal in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC NORMOCHROMIC BLOOD PICTURE WITH NEUTROPHILIA <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Normocytic Normochromic</b><br/><br><b>WBCs</b>:  Total count normal in number with increase in eosinophils and have normal morphology . No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Normal in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC NORMOCHROMIC BLOOD PICTURE WITH EOSINOPHILIA. <br/><br><b>Kindly correlate clinically and advise IgE levels for further evaluation.</b>: undefined <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Normocytic Normochromic</b><br/><br><b>WBCs</b>:  Total count decreased in number with normal morphology and distribution. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Normal in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC NORMOCHROMIC BLOOD PICTURE WITH LEUCOPENIA <br/><br><b>Variants detected.Kindly correlate clinically and advise hemoglobin electrophoresis for detection of any hemoglobin variant.</b>: undefined <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Normocytic Normochromic</b><br/><br><b>WBCs</b>:  Total count increase in number with increase in the number of lymphocytes. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Normal in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  NORMOCYTIC NORMOCHROMIC BLOOD PICTURE WITH LEUCOCYTOSIS AND LYMPHOCYTOSIS <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('','peripheral_smear','<b>RBCs: Show mild  anisopoikilocytosis with normocytic hypochromic and microcytic hypochromic RBCs. Few elliptocytes seen.</b><br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen <br/><br><b>Platelets</b>:  Normal in number with normal morphology <br/><br><b>Haemoparasites</b>:  Not seen <br/><br><b>Impression</b>:  MICROCYTIC HYPOCHROMIC BLOOD PICTURE <br/>','2024-08-24 21:53:09.485456+05:30','2024-08-24 21:53:09.485456+05:30','2024-09-12 13:28:28.470984+05:30',0,0,NULL,0),
('Normocytic Normochromic Blood Picture','peripheral_smear','<b>RBCs</b>:  Normocytic Normochromic<br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen<br/><br><b>Platelets</b>:  Normal in number with normal morphology<br/><br><b>Haemoparasites</b>:  Not seen<br/><br><b>Impression: NORMOCYTIC NORMOCHROMIC BLOOD PICTURE</b><br/>','2024-09-12 13:28:28.470984+05:30','2024-09-25 10:59:04.954244+05:30',NULL,25,25,NULL,1),
('Microcytic Hypochromic Anaemia','peripheral_smear','<b>RBCs</b>:  Show mild to moderate anisopoikilocytosis with normocytic hypochromic and microcytic hypochromic RBCs. Few elliptocytes seen.<br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen<br/><br><b>Platelets</b>:  Normal in number with normal morphology<br/><br><b>Haemoparasites</b>:  Not seen<br/><br><b>Impression: MICROCYTIC HYPOCHROMIC ANAEMIA</b><br/><br><b>Comments</b>:  Kindly correlate clinically and with iron studies for further evaluation<br/>','2024-09-12 13:28:28.470984+05:30','2024-09-25 10:59:04.97988+05:30',NULL,25,25,NULL,2),
('Normocytic Normochromic Anaemia','peripheral_smear','<b>RBCs</b>:  Predominantly normocytic normochromic RBCs. Few hypochromic RBCs and elliptocytes seen.<br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen<br/><br><b>Platelets</b>:  Normal in number with normal morphology<br/><br><b>Haemoparasites</b>:  Not seen<br/><br><b>Impression: NORMOCYTIC NORMOCHROMIC ANAEMIA.</b><br/>','2024-09-12 13:28:28.470984+05:30','2024-09-25 10:59:05.004091+05:30',NULL,25,25,NULL,3),
('Normocytic Hypochromic Anaemia','peripheral_smear','<b>RBCs</b>:  Predominantly normocytic hypochromic RBCs. Few elliptocytes seen.<br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen<br/><br><b>Platelets</b>:  Normal in number with normal morphology<br/><br><b>Haemoparasites</b>:  Not seen<br/><br><b>Impression: NORMOCYTIC HYPOCHROMIC ANAEMIA.</b><br/><br><b>Comment</b>:  Kindly correlate clinically and suggest follow-up<br/>','2024-09-12 13:28:28.470984+05:30','2024-09-25 10:59:05.028065+05:30',NULL,25,25,NULL,4),
('Normocytic Normochromic Blood Picture With Leucocytosis','peripheral_smear','<b>RBCs</b>:  Normocytic Normochromic<br/><br><b>WBCs</b>:  Total count increased in number with normal morphology and distribution. No abnormal or immature cells seen<br/><br><b>Platelets</b>:  Normal in number with normal morphology<br/><br><b>Haemoparasites</b>:  Not seen<br/><br><b>Impression: NORMOCYTIC NORMOCHROMIC BLOOD PICTURE WITH LEUCOCYTOSIS</b><br/><br><b>Comment</b>:  Kindly correlate clinically and suggest follow-up.<br/>','2024-09-12 13:28:28.470984+05:30','2024-09-25 10:59:05.092498+05:30',NULL,25,25,NULL,5),
('Normocytic Normochromic Blood Picture With Neutrophilic Leucocytosis','peripheral_smear','<b>RBCs</b>:  Normocytic Normochromic<br/><br><b>WBCs</b>:  Total count increased in number with increased neutrophils and normal morphology. Toxic granules noted. No abnormal or immature cells seen<br/><br><b>Platelets</b>:  Normal in number with normal morphology<br/><br><b>Haemoparasites</b>:  Not seen<br/><br><b>Impression: NORMOCYTIC NORMOCHROMIC BLOOD PICTURE WITH NEUTROPHILIC LEUCOCYTOSIS.</b><br/><br><b>Comment</b>:  Kindly correlate clinically and advise follow up<br/>','2024-09-12 13:28:28.470984+05:30','2024-09-25 10:59:05.117217+05:30',NULL,25,25,NULL,6),
('Normocytic Normochromic Blood Picture With Relative Lymphocytosis','peripheral_smear','<b>RBCs</b>:  Normocytic Normochromic<br/><br><b>WBCs</b>:  Total count normal in number with relative increase in  the number of lymphocytes. No abnormal or immature cells seen<br/><br><b>Platelets</b>:  Normal  in number with normal morphology<br/><br><b>Haemoparasites</b>:  Not seen<br/><br><b>Impression: NORMOCYTIC NORMOCHROMIC BLOOD PICTURE WITH RELATIVE LYMPHOCYTOSIS</b><br/><br><b>Comment</b>:  Kindly correlate clinically and suggest follow-up.<br/>','2024-09-12 13:28:28.470984+05:30','2024-09-25 10:59:05.141447+05:30',NULL,25,25,NULL,7),
('Normocytic Normochromic Blood Picture With Thrombocytopenia','peripheral_smear','<b>RBCs</b>:  Normocytic Normochromic.<br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen<br/><br><b>Platelets</b>:  Decreased in number with large platelets present.<br/><br><b>Haemoparasites</b>:  Not seen<br/><br><b>Impression: NORMOCYTIC NORMOCHROMIC BLOOD PICTURE WITH THROMBOCYTOPENIA</b><br/><br><b>Comment</b>:  Kindly correlate clinically and suggest follow-up.<br/>','2024-09-12 13:28:28.470984+05:30','2024-09-25 10:59:05.167681+05:30',NULL,25,25,NULL,8),
('Normocytic Normochromic Blood Picture With Eosinophilia','peripheral_smear','<b>RBCs</b>:  Normocytic Normochromic<br/><br><b>WBCs</b>:  Total count normal in number with increased number of eosinophils. No abnormal or immature cells seen<br/><br><b>Platelets</b>:  Normal in number with normal morphology<br/><br><b>Haemoparasites</b>:  Not seen<br/><br><b>Impression: NORMOCYTIC NORMOCHROMIC BLOOD PICTURE WITH EOSINOPHILIA</b><br/><br><b>Comment</b>:  Kindly correlate clinically and with IgE levels for further evaluation.<br/>','2024-09-12 13:28:28.470984+05:30','2024-09-25 10:59:05.191551+05:30',NULL,25,25,NULL,9),
('Microcytic Hypochromic Blood Picture','peripheral_smear','<b>RBCs</b>:  Show mild to moderate anisopoikilocytosis with normocytic hypochromic and microcytic hypochromic RBCs. Few elliptocytes and target cells seen.<br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen<br/><br><b>Platelets</b>:  Normal in number with normal morphology<br/><br><b>Haemoparasites</b>:  Not seen<br/><br><b>Impression: MICROCYTIC HYPOCHROMIC BLOOD PICTURE</b><br/><br><b>Comments</b>:  Kindly correlate clinically and advise iron studies and haemoglobin electrophoresis as Mentzer index <14 and Sehgal index <972 for further evaluation<br/>','2024-09-12 13:28:28.470984+05:30','2024-09-25 10:59:05.215849+05:30',NULL,25,25,NULL,10),
('Microcytic Hypochromic Anaemia With Thrombocytosis','peripheral_smear','<b>RBCs</b>:  Show mild to moderate anisopoikilocytosis with normocytic hypochromic and microcytic hypochromic RBCs. Few elliptocytes seen.<br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen<br/><br><b>Platelets</b>:  Increased in number with normal morphology<br/><br><b>Haemoparasites</b>:  Not seen<br/><br><b>Impression: MICROCYTIC HYPOCHROMIC ANAEMIA WITH THROMBOCYTOSIS</b><br/><br><b>Comments</b>:  Kindly correlate clinically and advise iron studies for further evaluation.<br/>','2024-09-12 13:28:28.470984+05:30','2024-09-25 10:59:05.239364+05:30',NULL,25,25,NULL,11),
('Normocytic Normochromic Anaemia With Leucopenia And  Thrombocytopenia','peripheral_smear','<b>RBCs</b>:  Predominantly normocytic normochromic RBCs. Few hypochromic RBCs and elliptocytes seen.<br/><br><b>WBCs</b>:  Total count decreased in number with normal morphology and distribution. No abnormal or immature cells seen<br/><br><b>Platelets</b>:  Decreased  in number with normal morphology<br/><br><b>Haemoparasites</b>:  Not seen<br/><br><b>Impression: NORMOCYTIC NORMOCHROMIC ANAEMIA WITH LEUCOPENIA AND  THROMBOCYTOPENIA.</b><br/><br><b>Comment</b>:  Pancytopenic Picture<br/>','2024-09-12 13:28:28.470984+05:30','2024-09-25 10:59:05.263049+05:30',NULL,25,25,NULL,12),
('Macrocytic Anaemia','peripheral_smear','<b>RBCs</b>:  Show mild anisopoikilocytosis with many macrocytes and normocytes. Few hypochromic RBCs and elliptocytes seen.<br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. Few hypersegmented neutrophils seen. No abnormal or immature cells seen.<br/><br><b>Platelets</b>:  Normal in number with normal morphology<br/><br><b>Haemoparasites</b>:  Not seen<br/><br><b>Impression: MACROCYTIC ANAEMIA</b><br/><br><b>Comments</b>:  Kindly correlate clinically and with Vitamin B12 and folic acid levels for further evaluation.<br/>','2024-09-12 13:28:28.470984+05:30','2024-09-25 10:59:05.289236+05:30',NULL,25,25,NULL,13),
('Dimorphic Anaemia','peripheral_smear','<b>RBCs</b>:  Show mild anisopoikilocytosis with a mixed population of macrocytes and microcytes. Few normocytes seen. Few hypochromic RBCs and elliptocytes seen.<br/><br><b>WBCs</b>:  Total count normal in number with normal morphology and distribution. No abnormal or immature cells seen<br/><br><b>Platelets</b>:  Normal in number with normal morphology<br/><br><b>Haemoparasites</b>:  Not seen<br/><br><b>Impression: DIMORPHIC ANAEMIA</b><br/><br><b>Comments</b>:  Kindly correlate clinically and with iron studies, Vitamin B12 and folic acid levels for further evaluation.<br/>','2024-09-12 13:28:28.470984+05:30','2024-09-25 10:59:05.313111+05:30',NULL,25,25,NULL,14),
('Normocytic Normochromic','peripheral_smear','<b>RBCs</b>:  Normocytic Normochromic<br/><br><b>WBCs</b>:  Total count increased in number. 1% Blasts noted with high N<br/><br><b>Platelets</b>:  Normal/Decreased in number<br/><br><b>Haemoparasites</b>:  Not seen<br/><br><b>Impression/Comment : In view of Atypical cells/ Blasts suggested Bone Marrow Aspiration and Flow Cytometry for further evaluation.</b><br/>','2024-09-12 13:28:28.470984+05:30','2024-09-25 10:59:05.337539+05:30',NULL,25,25,NULL,15),
('','rerun_reason','Value confirmation','2024-09-13 09:59:27.292543+05:30','2024-09-25 10:48:46.248891+05:30',NULL,25,25,NULL,6),
('','rerun_reason','History required','2024-09-13 09:59:27.292543+05:30','2024-09-25 10:48:46.27402+05:30',NULL,25,25,NULL,7),
('','rerun_reason','Pic required','2024-09-24 12:03:53.483538+05:30','2024-09-25 10:48:46.300036+05:30',NULL,25,25,NULL,8),
('','rerun_reason','Slide required','2024-09-24 12:04:03.103976+05:30','2024-09-25 10:48:46.324008+05:30',NULL,25,25,NULL,9),
('','withheld_reason','Lab communication required','2024-09-24 12:10:15.326553+05:30','2024-09-25 10:49:45.886066+05:30',NULL,25,25,NULL,2),
('','medical_remark','Manually checked.','2024-09-24 12:12:15.934746+05:30','2024-09-25 11:30:01.46624+05:30',NULL,25,25,NULL,2),
('','medical_remark','Reactive lymphocytes noted.','2024-09-24 12:12:15.96051+05:30','2024-09-25 11:30:01.675951+05:30',NULL,25,25,NULL,7),
('','medical_remark','Values rechecked.','2024-09-24 12:12:15.984891+05:30','2024-09-25 11:30:03.038441+05:30',NULL,25,25,NULL,43),
('','medical_remark','Values rechecked in dilution.','2024-09-24 12:12:16.01151+05:30','2024-09-25 11:30:03.07873+05:30',NULL,25,25,NULL,44),
('','medical_remark','Direct value','2024-09-24 12:12:16.094042+05:30','2024-09-25 11:30:03.163036+05:30',NULL,25,25,NULL,47);


-- migrate:down
-- write rollback statements below this line
