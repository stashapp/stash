CREATE TABLE `scene_stash_ids` (
  `scene_id` integer,
  `endpoint` varchar(255),
  `stash_id` varchar(36),
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE
);

CREATE TABLE `performer_stash_ids` (
  `performer_id` integer,
  `endpoint` varchar(255),
  `stash_id` varchar(36),
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE
);

CREATE TABLE `studio_stash_ids` (
  `studio_id` integer,
  `endpoint` varchar(255),
  `stash_id` varchar(36),
  foreign key(`studio_id`) references `studios`(`id`) on delete CASCADE
);
