-- +migrate Up
INSERT INTO `permissions` (id , title) VALUES
(1,'user-delete'),
(2,'user-list');
INSERT INTO `access_controls` (`actor_id`,`actor_type`,`permission_id`) values(2,'role',1),
(2,'role',2);

-- +migrate Down
DELETE FROM `access_controls` WHERE id IN (1,2);
DELETE FROM `permissions` WHERE id IN (1,2);