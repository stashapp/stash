CREATE TABLE `scene_segments` (
  `id` integer not null primary key autoincrement,
  `scene_id` integer not null,
  `title` varchar(255) not null,
  `start_seconds` float not null,
  `end_seconds` float not null,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE
);

CREATE INDEX `scene_segments_scene_id` on `scene_segments` (`scene_id`);
