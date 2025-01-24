-- +migrate Up
CREATE TABLE users (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       name VARCHAR(255) ,
                       phone_number VARCHAR(20) not null unique,
                       created_at timestamp default current_timestamp
);

-- +migrate Down

drop table users;