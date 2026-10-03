-- add VR-related metadata columns to video_files.
-- nullable so existing rows are unchanged.
-- generic terminology (projection / stereo_mode / vr_corrections),
-- no VR-player-specific concepts live here.

ALTER TABLE `video_files` ADD COLUMN `projection` varchar(255);
ALTER TABLE `video_files` ADD COLUMN `stereo_mode` varchar(255);
ALTER TABLE `video_files` ADD COLUMN `vr_corrections` text;
