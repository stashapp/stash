CREATE TABLE `performers_tags` (
  `performer_id` integer NOT NULL,
  `tag_id` integer NOT NULL,
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE
);

CREATE INDEX `index_performers_tags_on_tag_id` on `performers_tags` (`tag_id`);
CREATE INDEX `index_performers_tags_on_performer_id` on `performers_tags` (`performer_id`);
