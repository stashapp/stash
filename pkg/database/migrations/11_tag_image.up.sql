CREATE TABLE `tags_image` (
  `tag_id` integer,
  `image` blob not null,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE
);

CREATE UNIQUE INDEX `index_tag_image_on_tag_id` on `tags_image` (`tag_id`);
