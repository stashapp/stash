CREATE TABLE `studio_custom_fields` (
  `studio_id` integer NOT NULL,
  `field` varchar(64) NOT NULL,
  `value` BLOB NOT NULL,
  PRIMARY KEY (`studio_id`, `field`),
  foreign key(`studio_id`) references `studios`(`id`) on delete CASCADE
);

CREATE INDEX `index_studio_custom_fields_field_value` ON `studio_custom_fields` (`field`, `value`);
