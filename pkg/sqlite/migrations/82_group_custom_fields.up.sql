CREATE TABLE `group_custom_fields` (
  `group_id` integer NOT NULL,
  `field` varchar(64) NOT NULL,
  `value` BLOB NOT NULL,
  PRIMARY KEY (`group_id`, `field`),
  foreign key(`group_id`) references `groups`(`id`) on delete CASCADE
);

CREATE INDEX `index_group_custom_fields_field_value` ON `group_custom_fields` (`field`, `value`);
