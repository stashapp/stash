CREATE TABLE `images` (
  `id` integer not null primary key autoincrement,
  `path` varchar(510) not null,
  `checksum` varchar(255) not null,
  `title` varchar(255),
  `rating` tinyint,
  `size` integer,
  `width` tinyint,
  `height` tinyint,
  `studio_id` integer,
  `o_counter` tinyint not null default 0,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL
);

CREATE INDEX `index_images_on_studio_id` on `images` (`studio_id`);

CREATE TABLE `performers_images` (
  `performer_id` integer,
  `image_id` integer,
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE,
  foreign key(`image_id`) references `images`(`id`) on delete CASCADE
);

CREATE INDEX `index_performers_images_on_image_id` on `performers_images` (`image_id`);
CREATE INDEX `index_performers_images_on_performer_id` on `performers_images` (`performer_id`);

CREATE TABLE `images_tags` (
  `image_id` integer,
  `tag_id` integer,
  foreign key(`image_id`) references `images`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE
);

CREATE INDEX `index_images_tags_on_tag_id` on `images_tags` (`tag_id`);
CREATE INDEX `index_images_tags_on_image_id` on `images_tags` (`image_id`);

-- need to recreate galleries to add foreign key
ALTER TABLE `galleries` rename to `_galleries_old`;

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
  `scene_id` integer,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`scene_id`) references `scenes`(`id`) on delete SET NULL,
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL
);

DROP INDEX IF EXISTS `index_galleries_on_scene_id`;
DROP INDEX IF EXISTS `galleries_path_unique`;
DROP INDEX IF EXISTS `galleries_checksum_unique`;

CREATE INDEX `index_galleries_on_scene_id` on `galleries` (`scene_id`);
CREATE UNIQUE INDEX `galleries_path_unique` on `galleries` (`path`);
CREATE UNIQUE INDEX `galleries_checksum_unique` on `galleries` (`checksum`);
CREATE INDEX `index_galleries_on_studio_id` on `galleries` (`studio_id`);

CREATE TABLE `galleries_images` (
  `gallery_id` integer,
  `image_id` integer,
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE,
  foreign key(`image_id`) references `images`(`id`) on delete CASCADE
);

CREATE INDEX `index_galleries_images_on_image_id` on `galleries_images` (`image_id`);
CREATE INDEX `index_galleries_images_on_gallery_id` on `galleries_images` (`gallery_id`);

CREATE TABLE `performers_galleries` (
  `performer_id` integer,
  `gallery_id` integer,
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE,
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE
);

CREATE INDEX `index_performers_galleries_on_gallery_id` on `performers_galleries` (`gallery_id`);
CREATE INDEX `index_performers_galleries_on_performer_id` on `performers_galleries` (`performer_id`);

CREATE TABLE `galleries_tags` (
  `gallery_id` integer,
  `tag_id` integer,
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE
);

CREATE INDEX `index_galleries_tags_on_tag_id` on `galleries_tags` (`tag_id`);
CREATE INDEX `index_galleries_tags_on_gallery_id` on `galleries_tags` (`gallery_id`);

INSERT INTO `galleries`
  (
    `id`,
    `path`,
    `checksum`,
    `scene_id`,
    `created_at`,
    `updated_at`
  )
  SELECT 
    `id`,
    `path`,
    `checksum`,
    `scene_id`,
    `created_at`,
    `updated_at`
  FROM `_galleries_old`;

DROP TABLE `_galleries_old`;
