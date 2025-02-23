-- +migrate Up
CREATE TABLE `users` (
                       `id` INT AUTO_INCREMENT PRIMARY KEY,
                       `name` VARCHAR(191) ,
                       `phone_number` VARCHAR(20) NOT NULL UNIQUE,
                    `role` ENUM('user','admin') NOT NULL,
                       `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down

DROP TABLE `users`;