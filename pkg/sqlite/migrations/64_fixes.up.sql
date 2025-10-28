PRAGMA foreign_keys=OFF;

-- recreate scenes_view_dates adding not null to scene_id and adding indexes
CREATE TABLE `scenes_view_dates_new` (
  `scene_id` integer not null,
  `view_date` datetime not null,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE
);

INSERT INTO `scenes_view_dates_new`
  (
    `scene_id`,
    `view_date`
  )
  SELECT 
    `scene_id`,
    `view_date`
  FROM `scenes_view_dates`
  WHERE `scenes_view_dates`.`scene_id` IS NOT NULL;

DROP INDEX IF EXISTS `index_scenes_view_dates`;
DROP TABLE `scenes_view_dates`;
ALTER TABLE `scenes_view_dates_new` rename to `scenes_view_dates`;
CREATE INDEX `index_scenes_view_dates` ON `scenes_view_dates` (`scene_id`);

-- recreate scenes_o_dates adding not null to scene_id and adding indexes
CREATE TABLE `scenes_o_dates_new` (
  `scene_id` integer not null,
  `o_date` datetime not null,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE
);

INSERT INTO `scenes_o_dates_new`
  (
    `scene_id`,
    `o_date`
  )
  SELECT 
    `scene_id`,
    `o_date`
  FROM `scenes_o_dates`
  WHERE `scenes_o_dates`.`scene_id` IS NOT NULL;

DROP INDEX IF EXISTS `index_scenes_o_dates`;
DROP TABLE `scenes_o_dates`;
ALTER TABLE `scenes_o_dates_new` rename to `scenes_o_dates`;
CREATE INDEX `index_scenes_o_dates` ON `scenes_o_dates` (`scene_id`);

PRAGMA foreign_keys=ON;