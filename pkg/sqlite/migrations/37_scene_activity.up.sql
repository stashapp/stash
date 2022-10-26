ALTER TABLE `scenes` ADD COLUMN `continue_position` float not null default 0;
ALTER TABLE `scenes` ADD COLUMN `last_played_at` datetime not null;