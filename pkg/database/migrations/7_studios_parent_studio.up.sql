CREATE TABLE studios_dg_tmp
(
    id INTEGER NOT NULL
        PRIMARY KEY AUTOINCREMENT,
    image BLOB NOT NULL,
    checksum VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    url VARCHAR(255),
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    parent_studio_id INTEGER DEFAULT NULL
        CONSTRAINT studios_studios_id_fk
            REFERENCES studios
        CHECK ( id IS NOT parent_studio_id )
);

INSERT INTO studios_dg_tmp(id, image, checksum, name, url, created_at, updated_at) SELECT id, image, checksum, name, url, created_at, updated_at FROM studios;

DROP TABLE studios;

ALTER TABLE studios_dg_tmp RENAME TO studios;

CREATE INDEX index_studios_on_checksum
    ON studios (checksum);

CREATE INDEX index_studios_on_name
    ON studios (name);

CREATE UNIQUE INDEX studios_checksum_unique
    ON studios (checksum);