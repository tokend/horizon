-- +migrate Up

create table deferred_payments
(
    id    bigint      not null,
    amount NUMERIC(20, 0) not null,
    source_fixed_fee        NUMERIC(20, 0)              not null,
    source_percent_fee      NUMERIC(20, 0)              not null,
    destination_fixed_fee   NUMERIC(20, 0)              not null,
    destination_percent_fee NUMERIC(20, 0)              not null,
    source_pays_for_dest         BOOLEAN     NOT NULL,
    source_account VARCHAR(56) not null,
    source_balance VARCHAR(56) not null,
    destination_account VARCHAR(56) not null,
    primary key (id)
);


-- +migrate Down

drop table if exists deferred_payments cascade;