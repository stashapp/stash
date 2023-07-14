PRAGMA foreign_keys=OFF;

-- Cleanup old invalid dates
UPDATE `scenes` SET `date` = NULL WHERE `date` = '0001-01-01' OR `date` = '';
UPDATE `galleries` SET `date` = NULL WHERE `date` = '0001-01-01' OR `date` = '';
UPDATE `performers` SET `birthdate` = NULL WHERE `birthdate` = '0001-01-01' OR `birthdate` = '';
UPDATE `performers` SET `death_date` = NULL WHERE `death_date` = '0001-01-01' OR `death_date` = '';

-- Delete scene markers with missing scenes
DELETE FROM `scene_markers` WHERE `scene_id` IS NULL;

-- make scene_id not null
DROP INDEX `index_scene_markers_on_scene_id`;
DROP INDEX `index_scene_markers_on_primary_tag_id`;

CREATE TABLE `scene_markers_new` (
  `id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  `title` VARCHAR(255) NOT NULL,
  `seconds` FLOAT NOT NULL,
  `primary_tag_id` INTEGER NOT NULL,
  `scene_id` INTEGER NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  FOREIGN KEY(`primary_tag_id`) REFERENCES `tags`(`id`),
  FOREIGN KEY(`scene_id`) REFERENCES `scenes`(`id`)
);
INSERT INTO `scene_markers_new` SELECT * FROM `scene_markers`;

DROP TABLE `scene_markers`;
ALTER TABLE `scene_markers_new` RENAME TO `scene_markers`;

CREATE INDEX `index_scene_markers_on_primary_tag_id` ON `scene_markers`(`primary_tag_id`);
CREATE INDEX `index_scene_markers_on_scene_id` ON `scene_markers`(`scene_id`);

-- drop unused scraped items table
DROP TABLE IF EXISTS `scraped_items`;

-- remove checksum from movies
DROP INDEX `movies_checksum_unique`;
DROP INDEX `movies_name_unique`;

CREATE TABLE `movies_new` (
  `id` integer not null primary key autoincrement,
  `name` varchar(255) not null,
  `aliases` varchar(255),
  `duration` integer,
  `date` date,
  `rating` tinyint,
  `studio_id` integer REFERENCES `studios`(`id`) ON DELETE SET NULL,
  `director` varchar(255),
  `synopsis` text,
  `url` varchar(255),
  `created_at` datetime not null,
  `updated_at` datetime not null, 
  `front_image_blob` varchar(255) REFERENCES `blobs`(`checksum`), 
  `back_image_blob` varchar(255) REFERENCES `blobs`(`checksum`)
);

INSERT INTO `movies_new` SELECT `id`, `name`, `aliases`, `duration`, `date`, `rating`, `studio_id`, `director`, `synopsis`, `url`, `created_at`, `updated_at`, `front_image_blob`, `back_image_blob` FROM `movies`;

DROP TABLE `movies`;
ALTER TABLE `movies_new` RENAME TO `movies`;

CREATE UNIQUE INDEX `index_movies_on_name_unique` ON `movies`(`name`);
CREATE INDEX `index_movies_on_studio_id` on `movies` (`studio_id`);

-- remove checksum from studios
DROP INDEX `index_studios_on_checksum`;
DROP INDEX `index_studios_on_name`;
DROP INDEX `studios_checksum_unique`;

CREATE TABLE `studios_new` (
  `id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `url` VARCHAR(255),
  `parent_id` INTEGER DEFAULT NULL CHECK (`id` IS NOT `parent_id`) REFERENCES `studios`(`id`) ON DELETE SET NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `details` TEXT,
  `rating` TINYINT,
  `ignore_auto_tag` BOOLEAN NOT NULL DEFAULT FALSE,
  `image_blob` VARCHAR(255) REFERENCES `blobs`(`checksum`)
);
INSERT INTO `studios_new` SELECT `id`, `name`, `url`, `parent_id`, `created_at`, `updated_at`, `details`, `rating`, `ignore_auto_tag`, `image_blob` FROM `studios`;

DROP TABLE `studios`;
ALTER TABLE `studios_new` RENAME TO `studios`;

CREATE UNIQUE INDEX `index_studios_on_name_unique` ON `studios`(`name`);

PRAGMA foreign_keys=ON;