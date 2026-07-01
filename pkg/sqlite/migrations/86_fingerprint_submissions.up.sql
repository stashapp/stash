CREATE TABLE `fingerprint_submissions` (
  `endpoint` varchar(255) NOT NULL,
  `stash_id` varchar(36) NOT NULL,
  `scene_id` integer NOT NULL,
  `vote` varchar(20) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY (`endpoint`, `stash_id`),
  FOREIGN KEY(`scene_id`) REFERENCES `scenes`(`id`) ON DELETE CASCADE
);

CREATE INDEX `idx_fingerprint_submissions_endpoint` ON `fingerprint_submissions` (`endpoint`);
