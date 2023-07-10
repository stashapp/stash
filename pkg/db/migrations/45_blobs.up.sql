CREATE TABLE `blobs` (
    `checksum` varchar(255) NOT NULL PRIMARY KEY,
    `blob` blob
);

ALTER TABLE `tags` ADD COLUMN `image_blob` varchar(255) REFERENCES `blobs`(`checksum`);
ALTER TABLE `studios` ADD COLUMN `image_blob` varchar(255) REFERENCES `blobs`(`checksum`);
ALTER TABLE `performers` ADD COLUMN `image_blob` varchar(255) REFERENCES `blobs`(`checksum`);
ALTER TABLE `scenes` ADD COLUMN `cover_blob` varchar(255) REFERENCES `blobs`(`checksum`);

ALTER TABLE `movies` ADD COLUMN `front_image_blob` varchar(255) REFERENCES `blobs`(`checksum`);
ALTER TABLE `movies` ADD COLUMN `back_image_blob` varchar(255) REFERENCES `blobs`(`checksum`);

-- performed in the post-migration
-- DROP TABLE `tags_image`;
-- DROP TABLE `studios_image`;
-- DROP TABLE `performers_image`;
-- DROP TABLE `scenes_cover`;
-- DROP TABLE `movies_images`;
