ALTER TABLE `scenes` ADD COLUMN `captioned` boolean not null default '0';
ALTER TABLE `scenes` ADD COLUMN `captions` varchar(255) not null default '';