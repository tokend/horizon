-- +migrate Up

drop table if exists sale_participation cascade;

-- +migrate Down

create table sale_participation (
    id bigint not null,
    participant_id varchar(56) not null,
    sale_id bigint not null,
    base_amount numeric(20,0) not null,
    quote_amount numeric(20,0) not null,
    quote_asset varchar(16) not null,
    base_asset varchar(16) not null,
    price numeric(20,0) not null,
    primary key (id)
);
