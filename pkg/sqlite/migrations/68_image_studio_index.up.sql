-- with the existing index, if no images have a studio id, then the index is 
-- not used when filtering by studio id. The assumption with this change is that
-- most images don't have a studio id, so filtering by non-null studio id should
-- be faster with this index. This is a tradeoff, as filtering by null studio id
-- will be slower.
DROP INDEX index_images_on_studio_id;
CREATE INDEX `index_images_on_studio_id` on `images` (`studio_id`) WHERE `studio_id` IS NOT NULL;