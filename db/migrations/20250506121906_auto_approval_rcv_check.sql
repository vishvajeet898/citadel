-- migrate:up
-- write statements below this line

alter table investigation_results
add column auto_approval_failure_reason VARCHAR(50);

-- migrate:down
-- write rollback statements below this line

alter table investigation_results
drop column auto_approval_failure_reason;
