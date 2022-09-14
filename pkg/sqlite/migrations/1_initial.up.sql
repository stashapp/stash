CREATE TABLE `tags` (
  `id` integer not null primary key autoincrement,
  `name` varchar(255),
  `created_at` datetime not null,
  `updated_at` datetime not null
);
CREATE TABLE `studios` (
  `id` integer not null primary key autoincrement,
  `image` blob not null,
  `checksum` varchar(255) not null,
  `name` varchar(255),
  `url` varchar(255),
  `created_at` datetime not null,
  `updated_at` datetime not null
);
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
CREATE TABLE `scenes_tags` (
  `scene_id` integer,
  `tag_id` integer,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`)
);
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
CREATE TABLE `scene_markers_tags` (
  `scene_marker_id` integer,
  `tag_id` integer,
  foreign key(`scene_marker_id`) references `scene_markers`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`)
);
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
CREATE TABLE `performers_scenes` (
  `performer_id` integer,
  `scene_id` integer,
  foreign key(`performer_id`) references `performers`(`id`),
  foreign key(`scene_id`) references `scenes`(`id`)
);
CREATE TABLE `performers` (
  `id` integer not null primary key autoincrement,
  `image` blob not null,
  `checksum` varchar(255) not null,
  `name` varchar(255),
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
CREATE TABLE `galleries` (
  `id` integer not null primary key autoincrement,
  `path` varchar(510) not null,
  `checksum` varchar(255) not null,
  `scene_id` integer,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`scene_id`) references `scenes`(`id`)
);
CREATE UNIQUE INDEX `studios_checksum_unique` on `studios` (`checksum`);
CREATE UNIQUE INDEX `scenes_path_unique` on `scenes` (`path`);
CREATE UNIQUE INDEX `scenes_checksum_unique` on `scenes` (`checksum`);
CREATE UNIQUE INDEX `performers_checksum_unique` on `performers` (`checksum`);
CREATE INDEX `index_tags_on_name` on `tags` (`name`);
CREATE INDEX `index_studios_on_name` on `studios` (`name`);
CREATE INDEX `index_studios_on_checksum` on `studios` (`checksum`);
CREATE INDEX `index_scraped_items_on_studio_id` on `scraped_items` (`studio_id`);
CREATE INDEX `index_scenes_tags_on_tag_id` on `scenes_tags` (`tag_id`);
CREATE INDEX `index_scenes_tags_on_scene_id` on `scenes_tags` (`scene_id`);
CREATE INDEX `index_scenes_on_studio_id` on `scenes` (`studio_id`);
CREATE INDEX `index_scene_markers_tags_on_tag_id` on `scene_markers_tags` (`tag_id`);
CREATE INDEX `index_scene_markers_tags_on_scene_marker_id` on `scene_markers_tags` (`scene_marker_id`);
CREATE INDEX `index_scene_markers_on_scene_id` on `scene_markers` (`scene_id`);
CREATE INDEX `index_scene_markers_on_primary_tag_id` on `scene_markers` (`primary_tag_id`);
CREATE INDEX `index_performers_scenes_on_scene_id` on `performers_scenes` (`scene_id`);
CREATE INDEX `index_performers_scenes_on_performer_id` on `performers_scenes` (`performer_id`);
CREATE INDEX `index_performers_on_name` on `performers` (`name`);
CREATE INDEX `index_performers_on_checksum` on `performers` (`checksum`);
CREATE INDEX `index_galleries_on_scene_id` on `galleries` (`scene_id`);
CREATE UNIQUE INDEX `galleries_path_unique` on `galleries` (`path`);
CREATE UNIQUE INDEX `galleries_checksum_unique` on `galleries` (`checksum`);
