CREATE TABLE `group_aliases` (
  `group_id` integer NOT NULL,
  `alias` varchar(255) NOT NULL,
  PRIMARY KEY(`group_id`, `alias`),
  foreign key(`group_id`) references `groups`(`id`) on delete CASCADE
);

CREATE INDEX `group_aliases_alias` on `group_aliases` (`alias`);

-- copy data from group table where aliases is just a string
-- skip NULL/empty rows to avoid violating the NOT NULL constraints
INSERT INTO `group_aliases` (
    `group_id`,
    `alias`
) SELECT
    `id`,
    `aliases`
FROM `groups`
WHERE `aliases` IS NOT NULL AND TRIM(`aliases`) != '';

ALTER TABLE `groups` DROP COLUMN `aliases`;