DROP INDEX `performers_checksum_unique`;
DROP INDEX `index_performers_on_name`;
DROP INDEX `index_performers_on_checksum`;
ALTER TABLE `performers` RENAME TO `temp_old_performers`;
CREATE TABLE `performers` (
  `id` integer not null primary key autoincrement,
  `checksum` varchar(255) not null,
  `name` varchar(255),
  `gender` varchar(20),
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
  `updated_at` datetime not null,
  `image` blob not null
);
CREATE UNIQUE INDEX `performers_checksum_unique` on `performers` (`checksum`);
CREATE INDEX `index_performers_on_name` on `performers` (`name`);
INSERT INTO `performers` (
  `id`,
  `checksum`,
  `name`,
  `gender`,
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
  `updated_at`,
  `image`
)
SELECT 
  `id`,
  `checksum`,
  `name`,
  `gender`,
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
  `updated_at`,
  `image`
FROM `temp_old_performers`;

DROP INDEX `index_performers_scenes_on_scene_id`;
DROP INDEX `index_performers_scenes_on_performer_id`;
ALTER TABLE performers_scenes RENAME TO temp_old_performers_scenes;
CREATE TABLE `performers_scenes` (
  `performer_id` integer,
  `scene_id` integer,
  foreign key(`performer_id`) references `performers`(`id`),
  foreign key(`scene_id`) references `scenes`(`id`)
);
CREATE INDEX `index_performers_scenes_on_scene_id` on `performers_scenes` (`scene_id`);
CREATE INDEX `index_performers_scenes_on_performer_id` on `performers_scenes` (`performer_id`);
INSERT INTO `performers_scenes` (
  `performer_id`,
  `scene_id`
)
SELECT 
  `performer_id`,
  `scene_id`
FROM `temp_old_performers_scenes`;

DROP TABLE `temp_old_performers`;
DROP TABLE `temp_old_performers_scenes`;
