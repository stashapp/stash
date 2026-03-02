-- we cannot add basename column directly because we require it to be NOT NULL
-- recreate folders table with basename column
PRAGMA foreign_keys=OFF;

CREATE TABLE `folders_new` (
  `id` integer not null primary key autoincrement,
  `basename` varchar(255) NOT NULL,
  `path` varchar(255) NOT NULL,
  `parent_folder_id` integer,
  `zip_file_id` integer REFERENCES `files`(`id`),
  `mod_time` datetime not null,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`parent_folder_id`) references `folders`(`id`) on delete SET NULL
);

-- copy data from old table to new table, setting basename to path temporarily
INSERT INTO `folders_new` (
    `id`, 
    `basename`, 
    `path`, 
    `parent_folder_id`, 
    `zip_file_id`,
    `mod_time`, 
    `created_at`, 
    `updated_at`
) SELECT 
    `id`, 
    `path`, 
    `path`, 
    `parent_folder_id`, 
    `zip_file_id`,
    `mod_time`, 
    `created_at`, 
    `updated_at`
FROM `folders`;

DROP INDEX IF EXISTS `index_folders_on_parent_folder_id`;
DROP INDEX IF EXISTS `index_folders_on_path_unique`;
DROP INDEX IF EXISTS `index_folders_on_zip_file_id`;
DROP TABLE `folders`;

ALTER TABLE `folders_new` RENAME TO `folders`;

CREATE UNIQUE INDEX `index_folders_on_path_unique` on `folders` (`path`);
CREATE UNIQUE INDEX `index_folders_on_parent_folder_id_basename_unique` on `folders` (`parent_folder_id`, `basename`);
CREATE INDEX `index_folders_on_zip_file_id` on `folders` (`zip_file_id`) WHERE `zip_file_id` IS NOT NULL;
CREATE INDEX `index_folders_on_basename` on `folders` (`basename`);

PRAGMA foreign_keys=ON;