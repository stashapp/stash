ALTER TABLE `movies` RENAME TO `groups`;
ALTER TABLE `groups` RENAME COLUMN `synopsis` TO `description`;

DROP INDEX `index_movies_on_name`;
CREATE INDEX `index_groups_on_name` ON `groups`(`name`);
DROP INDEX `index_movies_on_studio_id`;
CREATE INDEX `index_groups_on_studio_id` on `groups` (`studio_id`);

ALTER TABLE `movie_urls` RENAME TO `group_urls`;
ALTER TABLE `group_urls` RENAME COLUMN `movie_id` TO `group_id`;

DROP INDEX `movie_urls_url`;
CREATE INDEX `group_urls_url` on `group_urls` (`url`);

ALTER TABLE `movies_tags` RENAME TO `groups_tags`;
ALTER TABLE `groups_tags` RENAME COLUMN `movie_id` TO `group_id`;

DROP INDEX `index_movies_tags_on_tag_id`;
CREATE INDEX `index_groups_tags_on_tag_id` on `groups_tags` (`tag_id`);
DROP INDEX `index_movies_tags_on_movie_id`;
CREATE INDEX `index_groups_tags_on_movie_id` on `groups_tags` (`group_id`);

ALTER TABLE `movies_scenes` RENAME TO `groups_scenes`;
ALTER TABLE `groups_scenes` RENAME COLUMN `movie_id` TO `group_id`;
