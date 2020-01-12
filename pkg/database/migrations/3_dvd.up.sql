CREATE TABLE `dvds` (
  `id` integer not null primary key autoincrement,
  `name` varchar(255),
  `aliases` varchar(255),
  `durationdvd` varchar(6),
  `year` varchar(4),
  `director` varchar(255),
  `synopsis` text,
  `frontimage` blob not null,
  `backimage` blob,
  `checksum` varchar(255) not null,
  `url` varchar(255),
  `created_at` datetime not null,
  `updated_at` datetime not null
);
ALTER TABLE `scenes` ADD COLUMN `dvd_id` integer;
ALTER TABLE `scraped_items` ADD COLUMN `dvd_id` integer;
CREATE UNIQUE INDEX `dvds_checksum_unique` on `dvds` (`checksum`);
CREATE INDEX `index_dvds_on_name` on `dvds` (`name`);
CREATE INDEX `index_dvds_on_checksum` on `dvds` (`checksum`);