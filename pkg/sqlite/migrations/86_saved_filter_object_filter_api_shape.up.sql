PRAGMA foreign_keys=OFF;

-- This migration updates the object_filter JSON structure in saved_filters
-- to match the API input shape. The actual logic runs in 86_postmigrate.go.

CREATE TABLE `saved_filters_new` (
  `id` integer not null primary key autoincrement,
  `name` varchar(510) not null,
  `mode` varchar(255) not null,
  `find_filter` blob,
  `object_filter` blob,
  `ui_options` blob
);

INSERT INTO `saved_filters_new`
  (
    `id`,
    `name`,
    `mode`,
    `find_filter`,
    `object_filter`,
    `ui_options`
  )
  SELECT
    `id`,
    `name`,
    `mode`,
    `find_filter`,
    `object_filter`,
    `ui_options`
  FROM `saved_filters`;

DROP INDEX `index_saved_filters_on_mode_name_unique`;
DROP TABLE `saved_filters`;
ALTER TABLE `saved_filters_new` rename to `saved_filters`;

CREATE UNIQUE INDEX `index_saved_filters_on_mode_name_unique` on `saved_filters` (`mode`, `name`);

PRAGMA foreign_keys=ON;
