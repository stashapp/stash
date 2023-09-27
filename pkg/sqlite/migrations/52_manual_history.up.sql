CREATE TABLE `scenes_playdates` (
  `id` integer not null primary key autoincrement,
  `scene_id` integer,
  `playdate` datetime not null,
  foreign key(`scene_id`) references `scenes`(`id`)
);

CREATE TABLE `scenes_odates` (
  `id` integer not null primary key autoincrement,
  `scene_id` integer,
  `odate` datetime not null,
  foreign key(`scene_id`) references `scenes`(`id`)
);

WITH max_play_count AS (
  SELECT MAX(play_count) AS max_count
  FROM scenes
), numbers AS (
  SELECT 1 AS n
  FROM max_play_count
  UNION ALL
  SELECT n + 1
  FROM numbers
  WHERE n < (SELECT max_count FROM max_play_count)
)
INSERT INTO scenes_playdates (scene_id, playdate)
SELECT scenes.id, 
       CASE 
         WHEN numbers.n = scenes.play_count THEN scenes.last_played_at
         ELSE scenes.created_at
       END AS playdate
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
INSERT INTO scenes_odates (scene_id, odate)
SELECT scenes.id, 
       CASE 
         WHEN numbers.n <= scenes.o_counter THEN scenes.created_at
       END AS odate
FROM scenes
CROSS JOIN numbers
WHERE numbers.n <= scenes.o_counter;