PRAGMA foreign_keys=OFF;

CREATE TABLE `performer_aliases` (
  `performer_id` integer NOT NULL,
  `alias` varchar(255) NOT NULL,
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE,
  PRIMARY KEY(`performer_id`, `alias`)
);

CREATE INDEX `performer_aliases_alias` on `performer_aliases` (`alias`);

DROP INDEX `performers_checksum_unique`;

-- drop aliases and checksum
-- add disambiguation

CREATE TABLE `performers_new` (
  `id` integer not null primary key autoincrement,
  `name` varchar(255),
  `disambiguation` varchar(255),
  `gender` varchar(20),
  `url` varchar(255),
  `twitter` varchar(255),
  `instagram` varchar(255),
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
  `ignore_auto_tag` boolean not null default '0'
);

INSERT INTO `performers_new`
  (
    `id`,
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
    `favorite`,
    `created_at`,
    `updated_at`,
    `details`,
    `death_date`,
    `hair_color`,
    `weight`,
    `rating`,
    `ignore_auto_tag`
  )
  SELECT 
    `id`,
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
    `favorite`,
    `created_at`,
    `updated_at`,
    `details`,
    `death_date`,
    `hair_color`,
    `weight`,
    `rating`,
    `ignore_auto_tag`
  FROM `performers`;

INSERT INTO `performer_aliases`
  (
    `performer_id`,
    `alias`
  )
  SELECT 
    `id`,
    `aliases`
  FROM `performers`
  WHERE `performers`.`aliases` IS NOT NULL AND `performers`.`aliases` != '';

DROP TABLE `performers`;
ALTER TABLE `performers_new` rename to `performers`;


-- these will be executed in the post-migration
-- CREATE UNIQUE INDEX `performers_name_disambiguation_unique` on `performers` (`name`, `disambiguation`) WHERE `disambiguation` IS NOT NULL;
-- CREATE UNIQUE INDEX `performers_name_unique` on `performers` (`name`) WHERE `disambiguation` IS NULL;

PRAGMA foreign_keys=ON;
