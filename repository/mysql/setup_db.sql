CREATE TABLE users (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       name VARCHAR(255) default 'Daniel',
                       phone_number VARCHAR(20) not null unique
);
