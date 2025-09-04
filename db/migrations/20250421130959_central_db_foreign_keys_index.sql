-- migrate:up
-- write statements below this line

-- foreign keys

ALTER TABLE "rerun_investigation_results" ADD CONSTRAINT "rerun_investigation_results_rerun_triggered_by_foreign" FOREIGN KEY("rerun_triggered_by") REFERENCES "users"("id");
ALTER TABLE "rerun_investigation_results" ADD CONSTRAINT "rerun_investigation_results_test_details_id_foreign" FOREIGN KEY("test_details_id") REFERENCES "test_details"("id");

ALTER TABLE "investigation_data" ADD CONSTRAINT "investigation_data_investigation_result_id_foreign" FOREIGN KEY("investigation_result_id") REFERENCES "investigation_results"("id");

ALTER TABLE "task_visit_mapping" ADD CONSTRAINT "task_visit_mapping_task_id_foreign" FOREIGN KEY("task_id") REFERENCES "tasks"("id");

ALTER TABLE "remarks" ADD CONSTRAINT "remarks_investigation_result_id_foreign" FOREIGN KEY("investigation_result_id") REFERENCES "investigation_results"("id");

ALTER TABLE "task_pathologist_mapping" ADD CONSTRAINT "task_pathologist_mapping_task_id_foreign" FOREIGN KEY("task_id") REFERENCES "tasks"("id");
ALTER TABLE "task_pathologist_mapping" ADD CONSTRAINT "task_pathologist_mapping_pathologist_id_foreign" FOREIGN KEY("pathologist_id") REFERENCES "users"("id");

ALTER TABLE "attachments" ADD CONSTRAINT "attachments_task_id_foreign" FOREIGN KEY("task_id") REFERENCES "tasks"("id");

ALTER TABLE "tasks" ADD CONSTRAINT "tasks_patient_details_id_foreign" FOREIGN KEY("patient_details_id") REFERENCES "patient_details"("id");

ALTER TABLE "test_details_metadata" ADD CONSTRAINT "test_details_metadata_test_details_id_foreign" FOREIGN KEY("test_details_id") REFERENCES "test_details"("id");

ALTER TABLE "investigation_results" ADD CONSTRAINT "investigation_results_test_details_id_foreign" FOREIGN KEY("test_details_id") REFERENCES "test_details"("id");

ALTER TABLE "task_metadata" ADD CONSTRAINT "task_metadata_task_id_foreign" FOREIGN KEY("task_id") REFERENCES "tasks"("id");

ALTER TABLE "co_authorized_pathologists" ADD CONSTRAINT "co_authorized_pathologists_co_authorized_by_foreign" FOREIGN KEY("co_authorized_by") REFERENCES "users"("id");
ALTER TABLE "co_authorized_pathologists" ADD CONSTRAINT "co_authorized_pathologists_co_authorized_to_foreign" FOREIGN KEY("co_authorized_to") REFERENCES "users"("id");
ALTER TABLE "co_authorized_pathologists" ADD CONSTRAINT "co_authorized_pathologists_task_id_foreign" FOREIGN KEY("task_id") REFERENCES "tasks"("id");

ALTER TABLE "sample_metadata" ADD CONSTRAINT "sample_metadata_sample_id_foreign" FOREIGN KEY("sample_id") REFERENCES "samples"("id");

ALTER TABLE "test_sample_mapping" ADD CONSTRAINT "test_sample_mapping_sample_id_foreign" FOREIGN KEY ("sample_id") REFERENCES "samples"("id");



-- index

CREATE INDEX IF NOT EXISTS task_status_idx ON tasks(status);
CREATE INDEX IF NOT EXISTS task_deleted_at_idx ON tasks(deleted_at);
CREATE INDEX IF NOT EXISTS task_doctor_tat_idx ON tasks(doctor_tat);
CREATE INDEX IF NOT EXISTS task_order_city_idx ON tasks(order_id, city_code);
CREATE INDEX IF NOT EXISTS task_request_id_idx ON tasks(request_id);
CREATE INDEX IF NOT EXISTS task_order_id_idx ON tasks(order_id);
CREATE INDEX IF NOT EXISTS task_lab_id_idx ON tasks(lab_id);
CREATE INDEX IF NOT EXISTS task_order_type_idx ON tasks(order_type);
CREATE INDEX IF NOT EXISTS tasks_status_deleted_idx ON tasks(status, deleted_at);
CREATE INDEX IF NOT EXISTS task_updated_at_idx ON tasks(updated_at);
CREATE INDEX IF NOT EXISTS tasks_oms_order_id_idx ON tasks(oms_order_id);
CREATE INDEX IF NOT EXISTS idx_tasks_patient_status ON tasks(patient_details_id, status);

CREATE INDEX IF NOT EXISTS task_metadata_task_id_idx ON task_metadata(task_id);
CREATE INDEX IF NOT EXISTS task_metadata_is_critical_idx ON task_metadata(is_critical);
CREATE INDEX IF NOT EXISTS task_metadata_last_event_idx ON task_metadata(last_event_sent_at);
CREATE INDEX IF NOT EXISTS task_metadata_deleted_at_idx ON task_metadata(deleted_at);
CREATE INDEX IF NOT EXISTS task_metadata_lower_doctor_name_idx ON task_metadata(lower(doctor_name));
CREATE INDEX IF NOT EXISTS task_metadata_lower_partner_name_idx ON task_metadata(lower(partner_name));
CREATE INDEX IF NOT EXISTS task_metadata_updated_at_idx ON task_metadata(updated_at);

CREATE INDEX IF NOT EXISTS attachments_task_id_idx ON attachments(task_id);
CREATE INDEX IF NOT EXISTS attachments_deleted_at_idx ON attachments(deleted_at);
CREATE INDEX IF NOT EXISTS attachments_updated_at_idx ON attachments(updated_at);

CREATE INDEX IF NOT EXISTS co_authorized_paths_task_id_idx ON co_authorized_pathologists(task_id);
CREATE INDEX IF NOT EXISTS co_authorized_paths_deleted_at_idx ON co_authorized_pathologists(deleted_at);
CREATE INDEX IF NOT EXISTS co_authorized_paths_by_idx ON co_authorized_pathologists(co_authorized_by);
CREATE INDEX IF NOT EXISTS co_authorized_paths_to_idx ON co_authorized_pathologists(co_authorized_to);
CREATE INDEX IF NOT EXISTS co_authorized_paths_updated_at_idx ON co_authorized_pathologists(updated_at);

CREATE INDEX IF NOT EXISTS investigation_data_investigation_id_idx ON investigation_data(investigation_result_id);
CREATE INDEX IF NOT EXISTS investigation_data_deleted_at_idx ON investigation_data(deleted_at);
CREATE INDEX IF NOT EXISTS investigation_data_updated_at_idx ON investigation_data(updated_at);

CREATE INDEX IF NOT EXISTS investigation_results_test_id_idx ON investigation_results(test_details_id);
CREATE INDEX IF NOT EXISTS investigation_results_deleted_at_idx ON investigation_results(deleted_at);
CREATE INDEX IF NOT EXISTS investigation_results_master_investigation_id_idx ON investigation_results(master_investigation_id);
CREATE INDEX IF NOT EXISTS investigation_results_status_idx ON investigation_results(investigation_status);
CREATE INDEX IF NOT EXISTS investigation_results_approved_at_idx ON investigation_results(approved_at);
CREATE INDEX IF NOT EXISTS investigation_results_updated_at_idx ON investigation_results(updated_at);

CREATE INDEX IF NOT EXISTS patient_details_patient_id_idx ON patient_details(system_patient_id);
CREATE INDEX IF NOT EXISTS patient_details_deleted_at_idx ON patient_details(deleted_at);
CREATE INDEX IF NOT EXISTS patient_details_lower_name_idx ON patient_details (lower(name));
CREATE INDEX IF NOT EXISTS patient_details_updated_at_idx ON patient_details(updated_at);

CREATE INDEX IF NOT EXISTS remarks_investigation_id_idx ON remarks(investigation_result_id);
CREATE INDEX IF NOT EXISTS remarks_deleted_at_idx ON remarks(deleted_at);
CREATE INDEX IF NOT EXISTS remarks_updated_at_idx ON remarks(updated_at);

CREATE INDEX IF NOT EXISTS rerun_test_details_test_id_idx ON rerun_investigation_results(test_details_id);
CREATE INDEX IF NOT EXISTS rerun_test_details_deleted_at_idx ON rerun_investigation_results(deleted_at);
CREATE INDEX IF NOT EXISTS rerun_test_details_updated_at_idx ON rerun_investigation_results(updated_at);

CREATE INDEX IF NOT EXISTS task_pathologist_mapping_task_id_idx ON task_pathologist_mapping(task_id);
CREATE INDEX IF NOT EXISTS task_pathologist_mapping_deleted_at_idx ON task_pathologist_mapping(deleted_at);
CREATE INDEX IF NOT EXISTS task_pathologist_mapping_updated_at_idx ON task_pathologist_mapping(updated_at);
CREATE INDEX IF NOT EXISTS task_pathologist_mapping_active_task_id_idx ON task_pathologist_mapping(is_active, task_id);

CREATE INDEX IF NOT EXISTS visit_task_id_idx ON task_visit_mapping(task_id);
CREATE INDEX IF NOT EXISTS visit_deleted_at_idx ON task_visit_mapping(deleted_at);
CREATE INDEX IF NOT EXISTS visit_id_idx ON task_visit_mapping(visit_id);
CREATE INDEX IF NOT EXISTS task_visit_mapping_updated_at_idx ON task_visit_mapping(updated_at);

CREATE INDEX IF NOT EXISTS test_details_status_idx ON test_details(status);
CREATE INDEX IF NOT EXISTS test_details_deleted_at_idx ON test_details(deleted_at);
CREATE INDEX IF NOT EXISTS test_details_task_id_idx ON test_details(task_id);
CREATE INDEX IF NOT EXISTS test_details_oms_test_id_idx ON test_details(oms_test_id);
CREATE INDEX IF NOT EXISTS test_details_updated_at_idx ON test_details(updated_at);
CREATE INDEX IF NOT EXISTS test_details_status_deleted_idx ON test_details(status, deleted_at);
CREATE INDEX IF NOT EXISTS test_details_oms_order_id_idx ON test_details(oms_order_id);
CREATE INDEX IF NOT EXISTS test_details_lab_id_idx ON test_details(lab_id);
CREATE INDEX IF NOT EXISTS test_details_processing_lab_id_idx ON test_details(processing_lab_id);
CREATE INDEX IF NOT EXISTS test_details_lab_eta_idx ON test_details(lab_eta);
CREATE INDEX IF NOT EXISTS test_details_master_package_id_idx ON test_details(master_package_id);
CREATE INDEX IF NOT EXISTS test_details_central_oms_test_id_idx ON test_details(central_oms_test_id);
CREATE INDEX IF NOT EXISTS test_details_cp_enabled_idx ON test_details(cp_enabled);

CREATE INDEX IF NOT EXISTS test_details_metadata_test_id_idx ON test_details_metadata(test_details_id);
CREATE INDEX IF NOT EXISTS test_details_metadata_deleted_at_idx ON test_details_metadata(deleted_at);
CREATE INDEX IF NOT EXISTS test_details_metadata_is_completed_idx ON test_details_metadata(is_completed_in_oms);
CREATE INDEX IF NOT EXISTS test_details_metadata_updated_at_idx ON test_details_metadata(updated_at);

CREATE INDEX IF NOT EXISTS users_email_idx ON users(email);
CREATE INDEX IF NOT EXISTS users_deleted_at_idx ON users(deleted_at);
CREATE INDEX IF NOT EXISTS users_type_idx ON users(user_type);

CREATE INDEX IF NOT EXISTS idx_external_investigation_results_contact_id ON external_investigation_results(contact_id);
CREATE INDEX IF NOT EXISTS idx_external_investigation_results_master_investigation_id ON external_investigation_results(master_investigation_id);
CREATE INDEX IF NOT EXISTS idx_external_investigation_results_system_external_report_id ON external_investigation_results(system_external_report_id);
CREATE INDEX IF NOT EXISTS idx_external_investigation_results_master_investigation_method_mapping_id ON external_investigation_results(master_investigation_method_mapping_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_external_investigation_results_system_external_investigation_result_id ON external_investigation_results(system_external_investigation_result_id);
CREATE INDEX IF NOT EXISTS idx_external_investigation_results_loinc_code ON external_investigation_results(loinc_code);
CREATE INDEX IF NOT EXISTS idx_external_investigation_results_reported_at ON external_investigation_results(reported_at);
CREATE INDEX IF NOT EXISTS idx_external_investigation_results_deleted_at ON external_investigation_results(deleted_at);
CREATE INDEX IF NOT EXISTS idx_external_investigation_results_abnormality ON external_investigation_results(abnormality);
CREATE INDEX IF NOT EXISTS ids_external_investigation_updated_at ON external_investigation_results(updated_at);

CREATE INDEX IF NOT EXISTS order_details_deleted_at_idx ON order_details(deleted_at);
CREATE INDEX IF NOT EXISTS order_details_oms_order_id_idx ON order_details(oms_order_id);
CREATE INDEX IF NOT EXISTS order_details_trf_id_idx ON order_details(trf_id);
CREATE INDEX IF NOT EXISTS order_details_updated_at_ids ON order_details(updated_at);
CREATE INDEX IF NOT EXISTS order_details_servicing_lab_id_idx ON order_details(servicing_lab_id);

CREATE INDEX IF NOT EXISTS samples_deleted_at_idx ON samples(deleted_at);
CREATE INDEX IF NOT EXISTS samples_oms_order_id_idx ON samples(oms_order_id);
CREATE INDEX IF NOT EXISTS sample_oms_request_id_idx ON samples(oms_request_id);
CREATE INDEX IF NOT EXISTS samples_visit_id_idx ON samples(visit_id);
CREATE INDEX IF NOT EXISTS samples_barcode_idx ON samples(barcode);
CREATE INDEX IF NOT EXISTS samples_status_idx ON samples(status);
CREATE INDEX IF NOT EXISTS samples_sample_number_idx ON samples(sample_number);
CREATE INDEX IF NOT EXISTS samples_vial_type_id_idx ON samples(vial_type_id);
CREATE INDEX IF NOT EXISTS samples_updated_at_ids ON samples(updated_at);
CREATE INDEX IF NOT EXISTS samples_audit_oms_order_id_idx ON samples_audit(oms_order_id);

CREATE INDEX IF NOT EXISTS sample_metadata_deleted_at_idx ON sample_metadata(deleted_at);
CREATE INDEX IF NOT EXISTS sample_metadata_sample_id_idx ON sample_metadata(sample_id);
CREATE INDEX IF NOT EXISTS sample_metadata_oms_order_id_idx ON sample_metadata(oms_order_id);
CREATE INDEX IF NOT EXISTS sample_metadata_task_sequence_idx ON sample_metadata(task_sequence);
CREATE INDEX IF NOT EXISTS sample_metadata_collect_later_reason_idx ON sample_metadata(collect_later_reason);
CREATE INDEX IF NOT EXISTS sample_metadata_updated_at_idx ON sample_metadata(updated_at);
CREATE INDEX IF NOT EXISTS sample_metadata_audit_oms_order_id_idx ON sample_metadata_audit(oms_order_id);

CREATE INDEX IF NOT EXISTS test_sample_mapping_deleted_at_idx ON test_sample_mapping(deleted_at);
CREATE INDEX IF NOT EXISTS test_sample_mapping_sample_id_idx ON test_sample_mapping(sample_id);
CREATE INDEX IF NOT EXISTS test_sample_mapping_oms_test_id_idx ON test_sample_mapping(oms_test_id);
CREATE INDEX IF NOT EXISTS test_sample_mapping_oms_order_id_idx ON test_sample_mapping(oms_order_id);
CREATE INDEX IF NOT EXISTS test_sample_mapping_sample_number_idx ON test_sample_mapping(sample_number);
CREATE INDEX IF NOT EXISTS test_sample_mapping_updated_at_idx ON test_sample_mapping(updated_at);
CREATE INDEX IF NOT EXISTS test_sample_mapping_is_rejected_idx ON test_sample_mapping(is_rejected);

CREATE INDEX IF NOT EXISTS ets_events_test_ids_idx ON ets_events(test_id);


-- migrate:down
-- write rollback statements below this line
