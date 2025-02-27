ALTER TABLE `performer_stash_ids` ADD COLUMN `updated_at` datetime not null default '1970-01-01T00:00:00Z';
ALTER TABLE `scene_stash_ids` ADD COLUMN `updated_at` datetime not null default '1970-01-01T00:00:00Z';
ALTER TABLE `studio_stash_ids` ADD COLUMN `updated_at` datetime not null default '1970-01-01T00:00:00Z';
