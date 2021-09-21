DROP INDEX IF EXISTS `images_path_unique`;

CREATE UNIQUE INDEX `images_path_unique` ON `images` (`path`);
