CREATE TABLE `galleries_chapters` (
  `id` integer not null primary key autoincrement,
  `title` varchar(255) not null,
  `image_index` integer not null,
  `gallery_id` integer not null,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE
);
CREATE INDEX `index_galleries_chapters_on_gallery_id` on `galleries_chapters` (`gallery_id`);
