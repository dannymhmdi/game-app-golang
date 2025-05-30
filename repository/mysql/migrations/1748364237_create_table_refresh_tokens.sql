-- +migrate Up
CREATE TABLE `refresh_tokens` (
                               `id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
                               `user_id` INT NOT NULL ,
                                `token_hash` TINYTEXT NOT NULL ,
                                `expires_at` TIMESTAMP NOT NULL ,
                                `revoked` TINYINT(1) NOT NULL DEFAULT 0,
                                `device_info` VARCHAR(256) ,
                                `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE  `refresh_tokens`;