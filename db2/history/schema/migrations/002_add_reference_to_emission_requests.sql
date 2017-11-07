-- +migrate Up
ALTER TABLE history_emission_requests ADD COLUMN reference VARCHAR(64) NOT NULL;

-- +migrate Down
ALTER TABLE history_emission_requests DROP COLUMN reference;