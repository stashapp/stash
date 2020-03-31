
PRAGMA foreign_keys=off;

-- need to re-create the performers table without the added column.
-- also need re-create the performers_scenes table due to the foreign key

-- rename existing performers table
ALTER TABLE `performers` RENAME TO `performers_old`;
ALTER TABLE `performers_scenes` RENAME TO `performers_scenes_old`;

-- drop the indexes
DROP INDEX IF EXISTS `index_performers_on_name`;
DROP INDEX IF EXISTS `index_performers_on_checksum`;
DROP INDEX IF EXISTS `index_performers_scenes_on_scene_id`;
DROP INDEX IF EXISTS `index_performers_scenes_on_performer_id`;

-- recreate the tables
CREATE TABLE `performers` (
  `id` integer not null primary key autoincrement,
  `image` blob not null,
  `checksum` varchar(255) not null,
  `name` varchar(255),
  `url` varchar(255),
  `twitter` varchar(255),
  `instagram` varchar(255),
  `birthdate` date,
  `ethnicity` varchar(255),
  `country` varchar(255),
  `eye_color` varchar(255),
  `height` varchar(255),
  `measurements` varchar(255),
  `fake_tits` varchar(255),
  `career_length` varchar(255),
  `tattoos` varchar(255),
  `piercings` varchar(255),
  `aliases` varchar(255),
  `favorite` boolean not null default '0',
  `created_at` datetime not null,
  `updated_at` datetime not null
);

CREATE TABLE `performers_scenes` (
  `performer_id` integer,
  `scene_id` integer,
  foreign key(`performer_id`) references `performers`(`id`),
  foreign key(`scene_id`) references `scenes`(`id`)
);

INSERT INTO `performers` 
  SELECT 
  `id`,
  `image`,
  `checksum`,
  `name`,
  `url`,
  `twitter`,
  `instagram`,
  `birthdate`,
  `ethnicity`,
  `country`,
  `eye_color`,
  `height`,
  `measurements`,
  `fake_tits`,
  `career_length`,
  `tattoos`,
  `piercings`,
  `aliases`,
  `favorite`,
  `created_at`,
  `updated_at`
  FROM `performers_old`;

INSERT INTO `performers_scenes`
  SELECT
  `performer_id`,
  `scene_id`
  FROM `performers_scenes_old`;

DROP TABLE `performers_scenes_old`;
DROP TABLE `performers_old`;

-- re-create the indexes after removing the old tables
CREATE INDEX `index_performers_on_name` on `performers` (`name`);
CREATE INDEX `index_performers_on_checksum` on `performers` (`checksum`);
CREATE INDEX `index_performers_scenes_on_scene_id` on `performers_scenes` (`scene_id`);
CREATE INDEX `index_performers_scenes_on_performer_id` on `performers_scenes` (`performer_id`);

PRAGMA foreign_keys=on;
