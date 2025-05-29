PRAGMA foreign_keys=OFF;

-- Create groups_aliases table
CREATE TABLE `groups_aliases` (
  `group_id` integer NOT NULL,
  `alias` varchar(255) NOT NULL,
  foreign key(`group_id`) references `groups`(`id`) on delete CASCADE,
  PRIMARY KEY(`group_id`, `alias`)
);

CREATE INDEX `groups_aliases_alias` on `groups_aliases` (`alias`);

-- Migrate existing aliases data
INSERT INTO `groups_aliases`
  (
    `group_id`,
    `alias`
  )
  SELECT 
    `id`,
    `aliases`
  FROM `groups`
  WHERE `groups`.`aliases` IS NOT NULL AND `groups`.`aliases` != '';

-- Create new groups table without aliases column
CREATE TABLE `groups_new` (
  `id` integer not null primary key autoincrement,
  `name` varchar(255) not null,
  `duration` integer,
  `date` date,
  `rating` tinyint,
  `studio_id` integer REFERENCES `studios`(`id`) ON DELETE SET NULL,
  `director` varchar(255),
  `description` text,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  `front_image_blob` varchar(255) REFERENCES `blobs`(`checksum`),
  `back_image_blob` varchar(255) REFERENCES `blobs`(`checksum`)
);

-- Copy data from old table to new table (excluding aliases)
INSERT INTO `groups_new`
  (
    `id`,
    `name`,
    `duration`,
    `date`,
    `rating`,
    `studio_id`,
    `director`,
    `description`,
    `created_at`,
    `updated_at`,
    `front_image_blob`,
    `back_image_blob`
  )
  SELECT 
    `id`,
    `name`,
    `duration`,
    `date`,
    `rating`,
    `studio_id`,
    `director`,
    `description`,
    `created_at`,
    `updated_at`,
    `front_image_blob`,
    `back_image_blob`
  FROM `groups`;

-- Drop old table and rename new one
DROP TABLE `groups`;
ALTER TABLE `groups_new` rename to `groups`;

-- Recreate indexes
CREATE INDEX `index_groups_on_name` ON `groups`(`name`);
CREATE INDEX `index_groups_on_studio_id` on `groups` (`studio_id`);

PRAGMA foreign_keys=ON; 