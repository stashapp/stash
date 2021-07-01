CREATE TABLE `saved_filters` (
  `id` integer not null primary key autoincrement,
  `name` varchar(510) not null,
  `mode` varchar(255) not null,
  `filter` blob not null
);

CREATE UNIQUE INDEX `index_saved_filters_on_mode_name_unique` on `saved_filters` (`mode`, `name`);
