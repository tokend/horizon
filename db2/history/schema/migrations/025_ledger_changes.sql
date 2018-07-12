-- +migrate Up

CREATE TABLE history_ledger_changes
(
  id            BIGSERIAL NOT NULL CHECK (id >= 0),
  tx_id         BIGINT    NOT NULL CHECK (tx_id >= 0),
  op_id         BIGINT    NOT NULL CHECK (op_id >= 0),
  order_number  INT       NOT NULL CHECK (order_number >= 0),
  effect        INT       NOT NULL,
  entry_type    INT       NOT NULL,
  ledger_key    TEXT      NOT NULL,
  PRIMARY KEY (id)
);

-- +migrate Down

DROP TABLE IF EXISTS history_ledger_changes;