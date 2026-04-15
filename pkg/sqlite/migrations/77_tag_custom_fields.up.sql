CREATE TABLE `tag_custom_fields` (
  `tag_id` integer NOT NULL,
  `field` varchar(64) NOT NULL,
  `value` BLOB NOT NULL,
  PRIMARY KEY (`tag_id`, `field`),
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE
);

CREATE INDEX `index_tag_custom_fields_field_value` ON `tag_custom_fields` (`field`, `value`);