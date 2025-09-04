-- migrate:up
-- write statements below this line

ALTER TABLE order_details ADD COLUMN srf_id VARCHAR(255) DEFAULT NULL;
CREATE INDEX IF NOT EXISTS order_details_srf_id_idx ON order_details(srf_id);
CREATE INDEX IF NOT EXISTS test_details_master_test_id_idx ON test_details(master_test_id);
CREATE INDEX idx_order_details_recent_active ON order_details (created_at, order_status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS order_details_created_at_idx ON order_details(created_at);

-- migrate:down
-- write rollback statements below this line

ALTER TABLE order_details DROP COLUMN srf_id;
DROP INDEX IF EXISTS idx_order_details_srf_id;
DROP INDEX IF EXISTS test_details_master_test_id_idx;
