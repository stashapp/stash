CREATE TABLE `users` (
    `id` INTEGER PRIMARY KEY,
    `username` TEXT NOT NULL,
    `notes` TEXT,
    `api_key` TEXT,
    `locked` BOOLEAN NOT NULL DEFAULT '0',
    `password_hash` TEXT,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL
);

CREATE UNIQUE INDEX `users_username_unique` on `users` (`username`);

CREATE TABLE `user_roles` (
    `user_id` INTEGER NOT NULL,
    `role` TEXT NOT NULL,
    PRIMARY KEY (`user_id`, `role`),
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);
-- existing user data is assigned to the initial user (id 1), which is created in the post-migration.
-- Foreign keys are disabled so that the data can be inserted before the user exists.
PRAGMA foreign_keys=OFF;

CREATE TABLE `scenes_user_data` (
  `scene_id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `rating` tinyint,
  `resume_time` float NOT NULL DEFAULT 0,
  `play_duration` float NOT NULL DEFAULT 0,
  PRIMARY KEY (`scene_id`, `user_id`),
  FOREIGN KEY (`scene_id`) REFERENCES `scenes`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE INDEX `index_scenes_user_data_on_user_id` ON `scenes_user_data` (`user_id`);

INSERT INTO `scenes_user_data` (`scene_id`, `user_id`, `rating`, `resume_time`, `play_duration`)
  SELECT `id`, 1, `rating`, `resume_time`, `play_duration` FROM `scenes`
  WHERE `rating` IS NOT NULL OR `resume_time` != 0 OR `play_duration` != 0;

ALTER TABLE `scenes` DROP COLUMN `rating`;
ALTER TABLE `scenes` DROP COLUMN `resume_time`;
ALTER TABLE `scenes` DROP COLUMN `play_duration`;

-- add user_id to scene o and view history
CREATE TABLE `scenes_o_dates_new` (
  `scene_id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `o_date` datetime NOT NULL,
  FOREIGN KEY (`scene_id`) REFERENCES `scenes`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

INSERT INTO `scenes_o_dates_new` (`scene_id`, `user_id`, `o_date`)
  SELECT `scene_id`, 1, `o_date` FROM `scenes_o_dates`;

DROP INDEX IF EXISTS `index_scenes_o_dates`;
DROP TABLE `scenes_o_dates`;
ALTER TABLE `scenes_o_dates_new` RENAME TO `scenes_o_dates`;
CREATE INDEX `index_scenes_o_dates` ON `scenes_o_dates` (`scene_id`, `user_id`);
CREATE INDEX `index_scenes_o_dates_on_user_id` ON `scenes_o_dates` (`user_id`);

CREATE TABLE `scenes_view_dates_new` (
  `scene_id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `view_date` datetime NOT NULL,
  FOREIGN KEY (`scene_id`) REFERENCES `scenes`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

INSERT INTO `scenes_view_dates_new` (`scene_id`, `user_id`, `view_date`)
  SELECT `scene_id`, 1, `view_date` FROM `scenes_view_dates`;

DROP INDEX IF EXISTS `index_scenes_view_dates`;
DROP TABLE `scenes_view_dates`;
ALTER TABLE `scenes_view_dates_new` RENAME TO `scenes_view_dates`;
CREATE INDEX `index_scenes_view_dates` ON `scenes_view_dates` (`scene_id`, `user_id`);
CREATE INDEX `index_scenes_view_dates_on_user_id` ON `scenes_view_dates` (`user_id`);

CREATE TABLE `images_user_data` (
  `image_id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `rating` tinyint,
  `o_counter` integer NOT NULL DEFAULT 0,
  PRIMARY KEY (`image_id`, `user_id`),
  FOREIGN KEY (`image_id`) REFERENCES `images`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE INDEX `index_images_user_data_on_user_id` ON `images_user_data` (`user_id`);

INSERT INTO `images_user_data` (`image_id`, `user_id`, `rating`, `o_counter`)
  SELECT `id`, 1, `rating`, `o_counter` FROM `images`
  WHERE `rating` IS NOT NULL OR `o_counter` != 0;

ALTER TABLE `images` DROP COLUMN `rating`;
ALTER TABLE `images` DROP COLUMN `o_counter`;

CREATE TABLE `galleries_user_data` (
  `gallery_id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `rating` tinyint,
  PRIMARY KEY (`gallery_id`, `user_id`),
  FOREIGN KEY (`gallery_id`) REFERENCES `galleries`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE INDEX `index_galleries_user_data_on_user_id` ON `galleries_user_data` (`user_id`);

INSERT INTO `galleries_user_data` (`gallery_id`, `user_id`, `rating`)
  SELECT `id`, 1, `rating` FROM `galleries`
  WHERE `rating` IS NOT NULL;

ALTER TABLE `galleries` DROP COLUMN `rating`;

CREATE TABLE `groups_user_data` (
  `group_id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `rating` tinyint,
  PRIMARY KEY (`group_id`, `user_id`),
  FOREIGN KEY (`group_id`) REFERENCES `groups`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE INDEX `index_groups_user_data_on_user_id` ON `groups_user_data` (`user_id`);

INSERT INTO `groups_user_data` (`group_id`, `user_id`, `rating`)
  SELECT `id`, 1, `rating` FROM `groups`
  WHERE `rating` IS NOT NULL;

ALTER TABLE `groups` DROP COLUMN `rating`;

CREATE TABLE `performers_user_data` (
  `performer_id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `rating` tinyint,
  `favorite` boolean NOT NULL DEFAULT false,
  PRIMARY KEY (`performer_id`, `user_id`),
  FOREIGN KEY (`performer_id`) REFERENCES `performers`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE INDEX `index_performers_user_data_on_user_id` ON `performers_user_data` (`user_id`);

INSERT INTO `performers_user_data` (`performer_id`, `user_id`, `rating`, `favorite`)
  SELECT `id`, 1, `rating`, `favorite` FROM `performers`
  WHERE `rating` IS NOT NULL OR `favorite` = 1;

ALTER TABLE `performers` DROP COLUMN `rating`;
ALTER TABLE `performers` DROP COLUMN `favorite`;

CREATE TABLE `studios_user_data` (
  `studio_id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `rating` tinyint,
  `favorite` boolean NOT NULL DEFAULT false,
  PRIMARY KEY (`studio_id`, `user_id`),
  FOREIGN KEY (`studio_id`) REFERENCES `studios`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE INDEX `index_studios_user_data_on_user_id` ON `studios_user_data` (`user_id`);

INSERT INTO `studios_user_data` (`studio_id`, `user_id`, `rating`, `favorite`)
  SELECT `id`, 1, `rating`, `favorite` FROM `studios`
  WHERE `rating` IS NOT NULL OR `favorite` = 1;

ALTER TABLE `studios` DROP COLUMN `rating`;
ALTER TABLE `studios` DROP COLUMN `favorite`;

CREATE TABLE `tags_user_data` (
  `tag_id` integer NOT NULL,
  `user_id` integer NOT NULL,
  `favorite` boolean NOT NULL DEFAULT false,
  PRIMARY KEY (`tag_id`, `user_id`),
  FOREIGN KEY (`tag_id`) REFERENCES `tags`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE INDEX `index_tags_user_data_on_user_id` ON `tags_user_data` (`user_id`);

INSERT INTO `tags_user_data` (`tag_id`, `user_id`, `favorite`)
  SELECT `id`, 1, `favorite` FROM `tags`
  WHERE `favorite` = 1;

ALTER TABLE `tags` DROP COLUMN `favorite`;

CREATE TABLE `saved_filters_new` (
  `id` integer not null primary key autoincrement,
  `user_id` integer not null,
  `name` varchar(510) not null,
  `mode` varchar(255) not null,
  `find_filter` blob,
  `object_filter` blob,
  `ui_options` blob,
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

INSERT INTO `saved_filters_new` (`id`, `user_id`, `name`, `mode`, `find_filter`, `object_filter`, `ui_options`)
  SELECT `id`, 1, `name`, `mode`, `find_filter`, `object_filter`, `ui_options` FROM `saved_filters`;

DROP INDEX `index_saved_filters_on_mode_name_unique`;
DROP TABLE `saved_filters`;
ALTER TABLE `saved_filters_new` RENAME TO `saved_filters`;

CREATE UNIQUE INDEX `index_saved_filters_on_user_id_mode_name_unique` ON `saved_filters` (`user_id`, `mode`, `name`);

-- users without a UI configuration use the UI configuration in the config file
CREATE TABLE `user_ui_config` (
  `user_id` integer NOT NULL PRIMARY KEY,
  `config` text NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

PRAGMA foreign_keys=ON;
