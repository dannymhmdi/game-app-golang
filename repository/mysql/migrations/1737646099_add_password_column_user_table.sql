-- +migrate Up
alter table users add column password varchar(255) not null;

-- +migrate Down
alter table users drop column password;