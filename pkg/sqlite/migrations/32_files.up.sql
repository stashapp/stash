CREATE TABLE `folders` (
  `id` integer not null primary key autoincrement,
  `path` varchar(255) NOT NULL,
  `parent_folder_id` integer,
  `mod_time` datetime not null,
  `missing_since` datetime,
  `last_scanned` datetime not null,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`parent_folder_id`) references `folders`(`id`) on delete CASCADE
);

CREATE INDEX `index_folders_on_parent_folder_id` on `folders` (`parent_folder_id`);

CREATE TABLE `files` (
  `id` integer not null primary key autoincrement,
  `basename` varchar(255) NOT NULL,
  `zip_file_id` integer,
  `parent_folder_id` integer not null,
  `size` integer NOT NULL,
  `mod_time` datetime not null,
  `missing_since` datetime,
  `last_scanned` datetime not null,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`parent_folder_id`) references `folders`(`id`) on delete CASCADE,
  foreign key(`zip_file_id`) references `files`(`id`) on delete CASCADE,
  CHECK (`basename` != '')
);

CREATE UNIQUE INDEX `index_files_zip_basename_unique` ON `files` (`zip_file_id`, `parent_folder_id`, `basename`);
CREATE INDEX `index_files_on_parent_folder_id_basename` on `files` (`parent_folder_id`, `basename`);
CREATE INDEX `index_files_on_basename` on `files` (`basename`);

ALTER TABLE `folders` ADD COLUMN `zip_file_id` integer REFERENCES `files`(`id`) ON DELETE CASCADE;
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

DROP TABLE `images`;
ALTER TABLE `images_new` rename to `images`;

CREATE INDEX `index_images_on_studio_id` on `images` (`studio_id`);


CREATE TABLE `galleries_new` (
  `id` integer not null primary key autoincrement,
  -- REMOVED: `path` varchar(510),
  -- REMOVED: `checksum` varchar(255) not null,
  -- REMOVED: `zip` boolean not null default '0',
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
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL
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

DROP TABLE `galleries`;
ALTER TABLE `galleries_new` rename to `galleries`;

CREATE INDEX `index_galleries_on_studio_id` on `galleries` (`studio_id`);

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

-- TODO - transfer fingerprint information
-- TODO - transfer captions

DROP TABLE `scenes`;

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
    `files`.`missing_since`,
    `files`.`last_scanned`,
    `files`.`zip_file_id`,
    `folders`.`id` as `parent_folder_id`,
    `folders`.`path` as `folder_path`,
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
    `files`.`id` as `file_id`,
    `files`.`basename`,
    `files`.`size`,
    `files`.`mod_time`,
    `files`.`missing_since`,
    `files`.`last_scanned`,
    `files`.`zip_file_id`,
    `folders`.`id` as `parent_folder_id`,
    `folders`.`path` as `folder_path`,
    `zip_files`.`basename` as `zip_basename`,
    `zip_files_folders`.`path` as `zip_folder_path`,
    `files_fingerprints`.`type` as `fingerprint_type`,
    `files_fingerprints`.`fingerprint`
  FROM `galleries`
  LEFT JOIN `performers_galleries` ON (`galleries`.`id` = `performers_galleries`.`gallery_id`) 
  LEFT JOIN `galleries_tags` ON (`galleries`.`id` = `galleries_tags`.`gallery_id`)
  LEFT JOIN `scenes_galleries` ON (`galleries`.`id` = `scenes_galleries`.`gallery_id`) 
  LEFT JOIN `galleries_files` ON (`galleries`.`id` = `galleries_files`.`gallery_id`) 
  LEFT JOIN `files` ON (`galleries_files`.`file_id` = `files`.`id`) 
  LEFT JOIN `folders` ON (`files`.`parent_folder_id` = `folders`.`id`) 
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
    `files`.`missing_since`,
    `files`.`last_scanned`,
    `files`.`zip_file_id`,
    `folders`.`id` as `parent_folder_id`,
    `folders`.`path` as `folder_path`,
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
