-- +migrate Up

ALTER TABLE reviewable_request ADD all_tasks INT DEFAULT 0;
ALTER TABLE reviewable_request ADD pending_tasks INT DEFAULT 0;
ALTER TABLE reviewable_request ADD external_details TEXT;

-- +migrate Down

ALTER TABLE reviewable_request DROP COLUMN all_tasks;
ALTER TABLE reviewable_request DROP COLUMN pending_tasks;
ALTER TABLE reviewable_request DROP COLUMN external_details;
