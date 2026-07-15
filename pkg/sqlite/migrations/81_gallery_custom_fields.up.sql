CREATE TABLE `gallery_custom_fields` (
  `gallery_id` integer NOT NULL,
  `field` varchar(64) NOT NULL,
  `value` BLOB NOT NULL,
  PRIMARY KEY (`gallery_id`, `field`),
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE
);

CREATE INDEX `index_gallery_custom_fields_field_value` ON `gallery_custom_fields` (`field`, `value`);
