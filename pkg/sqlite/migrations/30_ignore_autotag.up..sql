ALTER TABLE `performers` ADD COLUMN `ignore_auto_tag` boolean not null default '0'; 
ALTER TABLE `studios` ADD COLUMN `ignore_auto_tag` boolean not null default '0';
ALTER TABLE `tags` ADD COLUMN `ignore_auto_tag` boolean not null default '0';