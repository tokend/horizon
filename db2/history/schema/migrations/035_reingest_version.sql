-- +migrate Up


create table support_params
(
	reingest_version int not null check (reingest_version >= 0)
);

-- +migrate Down

drop table if exists support_params;
