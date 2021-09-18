DROP INDEX IF EXISTS `images_path_unique`;
DROP INDEX IF EXISTS `images_checksum_unique`;

CREATE UNIQUE INDEX `images_path_unique` ON `images` (`path`);
CREATE UNIQUE INDEX `images_checksum_unique` ON `images` (`checksum`);
