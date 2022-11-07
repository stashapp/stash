-- add primary keys to association tables that are missing them
PRAGMA foreign_keys=OFF;

CREATE TABLE `performers_new` (
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
  -- changed from varchar(255)
  `height` int,
  `measurements` varchar(255),
  `fake_tits` varchar(255),
  `career_length` varchar(255),
  `tattoos` varchar(255),
  `piercings` varchar(255),
  `aliases` varchar(255),
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
    `details`,
    `death_date`,
    `hair_color`,
    `weight`,
    `rating`,
    `ignore_auto_tag`
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
    CASE `height`
      WHEN '' THEN NULL
      WHEN NULL THEN NULL
      ELSE CAST(`height` as int)
    END,
    `measurements`,
    `fake_tits`,
    `career_length`,
    `tattoos`,
    `piercings`,
    `aliases`,
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

DROP TABLE `performers`;
ALTER TABLE `performers_new` rename to `performers`;

CREATE UNIQUE INDEX `performers_checksum_unique` on `performers` (`checksum`);
CREATE INDEX `index_performers_on_name` on `performers` (`name`);
