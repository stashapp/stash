ALTER TABLE studios 
    ADD COLUMN parent_id INTEGER DEFAULT NULL CHECK ( id IS NOT parent_id ) REFERENCES studios(id) on delete set null;
    CREATE INDEX index_studios_on_parent_id on studios (parent_id);