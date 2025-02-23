-- +migrate Up
CREATE TABLE `permissions` (
    `id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
     `title` VARCHAR(191) NOT NULL ,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
ALTER TABLE  `users` DROP COLUMN `role`;