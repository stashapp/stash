--------------------------------------------
-- audios definition
--
CREATE TABLE "audios" (
    `id` integer not null primary key autoincrement,
    `title` varchar(255),
    `details` text,
    `date` date,
    `rating` tinyint,
    `studio_id` integer,
    `organized` boolean not null default '0',
    `created_at` datetime not null,
    `updated_at` datetime not null,
    `code` text,
    `resume_time` float not null default 0,
    `play_duration` float not null default 0,
    "date_precision" TINYINT,
    foreign key(`studio_id`) references `studios`(`id`) on delete
    SET NULL
);
CREATE INDEX `index_audios_on_studio_id` on `audios` (`studio_id`);
--------------------------------------------
-- audios_o_dates definition
--
CREATE TABLE "audios_o_dates" (
    `audio_id` integer not null,
    `o_date` datetime not null,
    foreign key(`audio_id`) references `audios`(`id`) on delete CASCADE
);
CREATE INDEX `index_audios_o_dates` ON `audios_o_dates` (`audio_id`);
--------------------------------------------
-- audios_tags definition
--
CREATE TABLE "audios_tags" (
    `audio_id` integer,
    `tag_id` integer,
    foreign key(`audio_id`) references `audios`(`id`) on delete CASCADE,
    foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE,
    PRIMARY KEY(`audio_id`, `tag_id`)
);
CREATE INDEX `index_audios_tags_on_tag_id` on `audios_tags` (`tag_id`);
--------------------------------------------
-- audios_view_dates definition
--
CREATE TABLE "audios_view_dates" (
    `audio_id` integer not null,
    `view_date` datetime not null,
    foreign key(`audio_id`) references `audios`(`id`) on delete CASCADE
);
CREATE INDEX `index_audios_view_dates` ON `audios_view_dates` (`audio_id`);
--------------------------------------------
-- groups_audios definition
--
CREATE TABLE "groups_audios" (
    "group_id" integer,
    `audio_id` integer,
    `audio_index` tinyint,
    foreign key("group_id") references "groups"(`id`) on delete cascade,
    foreign key(`audio_id`) references `audios`(`id`) on delete cascade,
    PRIMARY KEY("group_id", `audio_id`)
);
CREATE INDEX `index_group_audios_on_group_id` on "groups_audios" ("group_id");
--------------------------------------------
-- performers_audios definition
--
CREATE TABLE "performers_audios" (
    `performer_id` integer,
    `audio_id` integer,
    foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE,
    foreign key(`audio_id`) references `audios`(`id`) on delete CASCADE,
    PRIMARY KEY (`audio_id`, `performer_id`)
);
CREATE INDEX `index_performers_audios_on_performer_id` on `performers_audios` (`performer_id`);
--------------------------------------------
-- audio_custom_fields definition
--
CREATE TABLE `audio_custom_fields` (
    `audio_id` integer NOT NULL,
    `field` varchar(64) NOT NULL,
    `value` BLOB NOT NULL,
    PRIMARY KEY (`audio_id`, `field`),
    foreign key(`audio_id`) references `audios`(`id`) on delete CASCADE
);
CREATE INDEX `index_audio_custom_fields_field_value` ON `audio_custom_fields` (`field`, `value`);
--------------------------------------------
-- audio_urls definition
--
CREATE TABLE `audio_urls` (
    `audio_id` integer NOT NULL,
    `position` integer NOT NULL,
    `url` varchar(255) NOT NULL,
    foreign key(`audio_id`) references `audios`(`id`) on delete CASCADE,
    PRIMARY KEY(`audio_id`, `position`, `url`)
);
CREATE INDEX `audio_urls_url` on `audio_urls` (`url`);
--------------------------------------------
-- audios_files definition
--
CREATE TABLE `audios_files` (
    `audio_id` integer NOT NULL,
    `file_id` integer NOT NULL,
    `primary` boolean NOT NULL,
    foreign key(`audio_id`) references `audios`(`id`) on delete CASCADE,
    foreign key(`file_id`) references `files`(`id`) on delete CASCADE,
    PRIMARY KEY(`audio_id`, `file_id`)
);
CREATE INDEX `index_audios_files_file_id` ON `audios_files` (`file_id`);
CREATE UNIQUE INDEX `unique_index_audios_files_on_primary` on `audios_files` (`audio_id`)
WHERE `primary` = 1;
--------------------------------------------
-- audio_files definition
--

-- TODO: think of better name for this, too close to `audios_files`
CREATE TABLE `audio_files` (
    `file_id` integer NOT NULL primary key,
    `duration` float NOT NULL,
    `format` varchar(255) NOT NULL,
    `audio_codec` varchar(255) NOT NULL,
    `sample_rate` float NOT NULL,
    `bit_rate` integer NOT NULL,
    foreign key(`file_id`) references `files`(`id`) on delete CASCADE
);