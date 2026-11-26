CREATE TABLE `scene_locked_fields` (
  `scene_id` integer NOT NULL,
  `field` varchar(255) NOT NULL,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE,
  PRIMARY KEY(`scene_id`, `field`)
);
