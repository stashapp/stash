-- Add performer_images table for multiple performer images
PRAGMA foreign_keys=OFF;

CREATE TABLE IF NOT EXISTS `performer_images` (
  `performer_id` integer not null,
  `image_blob` varchar(255) not null REFERENCES `blobs`(`checksum`),
  `sort_order` integer not null default 0,
  `created_at` datetime not null,
  foreign key (`performer_id`) references `performers`(`id`) on delete CASCADE
);

CREATE INDEX `performer_images_performer_id` on `performer_images` (`performer_id`);

-- Migrate existing single image to performer_images
INSERT INTO `performer_images` (`performer_id`, `image_blob`, `sort_order`, `created_at`)
SELECT `id`, `image_blob`, 0, `created_at`
FROM `performers`
WHERE `image_blob` IS NOT NULL;

PRAGMA foreign_keys=ON;
