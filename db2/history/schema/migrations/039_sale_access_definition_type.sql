-- +migrate Up

alter table sales add column access_definition_type int not null;

-- +migrate Down

alter table sales drop column if exists access_definition_type default 0;
