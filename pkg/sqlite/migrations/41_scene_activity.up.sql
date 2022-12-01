ALTER TABLE `scenes` ADD COLUMN `resume_time` float not null default 0;
ALTER TABLE `scenes` ADD COLUMN `last_played_at` datetime default null;
ALTER TABLE `scenes` ADD COLUMN `play_count` tinyint not null default 0;
ALTER TABLE `scenes` ADD COLUMN `play_duration` float not null default 0;