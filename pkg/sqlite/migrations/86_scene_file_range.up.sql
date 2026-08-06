-- 86_scene_file_range.up.sql
-- Feature: Support multiple scenes in a single file (#3530)
-- Add nullable start_time/end_time columns to scenes_files so multiple
-- scenes can reference the same physical file with different time ranges.
-- Both NULL = full file (existing behaviour, backward compatible).
-- Only start_time = play from that point to end of file.
-- Both set = bounded range.
-- Ranges are valid only for a scene's primary file.

ALTER TABLE `scenes_files` ADD COLUMN `start_time` REAL NULL;
ALTER TABLE `scenes_files` ADD COLUMN `end_time` REAL NULL;
