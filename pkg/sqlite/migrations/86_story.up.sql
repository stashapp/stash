-- Add story table for text files/stories (#1259)
PRAGMA foreign_keys=OFF;

CREATE TABLE `stories` (
  `id` integer not null primary key autoincrement,
  `title` varchar(255),
  `author` varchar(255),
  `url` varchar(255),
  `date` date,
  `language` varchar(10),
  `tag_line` text,
  `details` text,
  `studio_id` integer,
  `rating` tinyint,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  `front_image_blob` varchar(255) REFERENCES `blobs`(`checksum`),
  `back_image_blob` varchar(255) REFERENCES `blobs`(`checksum`),
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL
);

CREATE TABLE `stories_tags` (
  `story_id` integer,
  `tag_id` integer,
  foreign key(`story_id`) references `stories`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE,
  primary key (`story_id`, `tag_id`)
);

CREATE TABLE `performers_stories` (
  `performer_id` integer,
  `story_id` integer,
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE,
  foreign key(`story_id`) references `stories`(`id`) on delete CASCADE,
  primary key (`performer_id`, `story_id`)
);

CREATE TABLE `story_urls` (
  `story_id` integer not null,
  `url` varchar(255) not null,
  foreign key(`story_id`) references `stories`(`id`) on delete CASCADE
);

CREATE INDEX `stories_title_idx` on `stories` (`title`);
CREATE INDEX `stories_studio_id_idx` on `stories` (`studio_id`);

PRAGMA foreign_keys=ON;
