-- +migrate Up
ALTER TABLE `matches`
    MODIFY COLUMN `question_ids` JSON NOT NULL DEFAULT (JSON_ARRAY());

-- +migrate Down

ALTER TABLE `matches`
    ALTER COLUMN `question_ids` DROP DEFAULT;

