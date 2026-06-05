PRAGMA foreign_keys=OFF;

-- Update stories table: add audio_url for linking audio files, add performer_id for author
ALTER TABLE `stories` ADD COLUMN `audio_url` varchar(255);

PRAGMA foreign_keys=ON;
