CREATE TABLE `movies_tags` (
  `movie_id` integer NOT NULL,
  `tag_id` integer NOT NULL,
  foreign key(`movie_id`) references `movies`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE,
  PRIMARY KEY(`movie_id`, `tag_id`)
);

CREATE INDEX `index_movies_tags_on_tag_id` on `movies_tags` (`tag_id`);
CREATE INDEX `index_movies_tags_on_movie_id` on `movies_tags` (`movie_id`);
