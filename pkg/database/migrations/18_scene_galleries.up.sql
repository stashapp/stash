-- recreate the tables referencing galleries to correct their references
ALTER TABLE `galleries` rename to `_galleries_old`;
ALTER TABLE `galleries_images` rename to `_galleries_images_old`;
ALTER TABLE `galleries_tags` rename to `_galleries_tags_old`;
ALTER TABLE `performers_galleries` rename to `_performers_galleries_old`;

CREATE TABLE `galleries` (
  `id` integer not null primary key autoincrement,
  `path` varchar(510),
  `checksum` varchar(255) not null,
  `zip` boolean not null default '0',
  `title` varchar(255),
  `url` varchar(255),
  `date` date,
  `details` text,
  `studio_id` integer,
  `rating` tinyint,
  `file_mod_time` datetime,
  `organized` boolean not null default '0',
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL
);

DROP INDEX IF EXISTS `index_galleries_on_scene_id`;
DROP INDEX IF EXISTS `galleries_path_unique`;
DROP INDEX IF EXISTS `galleries_checksum_unique`;
DROP INDEX IF EXISTS `index_galleries_on_studio_id`;

CREATE UNIQUE INDEX `galleries_path_unique` on `galleries` (`path`);
CREATE UNIQUE INDEX `galleries_checksum_unique` on `galleries` (`checksum`);
CREATE INDEX `index_galleries_on_studio_id` on `galleries` (`studio_id`);

CREATE TABLE `scenes_galleries` (
  `scene_id` integer,
  `gallery_id` integer,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE,
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE
);

CREATE INDEX `index_scenes_galleries_on_scene_id` on `scenes_galleries` (`scene_id`);
CREATE INDEX `index_scenes_galleries_on_gallery_id` on `scenes_galleries` (`gallery_id`);

CREATE TABLE `galleries_images` (
  `gallery_id` integer,
  `image_id` integer,
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE,
  foreign key(`image_id`) references `images`(`id`) on delete CASCADE
);

DROP INDEX IF EXISTS `index_galleries_images_on_image_id`;
DROP INDEX IF EXISTS `index_galleries_images_on_gallery_id`;

CREATE INDEX `index_galleries_images_on_image_id` on `galleries_images` (`image_id`);
CREATE INDEX `index_galleries_images_on_gallery_id` on `galleries_images` (`gallery_id`);

CREATE TABLE `performers_galleries` (
  `performer_id` integer,
  `gallery_id` integer,
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE,
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE
);

DROP INDEX IF EXISTS `index_performers_galleries_on_gallery_id`;
DROP INDEX IF EXISTS `index_performers_galleries_on_performer_id`;

CREATE INDEX `index_performers_galleries_on_gallery_id` on `performers_galleries` (`gallery_id`);
CREATE INDEX `index_performers_galleries_on_performer_id` on `performers_galleries` (`performer_id`);

CREATE TABLE `galleries_tags` (
  `gallery_id` integer,
  `tag_id` integer,
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE
);

DROP INDEX IF EXISTS `index_galleries_tags_on_tag_id`;
DROP INDEX IF EXISTS `index_galleries_tags_on_gallery_id`;

CREATE INDEX `index_galleries_tags_on_tag_id` on `galleries_tags` (`tag_id`);
CREATE INDEX `index_galleries_tags_on_gallery_id` on `galleries_tags` (`gallery_id`);

-- populate from the old tables
INSERT INTO `galleries`
  (
    `id`,
    `path`,
    `checksum`,
    `zip`,
    `title`,
    `url`,
    `date`,
    `details`,
    `studio_id`,
    `rating`,
    `file_mod_time`,
    `organized`,
    `created_at`,
    `updated_at`
  )
  SELECT 
    `id`,
    `path`,
    `checksum`,
    `zip`,
    `title`,
    `url`,
    `date`,
    `details`,
    `studio_id`,
    `rating`,
    `file_mod_time`,
    `organized`,
    `created_at`,
    `updated_at`
  FROM `_galleries_old`;

INSERT INTO `scenes_galleries`
  (
    `scene_id`,
    `gallery_id`
  )
  SELECT
    `scene_id`,
    `id`
  FROM `_galleries_old`
  WHERE scene_id IS NOT NULL;

-- these tables are a direct copy
INSERT INTO `galleries_images` SELECT * from `_galleries_images_old`;
INSERT INTO `galleries_tags` SELECT * from `_galleries_tags_old`;
INSERT INTO `performers_galleries` SELECT * from `_performers_galleries_old`;

-- drop old tables
DROP TABLE `_galleries_old`;
DROP TABLE `_galleries_images_old`;
DROP TABLE `_galleries_tags_old`;
DROP TABLE `_performers_galleries_old`;
