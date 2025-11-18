CREATE TABLE `studio_urls` (
  `studio_id` integer NOT NULL,
  `position` integer NOT NULL,
  `url` varchar(255) NOT NULL,
  foreign key(`studio_id`) references `studios`(`id`) on delete CASCADE,
  PRIMARY KEY(`studio_id`, `position`, `url`)
);

CREATE INDEX `studio_urls_url` on `studio_urls` (`url`);

INSERT INTO `studio_urls`
  (
    `studio_id`,
    `position`,
    `url`
  )
  SELECT 
    `id`,
    '0',
    `url`
  FROM `studios`
  WHERE `studios`.`url` IS NOT NULL AND `studios`.`url` != '';

ALTER TABLE `studios` DROP COLUMN `url`;
