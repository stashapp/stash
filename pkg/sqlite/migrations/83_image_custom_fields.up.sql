CREATE TABLE `image_custom_fields` (
  `image_id` integer NOT NULL,
  `field` varchar(64) NOT NULL,
  `value` BLOB NOT NULL,
  PRIMARY KEY (`image_id`, `field`),
  foreign key(`image_id`) references `images`(`id`) on delete CASCADE
);

CREATE INDEX `index_image_custom_fields_field_value` ON `image_custom_fields` (`field`, `value`);
