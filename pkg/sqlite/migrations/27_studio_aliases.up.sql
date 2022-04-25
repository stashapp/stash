CREATE TABLE `studio_aliases` (
  `studio_id` integer,
  `alias` varchar(255) NOT NULL,
  foreign key(`studio_id`) references `studios`(`id`) on delete CASCADE
);

CREATE UNIQUE INDEX `studio_aliases_alias_unique` on `studio_aliases` (`alias`);
