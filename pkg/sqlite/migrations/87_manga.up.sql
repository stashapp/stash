-- Add manga table
PRAGMA foreign_keys=OFF;

CREATE TABLE `mangas` (
  `id` integer not null primary key autoincrement,
  `title` varchar(255),
  `url` varchar(255),
  `date` date,
  `details` text,
  `studio_id` integer,
  `rating` tinyint,
  `organized` boolean not null default '0',
  `created_at` datetime not null,
  `updated_at` datetime not null,
  `cover_image_blob` varchar(255) REFERENCES `blobs`(`checksum`),
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL
);

CREATE TABLE `mangas_tags` (
  `manga_id` integer,
  `tag_id` integer,
  foreign key(`manga_id`) references `mangas`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE,
  primary key (`manga_id`, `tag_id`)
);

CREATE TABLE `performers_mangas` (
  `performer_id` integer,
  `manga_id` integer,
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE,
  foreign key(`manga_id`) references `mangas`(`id`) on delete CASCADE,
  primary key (`performer_id`, `manga_id`)
);

CREATE TABLE `manga_urls` (
  `manga_id` integer not null,
  `url` varchar(255) not null,
  foreign key(`manga_id`) references `mangas`(`id`) on delete CASCADE
);

CREATE INDEX `mangas_title_idx` on `mangas` (`title`);
CREATE INDEX `mangas_studio_id_idx` on `mangas` (`studio_id`);

PRAGMA foreign_keys=ON;
