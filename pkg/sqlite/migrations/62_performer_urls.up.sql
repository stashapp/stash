PRAGMA foreign_keys=OFF;

CREATE TABLE `performer_urls` (
  `performer_id` integer NOT NULL,
  `position` integer NOT NULL,
  `url` varchar(255) NOT NULL,
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE,
  PRIMARY KEY(`performer_id`, `position`, `url`)
);

CREATE INDEX `performers_urls_url` on `performer_urls` (`url`);

-- drop url, twitter and instagram
-- make name not null
CREATE TABLE `performers_new` (
  `id` integer not null primary key autoincrement,
  `name` varchar(255) not null,
  `disambiguation` varchar(255),
  `gender` varchar(20),
  `birthdate` date,
  `ethnicity` varchar(255),
  `country` varchar(255),
  `eye_color` varchar(255),
  `height` int,
  `measurements` varchar(255),
  `fake_tits` varchar(255),
  `career_length` varchar(255),
  `tattoos` varchar(255),
  `piercings` varchar(255),
  `favorite` boolean not null default '0',
  `created_at` datetime not null,
  `updated_at` datetime not null,
  `details` text, 
  `death_date` date, 
  `hair_color` varchar(255), 
  `weight` integer, 
  `rating` tinyint, 
  `ignore_auto_tag` boolean not null default '0', 
  `image_blob` varchar(255) REFERENCES `blobs`(`checksum`), 
  `penis_length` float, 
  `circumcised` varchar[10]
);

INSERT INTO `performers_new`
  (
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
    `career_length`,
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
    `circumcised`
  )
  SELECT 
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
    `career_length`,
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
    `circumcised`
  FROM `performers`;

INSERT INTO `performer_urls`
  (
    `performer_id`,
    `position`,
    `url`
  )
  SELECT 
    `id`,
    '0',
    `url`
  FROM `performers`
  WHERE `performers`.`url` IS NOT NULL AND `performers`.`url` != '';

INSERT INTO `performer_urls`
  (
    `performer_id`,
    `position`,
    `url`
  )
  SELECT 
    `id`,
    (SELECT count(*) FROM `performer_urls` WHERE `performer_id` = `performers`.`id`)+1,
    CASE
      WHEN `twitter` LIKE 'http%://%' THEN `twitter`
      ELSE 'https://www.twitter.com/' || `twitter`
    END
  FROM `performers`
  WHERE `performers`.`twitter` IS NOT NULL AND `performers`.`twitter` != '';

INSERT INTO `performer_urls`
  (
    `performer_id`,
    `position`,
    `url`
  )
  SELECT 
    `id`,
    (SELECT count(*) FROM `performer_urls` WHERE `performer_id` = `performers`.`id`)+1,
    CASE
      WHEN `instagram` LIKE 'http%://%' THEN `instagram`
      ELSE 'https://www.instagram.com/' || `instagram`
    END
  FROM `performers`
  WHERE `performers`.`instagram` IS NOT NULL AND `performers`.`instagram` != '';

DROP INDEX IF EXISTS `performers_name_disambiguation_unique`;
DROP INDEX IF EXISTS `performers_name_unique`;
DROP TABLE IF EXISTS `performers`;
ALTER TABLE `performers_new` rename to `performers`;

CREATE UNIQUE INDEX `performers_name_disambiguation_unique` on `performers` (`name`, `disambiguation`) WHERE `disambiguation` IS NOT NULL;
CREATE UNIQUE INDEX `performers_name_unique` on `performers` (`name`) WHERE `disambiguation` IS NULL;

PRAGMA foreign_keys=ON;
