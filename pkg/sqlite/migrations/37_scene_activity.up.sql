ALTER TABLE `scenes` ADD COLUMN `continue_position` float not null default 0;
ALTER TABLE `scenes` ADD COLUMN `last_played_at` datetime not null;
ALTER TABLE `scenes` ADD COLUMN `view_count` tinyint not null default 0;
ALTER TABLE `scenes` ADD COLUMN `watch_time` float not null default 0;