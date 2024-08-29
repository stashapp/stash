CREATE TABLE `groups_relations` (
  `containing_id` integer not null,
  `sub_id` integer not null,
  `order_index` integer not null,
  `description` varchar(255),
  primary key (`containing_id`, `sub_id`),
  foreign key (`containing_id`) references `groups`(`id`) on delete cascade,
  foreign key (`sub_id`) references `groups`(`id`) on delete cascade,
  check (`containing_id` != `sub_id`)
);

CREATE INDEX `index_groups_relations_sub_id` ON `groups_relations` (`sub_id`);
CREATE UNIQUE INDEX `index_groups_relations_order_index_unique` ON `groups_relations` (`containing_id`, `order_index`);
