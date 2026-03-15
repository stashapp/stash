-- have to change the type of the career start/end columns so need to recreate the table
PRAGMA foreign_keys=OFF;

CREATE TABLE IF NOT EXISTS "performers_new" (
  `id` integer not null primary key autoincrement,
  `name` varchar(255) not null,
  `disambiguation` varchar(255),
  `gender` varchar(20),
  `birthdate` date,
  `birthdate_precision` TINYINT, 
  `ethnicity` varchar(255),
  `country` varchar(255),
  `eye_color` varchar(255),
  `height` int,
  `measurements` varchar(255),
  `fake_tits` varchar(255),
  `tattoos` varchar(255),
  `piercings` varchar(255),
  `favorite` boolean not null default '0',
  `created_at` datetime not null,
  `updated_at` datetime not null,
  `details` text, 
  `death_date` date, 
  `death_date_precision` TINYINT, 
  `hair_color` varchar(255), 
  `weight` integer, 
  `rating` tinyint, 
  `ignore_auto_tag` boolean not null default '0', 
  `penis_length` float, 
  `circumcised` varchar[10], 
  `career_start` date, 
  `career_start_precision` TINYINT, 
  `career_end` date, 
  `career_end_precision` TINYINT,
  `image_blob` varchar(255) REFERENCES `blobs`(`checksum`)
);

INSERT INTO `performers_new` (
    `id`,
    `name`,
    `disambiguation`,
    `gender`,
    `birthdate`,
    `ethnicity`,
    `country`,
    `eye_color`,
    `height`,
    `measurements`,
    `fake_tits`,
    `tattoos`,
    `piercings`,
    `favorite`,
    `created_at`,
    `updated_at`,
    `details`,
    `death_date`,
    `hair_color`,
    `weight`,
    `rating`,
    `ignore_auto_tag`,
    `image_blob`,
    `penis_length`,
    `circumcised`,
    `birthdate_precision`,
    `death_date_precision`,
    `career_start`,
    `career_end`
) SELECT 
    `id`,
    `name`,
    `disambiguation`,
    `gender`,
    `birthdate`,
    `ethnicity`,
    `country`,
    `eye_color`,
    `height`,
    `measurements`,
    `fake_tits`,
    `tattoos`,
    `piercings`,
    `favorite`,
    `created_at`,
    `updated_at`,
    `details`,
    `death_date`,
    `hair_color`,
    `weight`,
    `rating`,
    `ignore_auto_tag`,
    `image_blob`,
    `penis_length`,
    `circumcised`,
    `birthdate_precision`,
    `death_date_precision`,
    CAST(`career_start` AS TEXT),
    CAST(`career_end` AS TEXT)
FROM `performers`;

DROP INDEX IF EXISTS `performers_name_disambiguation_unique`;
DROP INDEX IF EXISTS `performers_name_unique`;
DROP TABLE `performers`;

ALTER TABLE `performers_new` RENAME TO `performers`;

UPDATE "performers" SET `career_start` = CONCAT(`career_start`, '-01-01'), "career_start_precision" = 2 WHERE "career_start" IS NOT NULL;
UPDATE "performers" SET `career_end` = CONCAT(`career_end`, '-01-01'), "career_end_precision" = 2 WHERE "career_end" IS NOT NULL;  

CREATE UNIQUE INDEX `performers_name_disambiguation_unique` on `performers` (`name`, `disambiguation`) WHERE `disambiguation` IS NOT NULL;
CREATE UNIQUE INDEX `performers_name_unique` on `performers` (`name`) WHERE `disambiguation` IS NULL;

PRAGMA foreign_keys=ON;