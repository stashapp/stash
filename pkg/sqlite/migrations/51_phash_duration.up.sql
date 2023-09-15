-- Migrated all phash to phash-old for videos shorter than 2.5 min
UPDATE `files_fingerprints`
SET `type` = 'phash-old'
WHERE `type` = 'phash'
  AND `file_id` IN (
    SELECT `file_id`
    FROM `video_files`
    WHERE `duration` < 151
);
