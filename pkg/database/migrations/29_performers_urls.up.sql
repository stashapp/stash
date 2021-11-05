CREATE TABLE `performers_urls` (
  `performer_id` integer NOT NULL,
  `url` varchar(255) NOT NULL,
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE
);

CREATE INDEX `index_performers_urls_on_performer_id` on `performers_urls` (`performer_id`);
