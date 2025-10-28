CREATE TABLE `studios_tags` (
  `studio_id` integer NOT NULL,
  `tag_id` integer NOT NULL,
  foreign key(`studio_id`) references `studios`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE,
  PRIMARY KEY(`studio_id`, `tag_id`)
);

CREATE INDEX `index_studios_tags_on_tag_id` on `studios_tags` (`tag_id`);