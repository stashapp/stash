-- revert VR metadata columns.
-- SQLite < 3.35 does not support DROP COLUMN; rely on full table
-- recreation in any future downgrade.

ALTER TABLE `video_files` DROP COLUMN `projection`;
ALTER TABLE `video_files` DROP COLUMN `stereo_mode`;
ALTER TABLE `video_files` DROP COLUMN `vr_corrections`;
