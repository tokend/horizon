-- +migrate Up
ALTER TABLE history_balances DROP COLUMN kyc;
ALTER TABLE history_balances ADD COLUMN exchange_name character varying(64);

-- +migrate Down

ALTER TABLE history_balances DROP COLUMN exchange_name;
ALTER TABLE history_balances ADD COLUMN kyc jsonb;