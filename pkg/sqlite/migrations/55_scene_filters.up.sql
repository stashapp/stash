PRAGMA foreign_keys=ON;

CREATE TABLE IF NOT EXISTS scene_filters (
    'id' INTEGER PRIMARY KEY autoincrement,
    'scene_id' INTEGER NOT NULL UNIQUE,
    'brightness' INTEGER NOT NULL,
    'contrast' INTEGER NOT NULL,
    'gamma' INTEGER NOT NULL,
    'saturate' INTEGER NOT NULL,
    'hue_rotate' INTEGER NOT NULL,
    'warmth' INTEGER NOT NULL,
    'red' INTEGER NOT NULL,
    'green' INTEGER NOT NULL,
    'blue' INTEGER NOT NULL,
    'blur' INTEGER NOT NULL,
    'rotate' REAL NOT NULL,
    'scale' INTEGER NOT NULL,
    'aspect_ratio' INTEGER NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    FOREIGN KEY ('scene_id') REFERENCES 'scenes'('id') ON DELETE CASCADE
);

CREATE INDEX `index_scene_filters_on_scene_id` ON `scene_filters` (`scene_id`);