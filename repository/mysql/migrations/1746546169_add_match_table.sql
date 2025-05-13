-- +migrate Up
CREATE TABLE `matches` (
                         `id` INT AUTO_INCREMENT PRIMARY KEY,
                         `player_ids` JSON NOT NULL COMMENT 'Stores []uint as JSON array, e.g., [1, 2, 3]',
                         `question_ids` JSON NOT NULL DEFAULT JSON_ARRAY()),
                         `category` ENUM('soccer','history') NOT NULL,
                         `start_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down

DROP TABLE `users`;

