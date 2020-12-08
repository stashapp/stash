ALTER TABLE `scenes` ADD COLUMN `organized` boolean not null default '0';
ALTER TABLE `images` ADD COLUMN `organized` boolean not null default '0';
ALTER TABLE `galleries` ADD COLUMN `organized` boolean not null default '0';
