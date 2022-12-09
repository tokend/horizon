-- +migrate Up

CREATE TABLE liquidity_pools (
    id                BIGINT NOT NULL,
    account           CHAR(56) NOT NULL,
    token_asset       VARCHAR(16) NOT NULL,
    first_balance     CHAR(56) NOT NULL,
    second_balance    CHAR(56) NOT NULL,
    tokens_amount     NUMERIC(20, 0) NOT NULL,
    first_reserve     NUMERIC(20, 0) NOT NULL,
    second_reserve    NUMERIC(20, 0) NOT NULL,
    first_asset_code  VARCHAR(16) NOT NULL,
    second_asset_code VARCHAR(16) NOT NULL,
    PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE liquidity_pools;