ALTER TABLE `movies` rename to `_movies_old`;
ALTER TABLE `movies_scenes` rename to `_movies_scenes_old`;

DROP INDEX IF EXISTS `movies_checksum_unique`;
DROP INDEX IF EXISTS `index_movie_id_scene_index_unique`;
DROP INDEX IF EXISTS `index_movies_scenes_on_movie_id`;
DROP INDEX IF EXISTS `index_movies_scenes_on_scene_id`;

-- recreate the movies table with fixed column types and constraints
CREATE TABLE `movies` (
  `id` integer not null primary key autoincrement,
  -- add not null
  `name` varchar(255) not null,
  `aliases` varchar(255),
  -- varchar(6) -> integer
  `duration` integer,
  `date` date,
  -- varchar(1) -> tinyint
  `rating` tinyint,
  `studio_id` integer,
  `director` varchar(255),
  `synopsis` text,
  `checksum` varchar(255) not null,
  `url` varchar(255),
  `created_at` datetime not null,
  `updated_at` datetime not null,
  `front_image` blob not null,
  `back_image` blob,
  foreign key(`studio_id`) references `studios`(`id`) on delete set null
);
CREATE TABLE `movies_scenes` (
  `movie_id` integer,
  `scene_id` integer,
  -- varchar(2) -> tinyint
  `scene_index` tinyint,
  foreign key(`movie_id`) references `movies`(`id`) on delete cascade,
  foreign key(`scene_id`) references `scenes`(`id`) on delete cascade
);

-- add unique index on movie name
CREATE UNIQUE INDEX `movies_name_unique` on `movies` (`name`);
CREATE UNIQUE INDEX `movies_checksum_unique` on `movies` (`checksum`);
-- remove unique index on movies_scenes
CREATE INDEX `index_movies_scenes_on_movie_id` on `movies_scenes` (`movie_id`);
CREATE INDEX `index_movies_scenes_on_scene_id` on `movies_scenes` (`scene_id`);
CREATE INDEX `index_movies_on_studio_id` on `movies` (`studio_id`);

-- custom functions cannot accept NULL values, so massage the old data
UPDATE `_movies_old` set `duration` = 0 WHERE `duration` IS NULL;

-- now populate from the old tables
INSERT INTO `movies`
  (
    `id`,
    `name`,
    `aliases`,
    `duration`,
    `date`,
    `rating`,
    `director`,
    `synopsis`,
    `front_image`,
    `back_image`,
    `checksum`,
    `url`,
    `created_at`,
    `updated_at`
  )
  SELECT 
    `id`,
    `name`,
    `aliases`,
    durationToTinyInt(`duration`),
    `date`,
    CAST(`rating` as tinyint),
    `director`,
    `synopsis`,
    `front_image`,
    `back_image`,
    `checksum`,
    `url`,
    `created_at`,
    `updated_at`
  FROM `_movies_old`
  -- ignore null named movies
  WHERE `name` is not null;

-- durationToTinyInt returns 0 if it cannot parse the string
-- set these values to null instead
UPDATE `movies` SET `duration` = NULL WHERE `duration` = 0;

INSERT INTO `movies_scenes` 
  (
    `movie_id`,
    `scene_id`,
    `scene_index`
  )
  SELECT
    `movie_id`,
    `scene_id`,
    CAST(`scene_index` as tinyint)
  FROM `_movies_scenes_old`;

-- drop old tables
DROP TABLE `_movies_scenes_old`;
DROP TABLE `_movies_old`;
