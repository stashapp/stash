ALTER TABLE `scenes` ADD COLUMN `cover_image_source` text CHECK (
  `cover_image_source` IS NULL
  OR `cover_image_source` = 'default'
  OR `cover_image_source` = 'clipboard'
  OR `cover_image_source` = 'userscript'
  OR `cover_image_source` LIKE 'url:%'
  OR `cover_image_source` LIKE 'stash:%'
  OR `cover_image_source` LIKE 'timestamp:%'
  OR `cover_image_source` LIKE 'file:%'
);
