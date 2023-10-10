PRAGMA foreign_keys=OFF;

CREATE TABLE `gallery_urls` (
  `gallery_id` integer NOT NULL,
  `position` integer NOT NULL,
  `url` varchar(255) NOT NULL,
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE,
  PRIMARY KEY(`gallery_id`, `position`, `url`)
);

CREATE INDEX `gallery_urls_url` on `gallery_urls` (`url`);

-- drop url
CREATE TABLE `galleries_new` (
  `id` integer not null primary key autoincrement,
  `folder_id` integer,
  `title` varchar(255),
  `date` date,
  `details` text,
  `studio_id` integer,
  `rating` tinyint,
  `organized` boolean not null default '0',
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL,
  foreign key(`folder_id`) references `folders`(`id`) on delete SET NULL
);

INSERT INTO `galleries_new`
  (
    `id`,
    `folder_id`,
    `title`,
    `date`,
    `details`,
    `studio_id`,
    `rating`,
    `organized`,
    `created_at`,
    `updated_at`
  )
  SELECT 
    `id`,
    `folder_id`,
    `title`,
    `date`,
    `details`,
    `studio_id`,
    `rating`,
    `organized`,
    `created_at`,
    `updated_at`
  FROM `galleries`;

INSERT INTO `gallery_urls`
  (
    `gallery_id`,
    `position`,
    `url`
  )
  SELECT 
    `id`,
    '0',
    `url`
  FROM `galleries`
  WHERE `galleries`.`url` IS NOT NULL AND `galleries`.`url` != '';

DROP INDEX `index_galleries_on_studio_id`;
DROP INDEX `index_galleries_on_folder_id_unique`;
DROP TABLE `galleries`;
ALTER TABLE `galleries_new` rename to `galleries`;

CREATE INDEX `index_galleries_on_studio_id` on `galleries` (`studio_id`);
CREATE UNIQUE INDEX `index_galleries_on_folder_id_unique` on `galleries` (`folder_id`);

PRAGMA foreign_keys=ON;
