CREATE TABLE `scene_captions` (
  `scene_id` integer,
  `caption` varchar(255) NOT NULL,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE
);

CREATE UNIQUE INDEX `scene_captions_caption_unique` on `scene_captions` (`caption`);
