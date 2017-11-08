-- +migrate Up
ALTER TABLE history_balances ALTER COLUMN asset TYPE character varying(16);
ALTER TABLE history_emission_requests ALTER COLUMN asset TYPE character varying(16);

-- +migrate Down
ALTER TABLE history_balances ALTER COLUMN asset TYPE character varying(6) USING substr(asset, 1, 6);
ALTER TABLE history_emission_requests ALTER COLUMN asset TYPE character varying(6) USING substr(asset, 1, 6);