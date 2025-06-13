PRAGMA foreign_keys=OFF;

-- Create new scenes table with studio_id column
CREATE TABLE "scenes_new" (
  `id` integer not null primary key autoincrement,
  `title` varchar(255),
  `code` varchar(255),
  `details` text,
  `director` varchar(255),
  `date` date,
  `rating` tinyint,
  `organized` boolean not null default '0',
  `studio_id` integer,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  `resume_time` float not null default 0,
  `play_duration` float not null default 0,
  foreign key(`studio_id`) references `studios`(`id`)
);

-- Copy data from old scenes table
INSERT INTO "scenes_new" 
SELECT 
  `id`,
  `title`,
  `code`,
  `details`,
  `director`,
  `date`,
  `rating`,
  `organized`,
  NULL as `studio_id`,
  `created_at`,
  `updated_at`,
  `resume_time`,
  `play_duration`
FROM "scenes";

-- Update studio_id from studios_scenes table (take first studio for each scene)
UPDATE "scenes_new" 
SET `studio_id` = (
  SELECT `studio_id` 
  FROM `studios_scenes` 
  WHERE `studios_scenes`.`scene_id` = "scenes_new".`id` 
  LIMIT 1
);

-- Drop old scenes table and rename new one
DROP TABLE "scenes";
ALTER TABLE "scenes_new" RENAME TO "scenes";

-- Create new galleries table with studio_id column
CREATE TABLE "galleries_new" (
  `id` integer not null primary key autoincrement,
  `title` varchar(255),
  `code` text,
  `date` date,
  `details` text,
  `photographer` text,
  `rating` tinyint,
  `organized` boolean not null default '0',
  `studio_id` integer,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  `folder_id` integer REFERENCES `folders`(`id`) ON DELETE CASCADE,
  foreign key(`studio_id`) references `studios`(`id`)
);

-- Copy data from old galleries table
INSERT INTO "galleries_new" 
SELECT 
  `id`,
  `title`,
  `code`,
  `date`,
  `details`,
  `photographer`,
  `rating`,
  `organized`,
  NULL as `studio_id`,
  `created_at`,
  `updated_at`,
  `folder_id`
FROM "galleries";

-- Update studio_id from studios_galleries table (take first studio for each gallery)
UPDATE "galleries_new" 
SET `studio_id` = (
  SELECT `studio_id` 
  FROM `studios_galleries` 
  WHERE `studios_galleries`.`gallery_id` = "galleries_new".`id` 
  LIMIT 1
);

-- Drop old galleries table and rename new one
DROP TABLE "galleries";
ALTER TABLE "galleries_new" RENAME TO "galleries";

-- Create new images table with studio_id column
CREATE TABLE "images_new" (
  `id` integer not null primary key autoincrement,
  `title` varchar(255),
  `code` text,
  `date` date,
  `details` text,
  `photographer` text,
  `rating` tinyint,
  `organized` boolean not null default '0',
  `o_counter` tinyint not null default 0,
  `studio_id` integer,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`studio_id`) references `studios`(`id`)
);

-- Copy data from old images table
INSERT INTO "images_new" 
SELECT 
  `id`,
  `title`,
  `code`,
  `date`,
  `details`,
  `photographer`,
  `rating`,
  `organized`,
  `o_counter`,
  NULL as `studio_id`,
  `created_at`,
  `updated_at`
FROM "images";

-- Update studio_id from studios_images table (take first studio for each image)
UPDATE "images_new" 
SET `studio_id` = (
  SELECT `studio_id` 
  FROM `studios_images` 
  WHERE `studios_images`.`image_id` = "images_new".`id` 
  LIMIT 1
);

-- Drop old images table and rename new one
DROP TABLE "images";
ALTER TABLE "images_new" RENAME TO "images";

-- Create new groups table with studio_id column
CREATE TABLE "groups_new" (
  `id` integer not null primary key autoincrement,
  `name` varchar(255) not null,
  `aliases` text,
  `duration` integer,
  `date` date,
  `rating` tinyint,
  `studio_id` integer,
  `director` text,
  `description` text,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  `front_image_blob` varchar(255) REFERENCES `blobs`(`checksum`),
  `back_image_blob` varchar(255) REFERENCES `blobs`(`checksum`),
  foreign key(`studio_id`) references `studios`(`id`)
);

-- Copy data from old groups table
INSERT INTO "groups_new" 
SELECT 
  `id`,
  `name`,
  `aliases`,
  `duration`,
  `date`,
  `rating`,
  NULL as `studio_id`,
  `director`,
  `description`,
  `created_at`,
  `updated_at`,
  `front_image_blob`,
  `back_image_blob`
FROM "groups";

-- Update studio_id from studios_groups table (take first studio for each group)
UPDATE "groups_new" 
SET `studio_id` = (
  SELECT `studio_id` 
  FROM `studios_groups` 
  WHERE `studios_groups`.`group_id` = "groups_new".`id` 
  LIMIT 1
);

-- Drop old groups table and rename new one
DROP TABLE "groups";
ALTER TABLE "groups_new" RENAME TO "groups";

-- Recreate indexes
CREATE INDEX `index_scenes_on_studio_id` on `scenes` (`studio_id`);
CREATE INDEX `index_galleries_on_studio_id` on `galleries` (`studio_id`);
CREATE INDEX `index_images_on_studio_id` on `images` (`studio_id`);
CREATE INDEX `index_groups_on_studio_id` on `groups` (`studio_id`);

-- Drop join tables
DROP TABLE `studios_scenes`;
DROP TABLE `studios_galleries`;
DROP TABLE `studios_images`;
DROP TABLE `studios_groups`;

PRAGMA foreign_keys=ON; 