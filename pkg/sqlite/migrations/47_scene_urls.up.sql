PRAGMA foreign_keys=OFF;

CREATE TABLE `scene_urls` (
  `scene_id` integer NOT NULL,
  `position` integer NOT NULL,
  `url` varchar(255) NOT NULL,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE,
  PRIMARY KEY(`scene_id`, `position`, `url`)
);

CREATE INDEX `scene_urls_url` on `scene_urls` (`url`);

-- drop url
CREATE TABLE "scenes_new" (
  `id` integer not null primary key autoincrement,
  `title` varchar(255),
  `details` text,
  `date` date,
  `rating` tinyint,
  `studio_id` integer,
  `o_counter` tinyint not null default 0,
  `organized` boolean not null default '0',
  `created_at` datetime not null,
  `updated_at` datetime not null, 
  `code` text, 
  `director` text, 
  `resume_time` float not null default 0, 
  `last_played_at` datetime default null, 
  `play_count` tinyint not null default 0, 
  `play_duration` float not null default 0, 
  `cover_blob` varchar(255) REFERENCES `blobs`(`checksum`),
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL
);

INSERT INTO `scenes_new`
  (
    `id`,
    `title`,
    `details`,
    `date`,
    `rating`,
    `studio_id`,
    `o_counter`,
    `organized`,
    `created_at`,
    `updated_at`,
    `code`,
    `director`,
    `resume_time`,
    `last_played_at`,
    `play_count`,
    `play_duration`,
    `cover_blob`
  )
  SELECT 
    `id`,
    `title`,
    `details`,
    `date`,
    `rating`,
    `studio_id`,
    `o_counter`,
    `organized`,
    `created_at`,
    `updated_at`,
    `code`,
    `director`,
    `resume_time`,
    `last_played_at`,
    `play_count`,
    `play_duration`,
    `cover_blob`
  FROM `scenes`;

INSERT INTO `scene_urls`
  (
    `scene_id`,
    `position`,
    `url`
  )
  SELECT 
    `id`,
    '0',
    `url`
  FROM `scenes`
  WHERE `scenes`.`url` IS NOT NULL AND `scenes`.`url` != '';

DROP INDEX `index_scenes_on_studio_id`;
DROP TABLE `scenes`;
ALTER TABLE `scenes_new` rename to `scenes`;

CREATE INDEX `index_scenes_on_studio_id` on `scenes` (`studio_id`);

PRAGMA foreign_keys=ON;
