CREATE TABLE `scene_custom_fields` (
  `scene_id` integer NOT NULL,
  `field` varchar(64) NOT NULL,
  `value` BLOB NOT NULL,
  PRIMARY KEY (`scene_id`, `field`),
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE
);

CREATE INDEX `index_scene_custom_fields_field_value` ON `scene_custom_fields` (`field`, `value`);