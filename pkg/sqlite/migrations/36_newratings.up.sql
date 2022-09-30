PRAGMA foreign_keys=off;
PRAGMA legacy_alter_table = TRUE;
CREATE TABLE "scenes_new" (
  `id` integer not null primary key autoincrement,
  -- REMOVED: `path` varchar(510) not null,
  -- REMOVED: `checksum` varchar(255),
  -- REMOVED: `oshash` varchar(255),
  `title` varchar(255),
  `details` text,
  `url` varchar(255),
  `date` date,
  `rating` tinyint NULL CHECK (`rating` >= 0 AND `rating` <= 100),
  -- REMOVED: `size` varchar(255),
  -- REMOVED: `duration` float,
  -- REMOVED: `video_codec` varchar(255),
  -- REMOVED: `audio_codec` varchar(255),
  -- REMOVED: `width` tinyint,
  -- REMOVED: `height` tinyint,
  -- REMOVED: `framerate` float,
  -- REMOVED: `bitrate` integer,
  `studio_id` integer,
  `o_counter` tinyint not null default 0,
  -- REMOVED: `format` varchar(255),
  `organized` boolean not null default '0',
  -- REMOVED: `interactive` boolean not null default '0',
  -- REMOVED: `interactive_speed` int,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  -- REMOVED: `file_mod_time` datetime,
  -- REMOVED: `phash` blob,
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL
  -- REMOVED: CHECK (`checksum` is not null or `oshash` is not null)
);
INSERT INTO "scenes_new" SELECT * FROM "scenes";
DROP TABLE "scenes";
ALTER TABLE "scenes_new" RENAME TO "scenes";

CREATE TABLE "galleries_new" (
  `id` integer not null primary key autoincrement,
  -- REMOVED: `path` varchar(510),
  -- REMOVED: `checksum` varchar(255) not null,
  -- REMOVED: `zip` boolean not null default '0',
  `folder_id` integer,
  `title` varchar(255),
  `url` varchar(255),
  `date` date,
  `details` text,
  `studio_id` integer,
  `rating` tinyint NULL CHECK (`rating` >= 0 AND `rating` <= 100),
  -- REMOVED: `file_mod_time` datetime,
  `organized` boolean not null default '0',
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL,
  foreign key(`folder_id`) references `folders`(`id`) on delete SET NULL
);
INSERT INTO "galleries_new" SELECT * FROM "galleries";
DROP TABLE "galleries";
ALTER TABLE "galleries_new" RENAME TO "galleries";

CREATE TABLE "images_new" (
  `id` integer not null primary key autoincrement,
  -- REMOVED: `path` varchar(510) not null,
  -- REMOVED: `checksum` varchar(255) not null,
  `title` varchar(255),
  `rating` tinyint NULL CHECK (`rating` >= 0 AND `rating` <= 100),
  -- REMOVED: `size` integer,
  -- REMOVED: `width` tinyint,
  -- REMOVED: `height` tinyint,
  `studio_id` integer,
  `o_counter` tinyint not null default 0,
  `organized` boolean not null default '0',
  -- REMOVED: `file_mod_time` datetime,
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`studio_id`) references `studios`(`id`) on delete SET NULL
);
INSERT INTO "images_new" SELECT * FROM "images";
DROP TABLE "images";
ALTER TABLE "images_new" RENAME TO "images";

CREATE TABLE "movies_new" (
  `id` integer not null primary key autoincrement,
  `name` varchar(255) not null,
  `aliases` varchar(255),
  `duration` integer,
  `date` date,
  `rating` tinyint NULL CHECK (`rating` >= 0 AND `rating` <= 100),
  `studio_id` integer,
  `director` varchar(255),
  `synopsis` text,
  `checksum` varchar(255) not null,
  `url` varchar(255),
  `created_at` datetime not null,
  `updated_at` datetime not null,
  foreign key(`studio_id`) references `studios`(`id`) on delete set null
);
INSERT INTO "movies_new" SELECT * FROM "movies";
DROP TABLE "movies";
ALTER TABLE "movies_new" RENAME TO "movies";

CREATE TABLE "performers_new" (
  `id` integer not null primary key autoincrement,
  `checksum` varchar(255) not null,
  `name` varchar(255),
  `gender` varchar(20),
  `url` varchar(255),
  `twitter` varchar(255),
  `instagram` varchar(255),
  `birthdate` date,
  `ethnicity` varchar(255),
  `country` varchar(255),
  `eye_color` varchar(255),
  `height` varchar(255),
  `measurements` varchar(255),
  `fake_tits` varchar(255),
  `career_length` varchar(255),
  `tattoos` varchar(255),
  `piercings` varchar(255),
  `aliases` varchar(255),
  `favorite` boolean not null default '0',
  `created_at` datetime not null,
  `updated_at` datetime not null, 
  `details` text, `death_date` date, 
  `hair_color` varchar(255), 
  `weight` integer, 
  `rating` tinyint NULL CHECK (`rating` >= 0 AND `rating` <= 100), 
  `ignore_auto_tag` boolean not null default '0'
  );
INSERT INTO "performers_new" SELECT * FROM "performers";
DROP TABLE "performers";
ALTER TABLE "performers_new" RENAME TO "performers";

CREATE TABLE "studios_new" (
  `id` integer not null primary key autoincrement,
  `checksum` varchar(255) not null,
  `name` varchar(255),
  `url` varchar(255),
  `parent_id` integer DEFAULT NULL CHECK ( id IS NOT parent_id ) REFERENCES studios(id) on delete set null,
  `created_at` datetime not null,
  `updated_at` datetime not null, 
  `details` text, 
  `rating` tinyint NULL CHECK (`rating` >= 0 AND `rating` <= 100), 
  `ignore_auto_tag` boolean not null default '0'
);
INSERT INTO "studios_new" SELECT * FROM "studios";
DROP TABLE "studios";
ALTER TABLE "studios_new" RENAME TO "studios";

PRAGMA foreign_keys=on;
PRAGMA legacy_alter_table = FALSE;

UPDATE `scenes` SET `rating` = (`rating` * 20) WHERE `rating` < 6;
UPDATE `galleries` SET `rating` = (`rating` * 20) WHERE `rating` < 6;
UPDATE `images` SET `rating` = (`rating` * 20) WHERE `rating` < 6;
UPDATE `movies` SET `rating` = (`rating` * 20) WHERE `rating` < 6;
UPDATE `performers` SET `rating` = (`rating` * 20) WHERE `rating` < 6;
UPDATE `studios` SET `rating` = (`rating` * 20) WHERE `rating` < 6;