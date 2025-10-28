PRAGMA foreign_keys=OFF;

CREATE TABLE `image_urls` (
  `image_id` integer NOT NULL,
  `position` integer NOT NULL,
  `url` varchar(255) NOT NULL,
  foreign key(`image_id`) references `images`(`id`) on delete CASCADE,
  PRIMARY KEY(`image_id`, `position`, `url`)
);

CREATE INDEX `image_urls_url` on `image_urls` (`url`);

-- drop url
CREATE TABLE "images_new" (
  `id` integer not null primary key autoincrement,
  `title` varchar(255),
  `rating` tinyint,
  `studio_id` integer,
  `o_counter` tinyint not null default 0,
  `organized` boolean not null default '0',
  `created_at` datetime not null,
  `updated_at` datetime not null, 
  `date` date,
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL
);

INSERT INTO `images_new`
  (
    `id`,
    `title`,
    `rating`,
    `studio_id`,
    `o_counter`,
    `organized`,
    `created_at`,
    `updated_at`,
    `date`
  )
  SELECT 
    `id`,
    `title`,
    `rating`,
    `studio_id`,
    `o_counter`,
    `organized`,
    `created_at`,
    `updated_at`,
    `date`
  FROM `images`;

INSERT INTO `image_urls`
  (
    `image_id`,
    `position`,
    `url`
  )
  SELECT 
    `id`,
    '0',
    `url`
  FROM `images`
  WHERE `images`.`url` IS NOT NULL AND `images`.`url` != '';

DROP INDEX `index_images_on_studio_id`;
DROP TABLE `images`;
ALTER TABLE `images_new` rename to `images`;

CREATE INDEX `index_images_on_studio_id` on `images` (`studio_id`);

PRAGMA foreign_keys=ON;
