CREATE TABLE `scene_captions` (
  `scene_id` integer,
  `language_code` varchar(255) NOT NULL,
  `filename` varchar(255) NOT NULL,
  `caption_type` varchar(255) NOT NULL,
  primary key (`scene_id`, `language_code`, `caption_type`),
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE
);
