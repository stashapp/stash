ALTER TABLE `files` ADD COLUMN `missing_since` datetime;
CREATE INDEX `files_missing_since_index` ON `files` (`missing_since`) WHERE `missing_since` IS NOT NULL;

ALTER TABLE `folders` ADD COLUMN `missing_since` datetime;
CREATE INDEX `folders_missing_since_index` ON `folders` (`missing_since`) WHERE `missing_since` IS NOT NULL;