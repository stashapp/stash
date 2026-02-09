UPDATE studios SET organized = 1 WHERE ignore_auto_tag = 1 AND organized = 0;
ALTER TABLE studios DROP COLUMN ignore_auto_tag;
