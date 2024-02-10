PRAGMA foreign_keys=OFF;

CREATE TABLE `scenes_view_dates` (
  `scene_id` integer,
  `view_date` datetime not null,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE
);

CREATE TABLE `scenes_o_dates` (
  `scene_id` integer,
  `o_date` datetime not null,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE
);

-- drop o_counter, play_count and last_played_at
CREATE TABLE "scenes_new" (
  `id` integer not null primary key autoincrement,
  `title` varchar(255),
  `details` text,
  `date` date,
  `rating` tinyint,
  `studio_id` integer,
  `organized` boolean not null default '0',
  `created_at` datetime not null,
  `updated_at` datetime not null, 
  `code` text, 
  `director` text, 
  `resume_time` float not null default 0, 
  `play_duration` float not null default 0, 
  `cover_blob` varchar(255) REFERENCES `blobs`(`checksum`),
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL
);

INSERT INTO `scenes_new`
  (
    `id`,
    `title`,
    `details`,
    `date`,
    `rating`,
    `studio_id`,
    `organized`,
    `created_at`,
    `updated_at`,
    `code`,
    `director`,
    `resume_time`,
    `play_duration`,
    `cover_blob`
  )
  SELECT 
    `id`,
    `title`,
    `details`,
    `date`,
    `rating`,
    `studio_id`,
    `organized`,
    `created_at`,
    `updated_at`,
    `code`,
    `director`,
    `resume_time`,
    `play_duration`,
    `cover_blob`
  FROM `scenes`;

WITH max_view_count AS (
  SELECT MAX(play_count) AS max_count
  FROM scenes
), numbers AS (
  SELECT 1 AS n
  FROM max_view_count
  UNION ALL
  SELECT n + 1
  FROM numbers
  WHERE n < (SELECT max_count FROM max_view_count)
)
INSERT INTO scenes_view_dates (scene_id, view_date)
SELECT scenes.id, 
       CASE 
         WHEN numbers.n = scenes.play_count THEN COALESCE(scenes.last_played_at, scenes.created_at) 
         ELSE scenes.created_at
       END AS view_date
FROM scenes
JOIN numbers
WHERE numbers.n <= scenes.play_count;

WITH numbers AS (
  SELECT 1 AS n
  UNION ALL
  SELECT n + 1
  FROM numbers
  WHERE n < (SELECT MAX(o_counter) FROM scenes)
)
INSERT INTO scenes_o_dates (scene_id, o_date)
SELECT scenes.id, 
       CASE 
         WHEN numbers.n <= scenes.o_counter THEN scenes.created_at
       END AS o_date
FROM scenes
CROSS JOIN numbers
WHERE numbers.n <= scenes.o_counter;

DROP INDEX `index_scenes_on_studio_id`;
DROP TABLE `scenes`;
ALTER TABLE `scenes_new` rename to `scenes`;

CREATE INDEX `index_scenes_on_studio_id` on `scenes` (`studio_id`);

PRAGMA foreign_keys=ON;