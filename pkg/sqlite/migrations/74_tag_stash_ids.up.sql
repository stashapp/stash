CREATE TABLE `tag_stash_ids` (
  `tag_id` integer,
  `endpoint` varchar(255),
  `stash_id` varchar(36),
  `updated_at` datetime not null default '1970-01-01T00:00:00Z',
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE
);