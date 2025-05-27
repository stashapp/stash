PRAGMA foreign_keys=OFF;

CREATE TABLE `studio_urls` (
  `studio_id` integer NOT NULL,
  `position` integer NOT NULL,
  `url` varchar(255) NOT NULL,
  foreign key(`studio_id`) references `studios`(`id`) on delete CASCADE,
  PRIMARY KEY(`studio_id`, `position`, `url`)
);

CREATE INDEX `studio_urls_url` on `studio_urls` (`url`);

-- insert existing URLs from studios.url into studio_urls table
INSERT INTO `studio_urls`
  (
    `studio_id`,
    `position`,
    `url`
  )
  SELECT 
    `id`,
    0,
    `url`
  FROM `studios`
  WHERE `studios`.`url` IS NOT NULL AND `studios`.`url` != '';

PRAGMA foreign_keys=ON; 