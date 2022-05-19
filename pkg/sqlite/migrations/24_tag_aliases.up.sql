CREATE TABLE `tag_aliases` (
  `tag_id` integer,
  `alias` varchar(255) NOT NULL,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE
);

CREATE UNIQUE INDEX `tag_aliases_alias_unique` on `tag_aliases` (`alias`);
