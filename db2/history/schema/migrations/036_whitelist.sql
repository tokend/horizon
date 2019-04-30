-- +migrate Up

create table account_specific_rules (
  id  bigint not null,
  address text,
  entry_type int not null,
  key jsonb not null,
  primary key (id)
);

-- +migrate Down

drop table if exists account_specific_rules cascade;
