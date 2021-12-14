CREATE TABLE `files` (
  `id` integer not null primary key autoincrement,
  `path` varchar(255) NOT NULL,
  `zip_file_id` integer,
  `checksum` varchar(255),
  `oshash` varchar(255),
  `size` integer NOT NULL,
  `duration` float,
  `video_codec` varchar(255),
  `audio_codec` varchar(255),
  `width` tinyint,
  `height` tinyint,
  `framerate` float,
  `bitrate` integer,
  `format` varchar(255),
  `mod_time` datetime,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`zip_file_id`) references `files`(`id`) on delete CASCADE,
  CHECK (`checksum` is not null or `oshash` is not null)
);

CREATE UNIQUE INDEX `index_file_path_unique` ON `files` (`path`);
CREATE INDEX `file_checksum` on `files` (`checksum`);
CREATE INDEX `file_oshash` on `files` (`oshash`);

CREATE TABLE `scenes_files` (
  `scene_id` integer,
  `file_id` integer,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE,
  foreign key(`file_id`) references `files`(`id`) on delete CASCADE
);

CREATE UNIQUE INDEX `index_scenes_files_unique` ON `scenes_files` ( `scene_id`, `file_id` );
CREATE INDEX `index_scenes_files_on_file_id` on `scenes_files` (`file_id`);
CREATE INDEX `index_scenes_files_on_scene_id` on `scenes_files` (`scene_id`);

CREATE TABLE `images_files` (
  `image_id` integer,
  `file_id` integer,
  foreign key(`image_id`) references `images`(`id`) on delete CASCADE,
  foreign key(`file_id`) references `files`(`id`) on delete CASCADE
);

CREATE UNIQUE INDEX `index_images_files_unique` ON `images_files` ( `image_id`, `file_id` );
CREATE INDEX `index_images_files_on_file_id` on `images_files` (`file_id`);
CREATE INDEX `index_images_files_on_image_id` on `images_files` (`image_id`);

CREATE TABLE `galleries_files` (
  `gallery_id` integer,
  `file_id` integer,
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE,
  foreign key(`file_id`) references `files`(`id`) on delete CASCADE
);

CREATE UNIQUE INDEX `index_galleries_files_unique` ON `galleries_files` ( `gallery_id`, `file_id` );
CREATE INDEX `index_galleries_files_on_file_id` on `galleries_files` (`file_id`);
CREATE INDEX `index_galleries_files_on_gallery_id` on `galleries_files` (`gallery_id`);

-- translate scenes, images and galleries into files
INSERT INTO `files`
 (
  `path`,
  `checksum`,
  `oshash`,
  `size`,
  `duration`,
  `video_codec`,
  `audio_codec`,
  `width`,
  `height`,
  `framerate`,
  `bitrate`,
  `format`,
  `mod_time`,
  `created_at`,
  `updated_at`
 )
 SELECT
  `path`,
  `checksum`,
  `oshash`,
  COALESCE(`size`, 0),
  `duration`,
  `video_codec`,
  `audio_codec`,
  `width`,
  `height`,
  `framerate`,
  `bitrate`,
  `format`,
  `file_mod_time`,
  `created_at`,
  `updated_at`
 FROM `scenes`;

INSERT INTO `scenes_files`
 (
  `scene_id`,
  `file_id`
 )
 SELECT
  `scenes`.`id`,
  `files`.`id`
 FROM `scenes`
 INNER JOIN `files` ON (`scenes`.`checksum` = `files`.`checksum` OR `scenes`.`oshash` = `files`.`oshash`);

INSERT INTO `files`
 (
  `path`,
  `checksum`,
  `size`,
  `width`,
  `height`,
  `mod_time`,
  `created_at`,
  `updated_at`
 )
 SELECT
  `path`,
  `checksum`,
  COALESCE(`size`, 0),
  `width`,
  `height`,
  `file_mod_time`,
  `created_at`,
  `updated_at`
 FROM `images`;

INSERT INTO `images_files`
 (
  `image_id`,
  `file_id`
 )
 SELECT
  `images`.`id`,
  `files`.`id`
 FROM `images`
 INNER JOIN `files` ON (`images`.`checksum` = `files`.`checksum`);

INSERT INTO `files`
 (
  `path`,
  `size`,
  `checksum`,
  `mod_time`,
  `created_at`,
  `updated_at`
 )
 SELECT
  `path`,
  0,
  `checksum`,
  `file_mod_time`,
  `created_at`,
  `updated_at`
 FROM `galleries` WHERE `zip` = '1';

INSERT INTO `galleries_files`
 (
  `gallery_id`,
  `file_id`
 )
 SELECT
  `galleries`.`id`,
  `files`.`id`
 FROM `galleries`
 INNER JOIN `files` ON (`galleries`.`checksum` = `files`.`checksum`);
 