-- +migrate Up

create table account_specific_rules (
  id  bigint not null,
  address varchar(64),
  forbids boolean not null,
  entry_type int not null,
  key jsonb not null,
  primary key (id)
);

create table sale_participation (
    id bigint not null,
    participant_id varchar(64) not null,
    sale_id bigint not null,
    base_amount numeric(23,0) not null,
    quote_amount numeric(23,0) not null,
    quote_asset varchar(16) not null,
    base_asset varchar(16) not null,
    price numeric(23,0) not null,
    primary key (id)
);

-- +migrate Down

drop table if exists account_specific_rules cascade;
drop table if exists sale_participation cascade;