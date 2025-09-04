-- migrate:up
-- write statements below this line

CREATE INDEX investigation_results_active_idx ON investigation_results(test_details_id) WHERE deleted_at IS NULL;
CREATE INDEX test_details_task_active_idx ON test_details(task_id) WHERE deleted_at IS NULL;


-- migrate:down
-- write rollback statements below this line

DROP INDEX investigation_results_active_idx;
DROP INDEX test_details_task_active_idx;
