PRAGMA foreign_keys=OFF;

-- remove filter column
CREATE TABLE `saved_filters_new` (
  `id` integer not null primary key autoincrement,
  `name` varchar(510) not null,
  `mode` varchar(255) not null,
  `find_filter` blob,
  `object_filter` blob,
  `ui_options` blob
);

-- move filter data into find_filter to be migrated in the post-migration
INSERT INTO `saved_filters_new`
  (
    `id`,
    `name`,
    `mode`,
    `find_filter`
  )
  SELECT 
    `id`,
    `name`,
    `mode`,
    `filter`
  FROM `saved_filters`;

DROP INDEX `index_saved_filters_on_mode_name_unique`;
DROP TABLE `saved_filters`;
ALTER TABLE `saved_filters_new` rename to `saved_filters`;

CREATE UNIQUE INDEX `index_saved_filters_on_mode_name_unique` on `saved_filters` (`mode`, `name`);

PRAGMA foreign_keys=ON;
