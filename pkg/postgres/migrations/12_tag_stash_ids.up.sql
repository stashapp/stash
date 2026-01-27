-- 74_tag_stash_ids.up.sql
CREATE TABLE "tag_stash_ids" (
  "tag_id" integer,
  "endpoint" text,
  "stash_id" uuid,
  "updated_at" timestamp not null default '1970-01-01T00:00:00Z',
  foreign key("tag_id") references "tags"("id") on delete CASCADE
);

CREATE UNIQUE INDEX tag_stash_ids_unique_idx ON tag_stash_ids (tag_id, endpoint, stash_id);
