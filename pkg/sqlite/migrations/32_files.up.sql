-- folders may be deleted independently. Don't cascade
CREATE TABLE `folders` (
  `id` integer not null primary key autoincrement,
  `path` varchar(255) NOT NULL,
  `parent_folder_id` integer,
  `mod_time` datetime not null,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`parent_folder_id`) references `folders`(`id`) on delete SET NULL
);

CREATE INDEX `index_folders_on_parent_folder_id` on `folders` (`parent_folder_id`);

-- require reference folders/zip files to be deleted manually first
CREATE TABLE `files` (
  `id` integer not null primary key autoincrement,
  `basename` varchar(255) NOT NULL,
  `zip_file_id` integer,
  `parent_folder_id` integer not null,
  `size` integer NOT NULL,
  `mod_time` datetime not null,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`parent_folder_id`) references `folders`(`id`),
  foreign key(`zip_file_id`) references `files`(`id`),
  CHECK (`basename` != '')
);

CREATE UNIQUE INDEX `index_files_zip_basename_unique` ON `files` (`zip_file_id`, `parent_folder_id`, `basename`);
CREATE INDEX `index_files_on_parent_folder_id_basename` on `files` (`parent_folder_id`, `basename`);
CREATE INDEX `index_files_on_basename` on `files` (`basename`);

ALTER TABLE `folders` ADD COLUMN `zip_file_id` integer REFERENCES `files`(`id`);
CREATE UNIQUE INDEX `index_folders_path_unique` on `folders` (`zip_file_id`, `path`);

CREATE TABLE `files_fingerprints` (
  `file_id` integer NOT NULL,
  `type` varchar(255) NOT NULL,
  `fingerprint` blob NOT NULL,
  foreign key(`file_id`) references `files`(`id`) on delete CASCADE,
  PRIMARY KEY (`file_id`, `type`, `fingerprint`)
);

CREATE INDEX `index_fingerprint_type_fingerprint` ON `files_fingerprints` (`type`, `fingerprint`);

CREATE TABLE `video_files` (
  `file_id` integer NOT NULL primary key,
  `duration` float NOT NULL,
	`video_codec` varchar(255) NOT NULL,
	`format` varchar(255) NOT NULL,
	`audio_codec` varchar(255) NOT NULL,
	`width` tinyint NOT NULL,
	`height` tinyint NOT NULL,
	`frame_rate` float NOT NULL,
	`bit_rate` integer NOT NULL,
  `interactive` boolean not null default '0',
  `interactive_speed` int,
  foreign key(`file_id`) references `files`(`id`) on delete CASCADE
);

CREATE TABLE `video_captions` (
  `file_id` integer NOT NULL,
  `language_code` varchar(255) NOT NULL,
  `filename` varchar(255) NOT NULL,
  `caption_type` varchar(255) NOT NULL,
  primary key (`file_id`, `language_code`, `caption_type`),
  foreign key(`file_id`) references `video_files`(`file_id`) on delete CASCADE
);

CREATE TABLE `image_files` (
  `file_id` integer NOT NULL primary key,
  `format` varchar(255) NOT NULL,
  `width` tinyint NOT NULL,
	`height` tinyint NOT NULL,
  foreign key(`file_id`) references `files`(`id`) on delete CASCADE
);

CREATE TABLE `images_files` (
    `image_id` integer NOT NULL,
    `file_id` integer NOT NULL,
    `primary` boolean NOT NULL,
    foreign key(`image_id`) references `images`(`id`) on delete CASCADE,
    foreign key(`file_id`) references `files`(`id`) on delete CASCADE,
    PRIMARY KEY(`image_id`, `file_id`)
);

CREATE INDEX `index_images_files_file_id` ON `images_files` (`file_id`);

CREATE TABLE `galleries_files` (
    `gallery_id` integer NOT NULL,
    `file_id` integer NOT NULL,
    `primary` boolean NOT NULL,
    foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE,
    foreign key(`file_id`) references `files`(`id`) on delete CASCADE,
    PRIMARY KEY(`gallery_id`, `file_id`)
);

CREATE INDEX `index_galleries_files_file_id` ON `galleries_files` (`file_id`);

CREATE TABLE `scenes_files` (
    `scene_id` integer NOT NULL,
    `file_id` integer NOT NULL,
    `primary` boolean NOT NULL,
    foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE,
    foreign key(`file_id`) references `files`(`id`) on delete CASCADE,
    PRIMARY KEY(`scene_id`, `file_id`)
);

CREATE INDEX `index_scenes_files_file_id` ON `scenes_files` (`file_id`);

PRAGMA foreign_keys=OFF;

CREATE TABLE `images_new` (
  `id` integer not null primary key autoincrement,
  -- REMOVED: `path` varchar(510) not null,
  -- REMOVED: `checksum` varchar(255) not null,
  `title` varchar(255),
  `rating` tinyint,
  -- REMOVED: `size` integer,
  -- REMOVED: `width` tinyint,
  -- REMOVED: `height` tinyint,
  `studio_id` integer,
  `o_counter` tinyint not null default 0,
  `organized` boolean not null default '0',
  -- REMOVED: `file_mod_time` datetime,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL
);

INSERT INTO `images_new`
  (
    `id`,
    `title`,
    `rating`,
    `studio_id`,
    `o_counter`,
    `organized`,
    `created_at`,
    `updated_at`
  )
  SELECT 
    `id`,
    `title`,
    `rating`,
    `studio_id`,
    `o_counter`,
    `organized`,
    `created_at`,
    `updated_at`
  FROM `images`;

-- create temporary placeholder folder
INSERT INTO `folders` (`path`, `mod_time`, `created_at`, `updated_at`) VALUES ('', '1970-01-01 00:00:00', '1970-01-01 00:00:00', '1970-01-01 00:00:00');

-- insert image files - we will fix these up in the post-migration
INSERT INTO `files`
  (
    `basename`,
    `parent_folder_id`,
    `size`,
    `mod_time`,
    `created_at`,
    `updated_at`
  )
  SELECT
    `path`,
    1,
    COALESCE(`size`, 0),
    -- set mod time to epoch so that it the format/size is calculated on scan
    '1970-01-01 00:00:00',
    `created_at`,
    `updated_at`
  FROM `images`;

INSERT INTO `image_files`
  (
    `file_id`,
    `format`,
    `width`,
    `height`
  )
  SELECT
    `files`.`id`,
    '',
    COALESCE(`images`.`width`, 0),
    COALESCE(`images`.`height`, 0)
  FROM `images` INNER JOIN `files` ON `images`.`path` = `files`.`basename` AND `files`.`parent_folder_id` = 1;

INSERT INTO `images_files`
  (
    `image_id`,
    `file_id`,
    `primary`
  )
  SELECT
    `images`.`id`,
    `files`.`id`,
    1
  FROM `images` INNER JOIN `files` ON `images`.`path` = `files`.`basename` AND `files`.`parent_folder_id` = 1;

INSERT INTO `files_fingerprints`
  (
    `file_id`,
    `type`,
    `fingerprint`
  )
  SELECT
    `files`.`id`,
    'md5',
    `images`.`checksum`
  FROM `images` INNER JOIN `files` ON `images`.`path` = `files`.`basename` AND `files`.`parent_folder_id` = 1;

DROP TABLE `images`;
ALTER TABLE `images_new` rename to `images`;

CREATE INDEX `index_images_on_studio_id` on `images` (`studio_id`);


CREATE TABLE `galleries_new` (
  `id` integer not null primary key autoincrement,
  -- REMOVED: `path` varchar(510),
  -- REMOVED: `checksum` varchar(255) not null,
  -- REMOVED: `zip` boolean not null default '0',
  `folder_id` integer,
  `title` varchar(255),
  `url` varchar(255),
  `date` date,
  `details` text,
  `studio_id` integer,
  `rating` tinyint,
  -- REMOVED: `file_mod_time` datetime,
  `organized` boolean not null default '0',
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL,
  foreign key(`folder_id`) references `folders`(`id`) on delete SET NULL
);

INSERT INTO `galleries_new`
  (
    `id`,
    `title`,
    `url`,
    `date`,
    `details`,
    `studio_id`,
    `rating`,
    `organized`,
    `created_at`,
    `updated_at`
  )
  SELECT 
    `id`,
    `title`,
    `url`,
    `date`,
    `details`,
    `studio_id`,
    `rating`,
    `organized`,
    `created_at`,
    `updated_at`
  FROM `galleries`;

-- insert gallery files - we will fix these up in the post-migration
INSERT INTO `files`
  (
    `basename`,
    `parent_folder_id`,
    `size`,
    `mod_time`,
    `created_at`,
    `updated_at`
  )
  SELECT
    `path`,
    1,
    0,
    '1970-01-01 00:00:00', -- set to placeholder so that size is updated
    `created_at`,
    `updated_at`
  FROM `galleries`
  WHERE `galleries`.`path` IS NOT NULL AND `galleries`.`zip` = '1';

-- insert gallery zip folders - we will fix these up in the post-migration
INSERT INTO `folders`
  (
    `path`,
    `zip_file_id`,
    `mod_time`,
    `created_at`,
    `updated_at`
  )
  SELECT
    `galleries`.`path`,
    `files`.`id`,
    '1970-01-01 00:00:00',
    `galleries`.`created_at`,
    `galleries`.`updated_at`
  FROM `galleries` 
  INNER JOIN `files` ON `galleries`.`path` = `files`.`basename` AND `files`.`parent_folder_id` = 1
  WHERE `galleries`.`path` IS NOT NULL AND `galleries`.`zip` = '1';

-- set the zip file id of the zip folders
UPDATE `folders` SET `zip_file_id` = (SELECT `files`.`id` FROM `files` WHERE `folders`.`path` = `files`.`basename`); 

-- insert gallery folders - we will fix these up in the post-migration
INSERT INTO `folders`
  (
    `path`,
    `mod_time`,
    `created_at`,
    `updated_at`
  )
  SELECT
    `path`,
    '1970-01-01 00:00:00',
    `created_at`,
    `updated_at`
  FROM `galleries`
  WHERE `galleries`.`path` IS NOT NULL AND `galleries`.`zip` = '0';

UPDATE `galleries_new` SET `folder_id` = (
  SELECT `folders`.`id` FROM `folders` INNER JOIN `galleries` ON `galleries_new`.`id` = `galleries`.`id` WHERE `folders`.`path` = `galleries`.`path` AND `galleries`.`zip` = '0'
);

INSERT INTO `galleries_files`
  (
    `gallery_id`,
    `file_id`,
    `primary`
  )
  SELECT
    `galleries`.`id`,
    `files`.`id`,
    1
  FROM `galleries` INNER JOIN `files` ON `galleries`.`path` = `files`.`basename` AND `files`.`parent_folder_id` = 1;

INSERT INTO `files_fingerprints`
  (
    `file_id`,
    `type`,
    `fingerprint`
  )
  SELECT
    `files`.`id`,
    'md5',
    `galleries`.`checksum`
  FROM `galleries` INNER JOIN `files` ON `galleries`.`path` = `files`.`basename` AND `files`.`parent_folder_id` = 1;

DROP TABLE `galleries`;
ALTER TABLE `galleries_new` rename to `galleries`;

CREATE INDEX `index_galleries_on_studio_id` on `galleries` (`studio_id`);
-- should only be possible to create a single gallery per folder
CREATE UNIQUE INDEX `index_galleries_on_folder_id_unique` on `galleries` (`folder_id`);

CREATE TABLE `scenes_new` (
  `id` integer not null primary key autoincrement,
  -- REMOVED: `path` varchar(510) not null,
  -- REMOVED: `checksum` varchar(255),
  -- REMOVED: `oshash` varchar(255),
  `title` varchar(255),
  `details` text,
  `url` varchar(255),
  `date` date,
  `rating` tinyint,
  -- REMOVED: `size` varchar(255),
  -- REMOVED: `duration` float,
  -- REMOVED: `video_codec` varchar(255),
  -- REMOVED: `audio_codec` varchar(255),
  -- REMOVED: `width` tinyint,
  -- REMOVED: `height` tinyint,
  -- REMOVED: `framerate` float,
  -- REMOVED: `bitrate` integer,
  `studio_id` integer,
  `o_counter` tinyint not null default 0,
  -- REMOVED: `format` varchar(255),
  `organized` boolean not null default '0',
  -- REMOVED: `interactive` boolean not null default '0',
  -- REMOVED: `interactive_speed` int,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  -- REMOVED: `file_mod_time` datetime,
  -- REMOVED: `phash` blob,
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL
  -- REMOVED: CHECK (`checksum` is not null or `oshash` is not null)
);

INSERT INTO `scenes_new`
  (
    `id`,
    `title`,
    `details`,
    `url`,
    `date`,
    `rating`,
    `studio_id`,
    `o_counter`,
    `organized`,
    `created_at`,
    `updated_at`
  )
  SELECT 
    `id`,
    `title`,
    `details`,
    `url`,
    `date`,
    `rating`,
    `studio_id`,
    `o_counter`,
    `organized`,
    `created_at`,
    `updated_at`
  FROM `scenes`;

-- insert scene files - we will fix these up in the post-migration
INSERT INTO `files`
  (
    `basename`,
    `parent_folder_id`,
    `size`,
    `mod_time`,
    `created_at`,
    `updated_at`
  )
  SELECT
    `path`,
    1,
    COALESCE(`size`, 0),
    -- set mod time to epoch so that it the format/size is calculated on scan
    '1970-01-01 00:00:00',
    `created_at`,
    `updated_at`
  FROM `scenes`;

INSERT INTO `video_files`
  (
    `file_id`,
    `duration`,
    `video_codec`,
    `format`,
    `audio_codec`,
    `width`,
    `height`,
    `frame_rate`,
    `bit_rate`,
    `interactive`,
    `interactive_speed`
  )
  SELECT
    `files`.`id`,
    `scenes`.`duration`,
    COALESCE(`scenes`.`video_codec`, ''),
    COALESCE(`scenes`.`format`, ''),
    COALESCE(`scenes`.`audio_codec`, ''),
    COALESCE(`scenes`.`width`, 0),
    COALESCE(`scenes`.`height`, 0),
    COALESCE(`scenes`.`framerate`, 0),
    COALESCE(`scenes`.`bitrate`, 0),
    `scenes`.`interactive`,
    `scenes`.`interactive_speed`
  FROM `scenes` INNER JOIN `files` ON `scenes`.`path` = `files`.`basename` AND `files`.`parent_folder_id` = 1;

INSERT INTO `scenes_files`
  (
    `scene_id`,
    `file_id`,
    `primary`
  )
  SELECT
    `scenes`.`id`,
    `files`.`id`,
    1
  FROM `scenes` INNER JOIN `files` ON `scenes`.`path` = `files`.`basename` AND `files`.`parent_folder_id` = 1;

INSERT INTO `files_fingerprints`
  (
    `file_id`,
    `type`,
    `fingerprint`
  )
  SELECT
    `files`.`id`,
    'md5',
    `scenes`.`checksum`
  FROM `scenes` INNER JOIN `files` ON `scenes`.`path` = `files`.`basename` AND `files`.`parent_folder_id` = 1
  WHERE `scenes`.`checksum` is not null;

INSERT INTO `files_fingerprints`
  (
    `file_id`,
    `type`,
    `fingerprint`
  )
  SELECT
    `files`.`id`,
    'oshash',
    `scenes`.`oshash`
  FROM `scenes` INNER JOIN `files` ON `scenes`.`path` = `files`.`basename` AND `files`.`parent_folder_id` = 1
  WHERE `scenes`.`oshash` is not null;

INSERT INTO `files_fingerprints`
  (
    `file_id`,
    `type`,
    `fingerprint`
  )
  SELECT
    `files`.`id`,
    'phash',
    `scenes`.`phash`
  FROM `scenes` INNER JOIN `files` ON `scenes`.`path` = `files`.`basename` AND `files`.`parent_folder_id` = 1
  WHERE `scenes`.`phash` is not null;

INSERT INTO `video_captions`
  (
    `file_id`,
    `language_code`,
    `filename`,
    `caption_type`
  )
  SELECT
    `files`.`id`,
    `scene_captions`.`language_code`,
    `scene_captions`.`filename`,
    `scene_captions`.`caption_type`
  FROM `scene_captions` 
  INNER JOIN `scenes` ON `scene_captions`.`scene_id` = `scenes`.`id`
  INNER JOIN `files` ON `scenes`.`path` = `files`.`basename` AND `files`.`parent_folder_id` = 1;

DROP TABLE `scenes`;
DROP TABLE `scene_captions`;

ALTER TABLE `scenes_new` rename to `scenes`;
CREATE INDEX `index_scenes_on_studio_id` on `scenes` (`studio_id`);

PRAGMA foreign_keys=ON;

-- create views to simplify queries

CREATE VIEW `images_query` AS 
  SELECT 
    `images`.`id`,
    `images`.`title`,
    `images`.`rating`,
    `images`.`organized`,
    `images`.`o_counter`,
    `images`.`studio_id`,
    `images`.`created_at`,
    `images`.`updated_at`,
    `galleries_images`.`gallery_id`,
    `images_tags`.`tag_id`,
    `performers_images`.`performer_id`,
    `image_files`.`format` as `image_format`,
    `image_files`.`width` as `image_width`,
    `image_files`.`height` as `image_height`,
    `files`.`id` as `file_id`,
    `files`.`basename`,
    `files`.`size`,
    `files`.`mod_time`,
    `files`.`zip_file_id`,
    `folders`.`id` as `parent_folder_id`,
    `folders`.`path` as `parent_folder_path`,
    `zip_files`.`basename` as `zip_basename`,
    `zip_files_folders`.`path` as `zip_folder_path`,
    `files_fingerprints`.`type` as `fingerprint_type`,
    `files_fingerprints`.`fingerprint`
  FROM `images`
  LEFT JOIN `performers_images` ON (`images`.`id` = `performers_images`.`image_id`) 
  LEFT JOIN `galleries_images` ON (`images`.`id` = `galleries_images`.`image_id`) 
  LEFT JOIN `images_tags` ON (`images`.`id` = `images_tags`.`image_id`)
  LEFT JOIN `images_files` ON (`images`.`id` = `images_files`.`image_id`) 
  LEFT JOIN `image_files` ON (`images_files`.`file_id` = `image_files`.`file_id`) 
  LEFT JOIN `files` ON (`images_files`.`file_id` = `files`.`id`) 
  LEFT JOIN `folders` ON (`files`.`parent_folder_id` = `folders`.`id`) 
  LEFT JOIN `files` AS `zip_files` ON (`files`.`zip_file_id` = `zip_files`.`id`)
  LEFT JOIN `folders` AS `zip_files_folders` ON (`zip_files`.`parent_folder_id` = `zip_files_folders`.`id`)
  LEFT JOIN `files_fingerprints` ON (`images_files`.`file_id` = `files_fingerprints`.`file_id`);

CREATE VIEW `galleries_query` AS 
  SELECT 
    `galleries`.`id`,
    `galleries`.`title`,
    `galleries`.`url`,
    `galleries`.`date`,
    `galleries`.`details`,
    `galleries`.`rating`,
    `galleries`.`organized`,
    `galleries`.`studio_id`,
    `galleries`.`created_at`,
    `galleries`.`updated_at`,
    `galleries_tags`.`tag_id`,
    `scenes_galleries`.`scene_id`,
    `performers_galleries`.`performer_id`,
    `galleries_folders`.`id` as `folder_id`,
    `galleries_folders`.`path` as `folder_path`,
    `files`.`id` as `file_id`,
    `files`.`basename`,
    `files`.`size`,
    `files`.`mod_time`,
    `files`.`zip_file_id`,
    `parent_folders`.`id` as `parent_folder_id`,
    `parent_folders`.`path` as `parent_folder_path`,
    `zip_files`.`basename` as `zip_basename`,
    `zip_files_folders`.`path` as `zip_folder_path`,
    `files_fingerprints`.`type` as `fingerprint_type`,
    `files_fingerprints`.`fingerprint`
  FROM `galleries`
  LEFT JOIN `performers_galleries` ON (`galleries`.`id` = `performers_galleries`.`gallery_id`) 
  LEFT JOIN `galleries_tags` ON (`galleries`.`id` = `galleries_tags`.`gallery_id`)
  LEFT JOIN `scenes_galleries` ON (`galleries`.`id` = `scenes_galleries`.`gallery_id`) 
  LEFT JOIN `folders` AS `galleries_folders` ON (`galleries`.`folder_id` = `galleries_folders`.`id`) 
  LEFT JOIN `galleries_files` ON (`galleries`.`id` = `galleries_files`.`gallery_id`) 
  LEFT JOIN `files` ON (`galleries_files`.`file_id` = `files`.`id`) 
  LEFT JOIN `folders` AS `parent_folders` ON (`files`.`parent_folder_id` = `parent_folders`.`id`) 
  LEFT JOIN `files` AS `zip_files` ON (`files`.`zip_file_id` = `zip_files`.`id`)
  LEFT JOIN `folders` AS `zip_files_folders` ON (`zip_files`.`parent_folder_id` = `zip_files_folders`.`id`)
  LEFT JOIN `files_fingerprints` ON (`galleries_files`.`file_id` = `files_fingerprints`.`file_id`);

CREATE VIEW `scenes_query` AS 
  SELECT 
    `scenes`.`id`,
    `scenes`.`title`,
    `scenes`.`details`,
    `scenes`.`url`,
    `scenes`.`date`,
    `scenes`.`rating`,
    `scenes`.`studio_id`,
    `scenes`.`o_counter`,
    `scenes`.`organized`,
    `scenes`.`created_at`,
    `scenes`.`updated_at`,
    `scenes_tags`.`tag_id`,
    `scenes_galleries`.`gallery_id`,
    `performers_scenes`.`performer_id`,
    `movies_scenes`.`movie_id`,
    `movies_scenes`.`scene_index`,
    `scene_stash_ids`.`stash_id`,
    `scene_stash_ids`.`endpoint`,
    `video_files`.`format` as `video_format`,
    `video_files`.`width` as `video_width`,
    `video_files`.`height` as `video_height`,
    `video_files`.`duration`,
    `video_files`.`video_codec`,
    `video_files`.`audio_codec`,
    `video_files`.`frame_rate`,
    `video_files`.`bit_rate`,
    `video_files`.`interactive`,
    `video_files`.`interactive_speed`,
    `files`.`id` as `file_id`,
    `files`.`basename`,
    `files`.`size`,
    `files`.`mod_time`,
    `files`.`zip_file_id`,
    `folders`.`id` as `parent_folder_id`,
    `folders`.`path` as `parent_folder_path`,
    `zip_files`.`basename` as `zip_basename`,
    `zip_files_folders`.`path` as `zip_folder_path`,
    `files_fingerprints`.`type` as `fingerprint_type`,
    `files_fingerprints`.`fingerprint`
  FROM `scenes`
  LEFT JOIN `performers_scenes` ON (`scenes`.`id` = `performers_scenes`.`scene_id`) 
  LEFT JOIN `scenes_tags` ON (`scenes`.`id` = `scenes_tags`.`scene_id`)
  LEFT JOIN `movies_scenes` ON (`scenes`.`id` = `movies_scenes`.`scene_id`)
  LEFT JOIN `scene_stash_ids` ON (`scenes`.`id` = `scene_stash_ids`.`scene_id`)
  LEFT JOIN `scenes_galleries` ON (`scenes`.`id` = `scenes_galleries`.`scene_id`) 
  LEFT JOIN `scenes_files` ON (`scenes`.`id` = `scenes_files`.`scene_id`) 
  LEFT JOIN `video_files` ON (`scenes_files`.`file_id` = `video_files`.`file_id`) 
  LEFT JOIN `files` ON (`scenes_files`.`file_id` = `files`.`id`) 
  LEFT JOIN `folders` ON (`files`.`parent_folder_id` = `folders`.`id`) 
  LEFT JOIN `files` AS `zip_files` ON (`files`.`zip_file_id` = `zip_files`.`id`)
  LEFT JOIN `folders` AS `zip_files_folders` ON (`zip_files`.`parent_folder_id` = `zip_files_folders`.`id`)
  LEFT JOIN `files_fingerprints` ON (`scenes_files`.`file_id` = `files_fingerprints`.`file_id`);
