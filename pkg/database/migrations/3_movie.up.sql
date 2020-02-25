CREATE TABLE `movies` (
  `id` integer not null primary key autoincrement,
  `name` varchar(255),
  `aliases` varchar(255),
  `duration` varchar(6),
  `date` date,
  `rating` varchar(1),
  `director` varchar(255),
  `synopsis` text,
  `front_image` blob not null,
  `back_image` blob,
  `checksum` varchar(255) not null,
  `url` varchar(255),
  `created_at` datetime not null,
  `updated_at` datetime not null
);
CREATE TABLE `movies_scenes` (
  `movie_id` integer,
  `scene_id` integer,
  `scene_index` varchar(2),
  foreign key(`movie_id`) references `movies`(`id`),
  foreign key(`scene_id`) references `scenes`(`id`)
);


ALTER TABLE `scraped_items` ADD COLUMN `movie_id` integer;
CREATE UNIQUE INDEX `movies_checksum_unique` on `movies` (`checksum`);
CREATE UNIQUE INDEX `index_movie_id_scene_index_unique` ON `movies_scenes` ( `movie_id`, `scene_index` );
CREATE INDEX `index_movies_scenes_on_movie_id` on `movies_scenes` (`movie_id`);
CREATE INDEX `index_movies_scenes_on_scene_id` on `movies_scenes` (`scene_id`);


