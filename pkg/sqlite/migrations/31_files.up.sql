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

CREATE UNIQUE INDEX `index_files_basename_unique` ON `files` (`zip_file_id`, `parent_folder_id`, `basename`);

ALTER TABLE `folders` ADD COLUMN `zip_file_id` integer REFERENCES `files`(`id`) ON DELETE CASCADE;
CREATE UNIQUE INDEX `index_folders_path_unique` on `folders` (`zip_file_id`, `path`);

CREATE TABLE `fingerprints` (
  `id` integer not null primary key autoincrement,
  `type` varchar(255) NOT NULL,
  `fingerprint` blob NOT NULL
);

CREATE UNIQUE INDEX `index_fingerprint_type_fingerprint_unique` ON `fingerprints` (`type`, `fingerprint`);

CREATE TABLE `files_fingerprints` (
  `file_id` integer NOT NULL,
  `fingerprint_id` integer NOT NULL,
  foreign key(`file_id`) references `files`(`id`) on delete CASCADE,
  foreign key(`fingerprint_id`) references `fingerprints`(`id`) on delete CASCADE,
  PRIMARY KEY (`file_id`, `fingerprint_id`)
);

CREATE INDEX `index_files_fingerprints_fingerprint_id` ON `files_fingerprints` (`fingerprint_id`);

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
  foreign key(`file_id`) references `files`(`id`) on delete CASCADE
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

CREATE TABLE `images_fingerprints` (
  `image_id` integer NOT NULL,
  `fingerprint_id` integer NOT NULL,
  foreign key(`image_id`) references `images`(`id`) on delete CASCADE,
  foreign key(`fingerprint_id`) references `fingerprints`(`id`) on delete CASCADE,
  PRIMARY KEY (`image_id`, `fingerprint_id`)
);

CREATE INDEX `index_images_fingerprints_fingerprint_id` ON `images_fingerprints` (`fingerprint_id`);

CREATE TABLE `galleries_files` (
    `gallery_id` integer NOT NULL,
    `file_id` integer NOT NULL,
    `primary` boolean NOT NULL,
    foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE,
    foreign key(`file_id`) references `files`(`id`) on delete CASCADE,
    PRIMARY KEY(`gallery_id`, `file_id`)
);

CREATE INDEX `index_galleries_files_file_id` ON `galleries_files` (`file_id`);

CREATE TABLE `galleries_fingerprints` (
  `gallery_id` integer NOT NULL,
  `fingerprint_id` integer NOT NULL,
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE,
  foreign key(`fingerprint_id`) references `fingerprints`(`id`) on delete CASCADE,
  PRIMARY KEY (`gallery_id`, `fingerprint_id`)
);

CREATE INDEX `index_galleries_fingerprints_fingerprint_id` ON `galleries_fingerprints` (`fingerprint_id`);

CREATE TABLE `scenes_files` (
    `scene_id` integer NOT NULL,
    `file_id` integer NOT NULL,
    `primary` boolean NOT NULL,
    foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE,
    foreign key(`file_id`) references `files`(`id`) on delete CASCADE,
    PRIMARY KEY(`scene_id`, `file_id`)
);

CREATE INDEX `index_scenes_files_file_id` ON `scenes_files` (`file_id`);

CREATE TABLE `scenes_fingerprints` (
  `scene_id` integer NOT NULL,
  `fingerprint_id` integer NOT NULL,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE,
  foreign key(`fingerprint_id`) references `fingerprints`(`id`) on delete CASCADE,
  PRIMARY KEY (`scene_id`, `fingerprint_id`)
);

CREATE INDEX `index_scenes_fingerprints_fingerprint_id` ON `scenes_fingerprints` (`fingerprint_id`);