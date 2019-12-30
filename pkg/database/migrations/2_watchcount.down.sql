--

PRAGMA foreign_keys=OFF;

ALTER TABLE `scenes` RENAME TO `scenes_old`;
DROP INDEX IF EXISTS `scenes_path_unique`;
DROP INDEX IF EXISTS `scenes_checksum_unique`;
DROP INDEX IF EXISTS `index_scenes_on_studio_id`;

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
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`studio_id`) references `studios`(`id`) on delete CASCADE
);
CREATE UNIQUE INDEX `scenes_path_unique` on `scenes` (`path`);
CREATE UNIQUE INDEX `scenes_checksum_unique` on `scenes` (`checksum`);
CREATE INDEX `index_scenes_on_studio_id` on `scenes` (`studio_id`);

INSERT INTO `scenes` (id, path, checksum, title, details, url, date,
rating, size, duration, video_codec, audio_codec, width, height, framerate,
bitrate, studio_id, created_at, updated_at)
SELECT id, path, checksum, title, details, url, date,
rating, size, duration, video_codec, audio_codec, width, height, framerate,
bitrate, studio_id, created_at, updated_at
FROM `scenes_old`;

DROP TABLE `scenes_old`;

PRAGMA foreign_keys=ON;