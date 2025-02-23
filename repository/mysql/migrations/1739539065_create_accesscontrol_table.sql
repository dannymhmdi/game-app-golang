-- +migrate Up
CREATE TABLE `access_controls` (
                               `id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
                               `actor_id` INT NOT NULL UNIQUE ,
                               `actor_type` ENUM('user','role') NOT NULL,
                               `permission_id` INT NOT NULL ,
                                `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               FOREIGN KEY (`permission_id`) REFERENCES `permissions`(`id`)
);

-- +migrate Down
DROP TABLE  `access_controls`;