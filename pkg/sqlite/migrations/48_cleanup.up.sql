-- Cleanup old invalid dates
UPDATE `scenes` SET `date` = NULL WHERE `date` = '0001-01-01' OR `date` = '';
UPDATE `galleries` SET `date` = NULL WHERE `date` = '0001-01-01' OR `date` = '';
UPDATE `performers` SET `birthdate` = NULL WHERE `birthdate` = '0001-01-01' OR `birthdate` = '';
UPDATE `performers` SET `death_date` = NULL WHERE `death_date` = '0001-01-01' OR `death_date` = '';

-- Delete scene markers with missing scenes
DELETE FROM `scene_markers` WHERE `scene_id` IS NULL;

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
