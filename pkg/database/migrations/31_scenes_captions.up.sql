CREATE TABLE `scene_captions` (
  `scene_id` integer,
  `language_code` varchar(255) NOT NULL,
  `caption_type` varchar(255) NOT NULL,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE
);