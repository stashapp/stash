CREATE TABLE `performer_aliases` (
  `performer_id` integer NOT NULL,
  `alias` varchar(255) NOT NULL,
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE,
  PRIMARY KEY(`performer_id`, `alias`)
);

CREATE INDEX `performer_aliases_alias` on `performer_aliases` (`alias`);

-- `performers`.`aliases` will be dropped in the post-migration

DROP INDEX `performers_checksum_unique`;
ALTER TABLE `performers` DROP COLUMN `checksum`;
ALTER TABLE `performers` ADD COLUMN `disambiguation` varchar(255);

CREATE UNIQUE INDEX `performers_name_disambiguation_unique` on `performers` (`name`, `disambiguation`) WHERE `disambiguation` IS NOT NULL;
CREATE UNIQUE INDEX `performers_name_unique` on `performers` (`name`) WHERE `disambiguation` IS NULL;
