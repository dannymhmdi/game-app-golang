CREATE TABLE users (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       name VARCHAR(255) ,
                       phone_number VARCHAR(20) not null unique,
    password varchar(255) not null,
    created_at timestamp default current_timestamp
);

INSERT INTO users (name, phone_number)
VALUES ('Daniel','09127275236' );
