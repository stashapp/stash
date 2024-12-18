PRAGMA foreign_keys=OFF;

CREATE TABLE `movie_urls` (
  `movie_id` integer NOT NULL,
  `position` integer NOT NULL,
  `url` varchar(255) NOT NULL,
  foreign key(`movie_id`) references `movies`(`id`) on delete CASCADE,
  PRIMARY KEY(`movie_id`, `position`, `url`)
);

CREATE INDEX `movie_urls_url` on `movie_urls` (`url`);

-- drop url
CREATE TABLE `movies_new` (
  `id` integer not null primary key autoincrement,
  `name` varchar(255) not null,
  `aliases` varchar(255),
  `duration` integer,
  `date` date,
  `rating` tinyint,
  `studio_id` integer REFERENCES `studios`(`id`) ON DELETE SET NULL,
  `director` varchar(255),
  `synopsis` text,
  `created_at` datetime not null,
  `updated_at` datetime not null, 
  `front_image_blob` varchar(255) REFERENCES `blobs`(`checksum`), 
  `back_image_blob` varchar(255) REFERENCES `blobs`(`checksum`)
);

INSERT INTO `movies_new`
  (
    `id`,
    `name`,
    `aliases`,
    `duration`,
    `date`,
    `rating`,
    `studio_id`,
    `director`,
    `synopsis`,
    `created_at`,
    `updated_at`,
    `front_image_blob`,
    `back_image_blob`
  )
  SELECT 
    `id`,
    `name`,
    `aliases`,
    `duration`,
    `date`,
    `rating`,
    `studio_id`,
    `director`,
    `synopsis`,
    `created_at`,
    `updated_at`,
    `front_image_blob`,
    `back_image_blob`
  FROM `movies`;

INSERT INTO `movie_urls`
  (
    `movie_id`,
    `position`,
    `url`
  )
  SELECT 
    `id`,
    '0',
    `url`
  FROM `movies`
  WHERE `movies`.`url` IS NOT NULL AND `movies`.`url` != '';

DROP INDEX `index_movies_on_name_unique`;
DROP INDEX `index_movies_on_studio_id`;
DROP TABLE `movies`;
ALTER TABLE `movies_new` rename to `movies`;

CREATE INDEX `index_movies_on_name` ON `movies`(`name`);
CREATE INDEX `index_movies_on_studio_id` on `movies` (`studio_id`);

PRAGMA foreign_keys=ON;
