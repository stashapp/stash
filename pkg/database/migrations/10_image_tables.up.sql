-- recreate scenes, studios and performers tables
ALTER TABLE `studios` rename to `_studios_old`;
ALTER TABLE `scenes` rename to `_scenes_old`;
ALTER TABLE `performers` RENAME TO `_performers_old`;
ALTER TABLE `movies` rename to `_movies_old`;

-- remove studio image
CREATE TABLE `studios` (
  `id` integer not null primary key autoincrement,
  `checksum` varchar(255) not null,
  `name` varchar(255),
  `url` varchar(255),
  `parent_id` integer DEFAULT NULL CHECK ( id IS NOT parent_id ) REFERENCES studios(id) on delete set null,
  `created_at` datetime not null,
  `updated_at` datetime not null
);

DROP INDEX `studios_checksum_unique`;
DROP INDEX `index_studios_on_name`;
DROP INDEX `index_studios_on_checksum`;

CREATE UNIQUE INDEX `studios_checksum_unique` on `studios` (`checksum`);
CREATE INDEX `index_studios_on_name` on `studios` (`name`);
CREATE INDEX `index_studios_on_checksum` on `studios` (`checksum`);

-- remove scene cover
CREATE TABLE `scenes` (
  `id` integer not null primary key autoincrement,
  `path` varchar(510) not null,
  `checksum` varchar(255) not null,
  `title` varchar(255),
  `details` text,
  `url` varchar(255),
  `date` date,
  `rating` tinyint,
  `size` varchar(255),
  `duration` float,
  `video_codec` varchar(255),
  `audio_codec` varchar(255),
  `width` tinyint,
  `height` tinyint,
  `framerate` float,
  `bitrate` integer,
  `studio_id` integer,
  `o_counter` tinyint not null default 0,
  `format` varchar(255),
  `created_at` datetime not null,
  `updated_at` datetime not null,
  -- changed from cascade delete
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL
);

DROP INDEX IF EXISTS `scenes_path_unique`;
DROP INDEX IF EXISTS `scenes_checksum_unique`;
DROP INDEX IF EXISTS `index_scenes_on_studio_id`;

CREATE UNIQUE INDEX `scenes_path_unique` on `scenes` (`path`);
CREATE UNIQUE INDEX `scenes_checksum_unique` on `scenes` (`checksum`);
CREATE INDEX `index_scenes_on_studio_id` on `scenes` (`studio_id`);

-- remove performer image
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
  `updated_at` datetime not null
);

DROP INDEX `performers_checksum_unique`;
DROP INDEX `index_performers_on_name`;

CREATE UNIQUE INDEX `performers_checksum_unique` on `performers` (`checksum`);
CREATE INDEX `index_performers_on_name` on `performers` (`name`);

-- remove front_image and back_image
CREATE TABLE `movies` (
  `id` integer not null primary key autoincrement,
  `name` varchar(255) not null,
  `aliases` varchar(255),
  `duration` integer,
  `date` date,
  `rating` tinyint,
  `studio_id` integer,
  `director` varchar(255),
  `synopsis` text,
  `checksum` varchar(255) not null,
  `url` varchar(255),
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`studio_id`) references `studios`(`id`) on delete set null
);

DROP INDEX `movies_name_unique`;
DROP INDEX `movies_checksum_unique`;
DROP INDEX `index_movies_on_studio_id`;

CREATE UNIQUE INDEX `movies_name_unique` on `movies` (`name`);
CREATE UNIQUE INDEX `movies_checksum_unique` on `movies` (`checksum`);
CREATE INDEX `index_movies_on_studio_id` on `movies` (`studio_id`);

-- recreate the tables referencing the above tables to correct their references
ALTER TABLE `galleries` rename to `_galleries_old`;
ALTER TABLE `performers_scenes` rename to `_performers_scenes_old`;
ALTER TABLE `scene_markers` rename to `_scene_markers_old`;
ALTER TABLE `scene_markers_tags` rename to `_scene_markers_tags_old`;
ALTER TABLE `scenes_tags` rename to `_scenes_tags_old`;
ALTER TABLE `movies_scenes` rename to `_movies_scenes_old`;
ALTER TABLE `scraped_items` rename to `_scraped_items_old`;

CREATE TABLE `galleries` (
  `id` integer not null primary key autoincrement,
  `path` varchar(510) not null,
  `checksum` varchar(255) not null,
  `scene_id` integer,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`scene_id`) references `scenes`(`id`)
);

DROP INDEX IF EXISTS `index_galleries_on_scene_id`;
DROP INDEX IF EXISTS `galleries_path_unique`;
DROP INDEX IF EXISTS `galleries_checksum_unique`;

CREATE INDEX `index_galleries_on_scene_id` on `galleries` (`scene_id`);
CREATE UNIQUE INDEX `galleries_path_unique` on `galleries` (`path`);
CREATE UNIQUE INDEX `galleries_checksum_unique` on `galleries` (`checksum`);

CREATE TABLE `performers_scenes` (
  `performer_id` integer,
  `scene_id` integer,
  foreign key(`performer_id`) references `performers`(`id`),
  foreign key(`scene_id`) references `scenes`(`id`)
);

DROP INDEX `index_performers_scenes_on_scene_id`;
DROP INDEX `index_performers_scenes_on_performer_id`;

CREATE INDEX `index_performers_scenes_on_scene_id` on `performers_scenes` (`scene_id`);
CREATE INDEX `index_performers_scenes_on_performer_id` on `performers_scenes` (`performer_id`);

CREATE TABLE `scene_markers` (
  `id` integer not null primary key autoincrement,
  `title` varchar(255) not null,
  `seconds` float not null,
  `primary_tag_id` integer not null,
  `scene_id` integer,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`primary_tag_id`) references `tags`(`id`),
  foreign key(`scene_id`) references `scenes`(`id`)
);

DROP INDEX `index_scene_markers_on_scene_id`;
DROP INDEX `index_scene_markers_on_primary_tag_id`;

CREATE INDEX `index_scene_markers_on_scene_id` on `scene_markers` (`scene_id`);
CREATE INDEX `index_scene_markers_on_primary_tag_id` on `scene_markers` (`primary_tag_id`);

CREATE TABLE `scene_markers_tags` (
  `scene_marker_id` integer,
  `tag_id` integer,
  foreign key(`scene_marker_id`) references `scene_markers`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`)
);

DROP INDEX `index_scene_markers_tags_on_tag_id`;
DROP INDEX `index_scene_markers_tags_on_scene_marker_id`;

CREATE INDEX `index_scene_markers_tags_on_tag_id` on `scene_markers_tags` (`tag_id`);
CREATE INDEX `index_scene_markers_tags_on_scene_marker_id` on `scene_markers_tags` (`scene_marker_id`);

CREATE TABLE `scenes_tags` (
  `scene_id` integer,
  `tag_id` integer,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`)
);

DROP INDEX `index_scenes_tags_on_tag_id`;
DROP INDEX `index_scenes_tags_on_scene_id`;

CREATE INDEX `index_scenes_tags_on_tag_id` on `scenes_tags` (`tag_id`);
CREATE INDEX `index_scenes_tags_on_scene_id` on `scenes_tags` (`scene_id`);

CREATE TABLE `movies_scenes` (
  `movie_id` integer,
  `scene_id` integer,
  `scene_index` tinyint,
  foreign key(`movie_id`) references `movies`(`id`) on delete cascade,
  foreign key(`scene_id`) references `scenes`(`id`) on delete cascade
);

DROP INDEX `index_movies_scenes_on_movie_id`;
DROP INDEX `index_movies_scenes_on_scene_id`;

CREATE INDEX `index_movies_scenes_on_movie_id` on `movies_scenes` (`movie_id`);
CREATE INDEX `index_movies_scenes_on_scene_id` on `movies_scenes` (`scene_id`);

-- remove movie_id since doesn't appear to be used
CREATE TABLE `scraped_items` (
  `id` integer not null primary key autoincrement,
  `title` varchar(255),
  `description` text,
  `url` varchar(255),
  `date` date,
  `rating` varchar(255),
  `tags` varchar(510),
  `models` varchar(510),
  `episode` integer,
  `gallery_filename` varchar(255),
  `gallery_url` varchar(510),
  `video_filename` varchar(255),
  `video_url` varchar(255),
  `studio_id` integer,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`studio_id`) references `studios`(`id`)
);

DROP INDEX `index_scraped_items_on_studio_id`;

CREATE INDEX `index_scraped_items_on_studio_id` on `scraped_items` (`studio_id`);

-- now populate from the old tables
-- these tables are changed so require the full column def
INSERT INTO `studios` 
  (
    `id`,
    `checksum`,
    `name`,
    `url`,
    `parent_id`,
    `created_at`,
    `updated_at`
  )
  SELECT 
    `id`,
    `checksum`,
    `name`,
    `url`,
    `parent_id`,
    `created_at`,
    `updated_at`
  FROM `_studios_old`;

INSERT INTO `scenes`
  (
    `id`,
    `path`,
    `checksum`,
    `title`,
    `details`,
    `url`,
    `date`,
    `rating`,
    `size`,
    `duration`,
    `video_codec`,
    `audio_codec`,
    `width`,
    `height`,
    `framerate`,
    `bitrate`,
    `studio_id`,
    `o_counter`,
    `format`,
    `created_at`,
    `updated_at`
  )
  SELECT 
    `id`,
    `path`,
    `checksum`,
    `title`,
    `details`,
    `url`,
    `date`,
    `rating`,
    `size`,
    `duration`,
    `video_codec`,
    `audio_codec`,
    `width`,
    `height`,
    `framerate`,
    `bitrate`,
    `studio_id`,
    `o_counter`,
    `format`,
    `created_at`,
    `updated_at`
  FROM `_scenes_old`;

INSERT INTO `performers` 
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
    `updated_at`
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
    `updated_at`
  FROM `_performers_old`;

INSERT INTO `movies`
  (
    `id`,
    `name`,
    `aliases`,
    `duration`,
    `date`,
    `rating`,
    `studio_id`,
    `director`,
    `synopsis`,
    `checksum`,
    `url`,
    `created_at`,
    `updated_at`
  )
  SELECT
    `id`,
    `name`,
    `aliases`,
    `duration`,
    `date`,
    `rating`,
    `studio_id`,
    `director`,
    `synopsis`,
    `checksum`,
    `url`,
    `created_at`,
    `updated_at`
  FROM `_movies_old`;

INSERT INTO `scraped_items`
  (
    `id`,
    `title`,
    `description`,
    `url`,
    `date`,
    `rating`,
    `tags`,
    `models`,
    `episode`,
    `gallery_filename`,
    `gallery_url`,
    `video_filename`,
    `video_url`,
    `studio_id`,
    `created_at`,
    `updated_at`
  )
  SELECT
    `id`,
    `title`,
    `description`,
    `url`,
    `date`,
    `rating`,
    `tags`,
    `models`,
    `episode`,
    `gallery_filename`,
    `gallery_url`,
    `video_filename`,
    `video_url`,
    `studio_id`,
    `created_at`,
    `updated_at`
  FROM `_scraped_items_old`;

-- these tables are a direct copy
INSERT INTO `galleries` SELECT * from `_galleries_old`;
INSERT INTO `performers_scenes` SELECT * from `_performers_scenes_old`;
INSERT INTO `scene_markers` SELECT * from `_scene_markers_old`;
INSERT INTO `scene_markers_tags` SELECT * from `_scene_markers_tags_old`;
INSERT INTO `scenes_tags` SELECT * from `_scenes_tags_old`;
INSERT INTO `movies_scenes` SELECT * from `_movies_scenes_old`;

-- populate covers in separate table
CREATE TABLE `scenes_cover` (
  `scene_id` integer,
  `cover` blob not null,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE
);

CREATE UNIQUE INDEX `index_scene_covers_on_scene_id` on `scenes_cover` (`scene_id`);

INSERT INTO `scenes_cover` 
  (
    `scene_id`,
    `cover`
  )
  SELECT `id`, `cover` from `_scenes_old` where `cover` is not null;

-- put performer images in separate table
CREATE TABLE `performers_image` (
  `performer_id` integer,
  `image` blob not null,
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE
);

CREATE UNIQUE INDEX `index_performer_image_on_performer_id` on `performers_image` (`performer_id`);

INSERT INTO `performers_image` 
  (
    `performer_id`,
    `image`
  )
  SELECT `id`, `image` from `_performers_old` where `image` is not null;

-- put studio images in separate table
CREATE TABLE `studios_image` (
  `studio_id` integer,
  `image` blob not null,
  foreign key(`studio_id`) references `studios`(`id`) on delete CASCADE
);

CREATE UNIQUE INDEX `index_studio_image_on_studio_id` on `studios_image` (`studio_id`);

INSERT INTO `studios_image` 
  (
    `studio_id`,
    `image`
  )
  SELECT `id`, `image` from `_studios_old` where `image` is not null;

-- put movie images in separate table
CREATE TABLE `movies_images` (
  `movie_id` integer,
  `front_image` blob not null,
  `back_image` blob,
  foreign key(`movie_id`) references `movies`(`id`) on delete CASCADE
);

CREATE UNIQUE INDEX `index_movie_images_on_movie_id` on `movies_images` (`movie_id`);

INSERT INTO `movies_images` 
  (
    `movie_id`,
    `front_image`,
    `back_image`
  )
  SELECT `id`, `front_image`, `back_image` from `_movies_old` where `front_image` is not null;

-- drop old tables
DROP TABLE `_scenes_old`;
DROP TABLE `_studios_old`;
DROP TABLE `_performers_old`;
DROP TABLE `_movies_old`;
DROP TABLE `_galleries_old`;
DROP TABLE `_performers_scenes_old`;
DROP TABLE `_scene_markers_old`;
DROP TABLE `_scene_markers_tags_old`;
DROP TABLE `_scenes_tags_old`;
DROP TABLE `_movies_scenes_old`;
DROP TABLE `_scraped_items_old`;
