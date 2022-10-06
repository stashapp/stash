
-- need to change scenes.checksum to be nullable
ALTER TABLE `scenes` rename to `_scenes_old`;

CREATE TABLE `scenes` (
  `id` integer not null primary key autoincrement,
  `path` varchar(510) not null,
  -- nullable
  `checksum` varchar(255),
  -- add oshash
  `oshash` varchar(255),
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
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL,
  -- add check to ensure at least one hash is set
  CHECK (`checksum` is not null or `oshash` is not null)
);

DROP INDEX IF EXISTS `scenes_path_unique`;
DROP INDEX IF EXISTS `scenes_checksum_unique`;
DROP INDEX IF EXISTS `index_scenes_on_studio_id`;

CREATE UNIQUE INDEX `scenes_path_unique` on `scenes` (`path`);
CREATE UNIQUE INDEX `scenes_checksum_unique` on `scenes` (`checksum`);
CREATE UNIQUE INDEX `scenes_oshash_unique` on `scenes` (`oshash`);
CREATE INDEX `index_scenes_on_studio_id` on `scenes` (`studio_id`);

-- recreate the tables referencing scenes to correct their references
ALTER TABLE `galleries` rename to `_galleries_old`;
ALTER TABLE `performers_scenes` rename to `_performers_scenes_old`;
ALTER TABLE `scene_markers` rename to `_scene_markers_old`;
ALTER TABLE `scene_markers_tags` rename to `_scene_markers_tags_old`;
ALTER TABLE `scenes_tags` rename to `_scenes_tags_old`;
ALTER TABLE `movies_scenes` rename to `_movies_scenes_old`;
ALTER TABLE `scenes_cover` rename to `_scenes_cover_old`;

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

CREATE TABLE `scenes_cover` (
  `scene_id` integer,
  `cover` blob not null,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE
);

DROP INDEX `index_scene_covers_on_scene_id`;

CREATE UNIQUE INDEX `index_scene_covers_on_scene_id` on `scenes_cover` (`scene_id`);

-- now populate from the old tables
-- these tables are changed so require the full column def
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

-- these tables are a direct copy
INSERT INTO `galleries` SELECT * from `_galleries_old`;
INSERT INTO `performers_scenes` SELECT * from `_performers_scenes_old`;
INSERT INTO `scene_markers` SELECT * from `_scene_markers_old`;
INSERT INTO `scene_markers_tags` SELECT * from `_scene_markers_tags_old`;
INSERT INTO `scenes_tags` SELECT * from `_scenes_tags_old`;
INSERT INTO `movies_scenes` SELECT * from `_movies_scenes_old`;
INSERT INTO `scenes_cover` SELECT * from `_scenes_cover_old`;

-- drop old tables
DROP TABLE `_scenes_old`;
DROP TABLE `_galleries_old`;
DROP TABLE `_performers_scenes_old`;
DROP TABLE `_scene_markers_old`;
DROP TABLE `_scene_markers_tags_old`;
DROP TABLE `_scenes_tags_old`;
DROP TABLE `_movies_scenes_old`;
DROP TABLE `_scenes_cover_old`;
