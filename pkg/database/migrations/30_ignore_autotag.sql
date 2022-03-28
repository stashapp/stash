ALTER TABLE `performers` ADD COLUMN `auto_tag_ignored` boolean not null default '0'; 
ALTER TABLE `studios` ADD COLUMN `auto_tag_ignored` boolean not null default '0';
ALTER TABLE `tags` ADD COLUMN `auto_tag_ignored` boolean not null default '0';