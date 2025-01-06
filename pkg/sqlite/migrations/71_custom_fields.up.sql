CREATE TABLE `performer_custom_fields` (
  `performer_id` integer NOT NULL,
  `field` varchar(64) NOT NULL,
  `value` BLOB NOT NULL,
  PRIMARY KEY (`performer_id`, `field`),
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE
);

CREATE INDEX `index_performer_custom_fields_field_value` ON `performer_custom_fields` (`field`, `value`);
