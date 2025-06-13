PRAGMA foreign_keys=OFF;

-- Create studios_scenes join table
CREATE TABLE `studios_scenes` (
  `studio_id` integer,
  `scene_id` integer,
  foreign key(`studio_id`) references `studios`(`id`) on delete CASCADE,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE,
  UNIQUE(`studio_id`, `scene_id`)
);

-- Create studios_galleries join table
CREATE TABLE `studios_galleries` (
  `studio_id` integer,
  `gallery_id` integer,
  foreign key(`studio_id`) references `studios`(`id`) on delete CASCADE,
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE,
  UNIQUE(`studio_id`, `gallery_id`)
);

-- Create studios_images join table
CREATE TABLE `studios_images` (
  `studio_id` integer,
  `image_id` integer,
  foreign key(`studio_id`) references `studios`(`id`) on delete CASCADE,
  foreign key(`image_id`) references `images`(`id`) on delete CASCADE,
  UNIQUE(`studio_id`, `image_id`)
);

-- Create studios_groups join table
CREATE TABLE `studios_groups` (
  `studio_id` integer,
  `group_id` integer,
  foreign key(`studio_id`) references `studios`(`id`) on delete CASCADE,
  foreign key(`group_id`) references `groups`(`id`) on delete CASCADE,
  UNIQUE(`studio_id`, `group_id`)
);

-- Create indexes for the new join tables
CREATE INDEX `index_studios_scenes_on_studio_id` on `studios_scenes` (`studio_id`);
CREATE INDEX `index_studios_scenes_on_scene_id` on `studios_scenes` (`scene_id`);

CREATE INDEX `index_studios_galleries_on_studio_id` on `studios_galleries` (`studio_id`);
CREATE INDEX `index_studios_galleries_on_gallery_id` on `studios_galleries` (`gallery_id`);

CREATE INDEX `index_studios_images_on_studio_id` on `studios_images` (`studio_id`);
CREATE INDEX `index_studios_images_on_image_id` on `studios_images` (`image_id`);

CREATE INDEX `index_studios_groups_on_studio_id` on `studios_groups` (`studio_id`);
CREATE INDEX `index_studios_groups_on_group_id` on `studios_groups` (`group_id`);

-- Migrate existing studio_id data to the new join tables
INSERT INTO `studios_scenes` (`studio_id`, `scene_id`)
SELECT `studio_id`, `id`
FROM `scenes`
WHERE `studio_id` IS NOT NULL;

INSERT INTO `studios_galleries` (`studio_id`, `gallery_id`)
SELECT `studio_id`, `id`
FROM `galleries`
WHERE `studio_id` IS NOT NULL;

INSERT INTO `studios_images` (`studio_id`, `image_id`)
SELECT `studio_id`, `id`
FROM `images`
WHERE `studio_id` IS NOT NULL;

INSERT INTO `studios_groups` (`studio_id`, `group_id`)
SELECT `studio_id`, `id`
FROM `groups`
WHERE `studio_id` IS NOT NULL;

-- Create new scenes table without studio_id
CREATE TABLE "scenes_new" (
  `id` integer not null primary key autoincrement,
  `title` varchar(255),
  `details` text,
  `date` date,
  `rating` tinyint,
  `organized` boolean not null default '0',
  `created_at` datetime not null,
  `updated_at` datetime not null, 
  `code` text, 
  `director` text, 
  `resume_time` float not null default 0, 
  `play_duration` float not null default 0, 
  `cover_blob` varchar(255) REFERENCES `blobs`(`checksum`)
);

-- Copy data from old scenes table (excluding studio_id)
INSERT INTO `scenes_new`
  (
    `id`,
    `title`,
    `details`,
    `date`,
    `rating`,
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
    `organized`,
    `created_at`,
    `updated_at`,
    `code`,
    `director`,
    `resume_time`,
    `play_duration`,
    `cover_blob`
  FROM `scenes`;

-- Drop old scenes table and rename new one
DROP INDEX `index_scenes_on_studio_id`;
DROP TABLE `scenes`;
ALTER TABLE `scenes_new` rename to `scenes`;

-- Create new galleries table without studio_id
CREATE TABLE "galleries_new" (
  `id` integer not null primary key autoincrement,
  `title` varchar(255),
  `code` text,
  `date` date,
  `details` text,
  `photographer` text,
  `rating` tinyint,
  `organized` boolean not null default '0',
  `created_at` datetime not null,
  `updated_at` datetime not null,
  `folder_id` integer REFERENCES `folders`(`id`) ON DELETE CASCADE
);

-- Copy data from old galleries table (excluding studio_id)
INSERT INTO `galleries_new`
  (
    `id`,
    `title`,
    `code`,
    `date`,
    `details`,
    `photographer`,
    `rating`,
    `organized`,
    `created_at`,
    `updated_at`,
    `folder_id`
  )
  SELECT 
    `id`,
    `title`,
    `code`,
    `date`,
    `details`,
    `photographer`,
    `rating`,
    `organized`,
    `created_at`,
    `updated_at`,
    `folder_id`
  FROM `galleries`;

-- Drop old galleries table and rename new one
DROP INDEX IF EXISTS `index_galleries_on_studio_id`;
DROP TABLE `galleries`;
ALTER TABLE `galleries_new` rename to `galleries`;

-- Create new images table without studio_id
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
  `created_at` datetime not null,
  `updated_at` datetime not null
);

-- Copy data from old images table (excluding studio_id)
INSERT INTO `images_new`
  (
    `id`,
    `title`,
    `code`,
    `date`,
    `details`,
    `photographer`,
    `rating`,
    `organized`,
    `o_counter`,
    `created_at`,
    `updated_at`
  )
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
    `created_at`,
    `updated_at`
  FROM `images`;

-- Drop old images table and rename new one
DROP INDEX IF EXISTS `index_images_on_studio_id`;
DROP TABLE `images`;
ALTER TABLE `images_new` rename to `images`;

-- Create new groups table without studio_id
CREATE TABLE "groups_new" (
  `id` integer not null primary key autoincrement,
  `name` varchar(255) not null,
  `aliases` text,
  `duration` integer,
  `date` date,
  `rating` tinyint,
  `director` text,
  `description` text,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  `front_image_blob` varchar(255) REFERENCES `blobs`(`checksum`),
  `back_image_blob` varchar(255) REFERENCES `blobs`(`checksum`)
);

-- Copy data from old groups table (excluding studio_id)
INSERT INTO `groups_new`
  (
    `id`,
    `name`,
    `aliases`,
    `duration`,
    `date`,
    `rating`,
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
    `aliases`,
    `duration`,
    `date`,
    `rating`,
    `director`,
    `description`,
    `created_at`,
    `updated_at`,
    `front_image_blob`,
    `back_image_blob`
  FROM `groups`;

-- Drop old groups table and rename new one
DROP INDEX IF EXISTS `index_groups_on_studio_id`;
DROP TABLE `groups`;
ALTER TABLE `groups_new` rename to `groups`;

PRAGMA foreign_keys=ON; 