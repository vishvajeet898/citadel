-- migrate:up
-- write statements below this line

UPDATE "users" SET "user_type" = 'pathologist' WHERE "user_type" = 'Pathologist';
UPDATE "users" SET "user_type" = 'system' WHERE "user_type" = 'System';

-- migrate:down
-- write rollback statements below this line

UPDATE "users" SET "user_type" = 'Pathologist' WHERE "user_type" = 'pathologist';
UPDATE "users" SET "user_type" = 'System' WHERE "user_type" = 'system';
