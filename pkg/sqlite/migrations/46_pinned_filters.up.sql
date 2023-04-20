CREATE TABLE `pinned_filters` (
  `id` integer not null primary key autoincrement,
  `name` varchar(510) not null,
  `mode` varchar(255) not null
);

CREATE UNIQUE INDEX `index_pinned_filters_on_mode_name_unique` on `pinned_filters` (`mode`, `name`);
