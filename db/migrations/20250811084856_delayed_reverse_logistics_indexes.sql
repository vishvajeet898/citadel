-- migrate:up
-- write statements below this line

CREATE INDEX IF NOT EXISTS order_details_camp_id_idx ON order_details(camp_id);
CREATE INDEX IF NOT EXISTS order_details_collection_type_idx ON order_details(collection_type);
CREATE INDEX IF NOT EXISTS sample_metadata_collected_at_idx ON sample_metadata(collected_at);
CREATE INDEX IF NOT EXISTS sample_metadata_transferred_at_idx ON sample_metadata(transferred_at);

-- migrate:down
-- write rollback statements below this line

DROP INDEX IF EXISTS order_details_camp_id_idx;
DROP INDEX IF EXISTS order_details_collection_type_idx;
DROP INDEX IF EXISTS sample_metadata_collected_at_idx;
DROP INDEX IF EXISTS sample_metadata_transferred_at_idx;
