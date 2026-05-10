ALTER TABLE `tags` ADD COLUMN `supports_numeric_rating` boolean not null default '0';

CREATE TABLE `scene_tag_ratings` (
  `scene_id` integer NOT NULL,
  `tag_id` integer NOT NULL,
  `rating` tinyint NOT NULL,
  PRIMARY KEY (`scene_id`, `tag_id`),
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE
);

CREATE INDEX `index_scene_tag_ratings_tag_id` ON `scene_tag_ratings` (`tag_id`);
