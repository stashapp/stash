ALTER TABLE tags ADD COLUMN `allow_autotag` INT(1) DEFAULT 1 NOT NULL;

CREATE INDEX `images_autotag_covering_index` ON `images` (`organized`, `path`, `id`);

CREATE INDEX `scenes_autotag_covering_index` ON `scenes` (`organized`, `path`, `id`);