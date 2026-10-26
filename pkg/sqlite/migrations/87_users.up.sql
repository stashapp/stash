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

PRAGMA foreign_keys=ON;
