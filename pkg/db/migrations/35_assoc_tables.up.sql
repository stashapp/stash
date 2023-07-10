-- add primary keys to association tables that are missing them
PRAGMA foreign_keys=OFF;

CREATE TABLE `performers_image_new` (
  `performer_id` integer primary key,
  `image` blob not null,
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE
);

INSERT INTO `performers_image_new`
  (
    `performer_id`,
    `image`
  )
  SELECT 
    `performer_id`,
    `image`
  FROM `performers_image` WHERE
  `performer_id` IS NOT NULL;

DROP TABLE `performers_image`;
ALTER TABLE `performers_image_new` rename to `performers_image`;

-- the following index is removed in favour of primary key
-- CREATE UNIQUE INDEX `index_performer_image_on_performer_id` on `performers_image` (`performer_id`);


CREATE TABLE `studios_image_new` (
  `studio_id` integer primary key,
  `image` blob not null,
  foreign key(`studio_id`) references `studios`(`id`) on delete CASCADE
);

INSERT INTO `studios_image_new`
  (
    `studio_id`,
    `image`
  )
  SELECT 
    `studio_id`,
    `image`
  FROM `studios_image` WHERE
  `studio_id` IS NOT NULL;

DROP TABLE `studios_image`;
ALTER TABLE `studios_image_new` rename to `studios_image`;

-- the following index is removed in favour of primary key
-- CREATE UNIQUE INDEX `index_studio_image_on_studio_id` on `studios_image` (`studio_id`);


CREATE TABLE `movies_images_new` (
  `movie_id` integer primary key,
  `front_image` blob not null,
  `back_image` blob,
  foreign key(`movie_id`) references `movies`(`id`) on delete CASCADE
);

INSERT INTO `movies_images_new`
  (
    `movie_id`,
    `front_image`,
    `back_image`
  )
  SELECT 
    `movie_id`,
    `front_image`,
    `back_image`
  FROM `movies_images` WHERE
  `movie_id` IS NOT NULL;

DROP TABLE `movies_images`;
ALTER TABLE `movies_images_new` rename to `movies_images`;

-- the following index is removed in favour of primary key
-- CREATE UNIQUE INDEX `index_movie_images_on_movie_id` on `movies_images` (`movie_id`);


CREATE TABLE `tags_image_new` (
  `tag_id` integer primary key,
  `image` blob not null,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE
);

INSERT INTO `tags_image_new`
  (
    `tag_id`,
    `image`
  )
  SELECT 
    `tag_id`,
    `image`
  FROM `tags_image` WHERE
  `tag_id` IS NOT NULL;

DROP TABLE `tags_image`;
ALTER TABLE `tags_image_new` rename to `tags_image`;

-- the following index is removed in favour of primary key
-- CREATE UNIQUE INDEX `index_tag_image_on_tag_id` on `tags_image` (`tag_id`);

-- add on delete cascade to foreign keys
CREATE TABLE `performers_scenes_new` (
  `performer_id` integer,
  `scene_id` integer,
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE,
  PRIMARY KEY (`scene_id`, `performer_id`)
);

INSERT INTO `performers_scenes_new`
  (
    `performer_id`,
    `scene_id`
  )
  SELECT 
    `performer_id`,
    `scene_id`
  FROM `performers_scenes` WHERE 
  `performer_id` IS NOT NULL AND `scene_id` IS NOT NULL
  ON CONFLICT (`scene_id`, `performer_id`) DO NOTHING;

DROP TABLE `performers_scenes`;
ALTER TABLE `performers_scenes_new` rename to `performers_scenes`;

CREATE INDEX `index_performers_scenes_on_performer_id` on `performers_scenes` (`performer_id`);

-- the following index is removed in favour of primary key
-- CREATE INDEX `index_performers_scenes_on_scene_id` on `performers_scenes` (`scene_id`);


CREATE TABLE `scene_markers_tags_new` (
  `scene_marker_id` integer,
  `tag_id` integer,
  foreign key(`scene_marker_id`) references `scene_markers`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE,
  PRIMARY KEY(`scene_marker_id`, `tag_id`)
);

INSERT INTO `scene_markers_tags_new`
  (
    `scene_marker_id`,
    `tag_id`
  )
  SELECT 
    `scene_marker_id`,
    `tag_id`
  FROM `scene_markers_tags` WHERE 
  `scene_marker_id` IS NOT NULL AND `tag_id` IS NOT NULL
  ON CONFLICT (`scene_marker_id`, `tag_id`) DO NOTHING;

DROP TABLE `scene_markers_tags`;
ALTER TABLE `scene_markers_tags_new` rename to `scene_markers_tags`;

CREATE INDEX `index_scene_markers_tags_on_tag_id` on `scene_markers_tags` (`tag_id`);

-- the following index is removed in favour of primary key
-- CREATE INDEX `index_scene_markers_tags_on_scene_marker_id` on `scene_markers_tags` (`scene_marker_id`);

-- add delete cascade to tag_id
CREATE TABLE `scenes_tags_new` (
  `scene_id` integer,
  `tag_id` integer,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE,
  PRIMARY KEY(`scene_id`, `tag_id`)
);

INSERT INTO `scenes_tags_new`
  (
    `scene_id`,
    `tag_id`
  )
  SELECT 
    `scene_id`,
    `tag_id`
  FROM `scenes_tags` WHERE 
  `scene_id` IS NOT NULL AND `tag_id` IS NOT NULL
  ON CONFLICT (`scene_id`, `tag_id`) DO NOTHING;

DROP TABLE `scenes_tags`;
ALTER TABLE `scenes_tags_new` rename to `scenes_tags`;

CREATE INDEX `index_scenes_tags_on_tag_id` on `scenes_tags` (`tag_id`);

-- the following index is removed in favour of primary key
-- CREATE INDEX `index_scenes_tags_on_scene_id` on `scenes_tags` (`scene_id`);


CREATE TABLE `movies_scenes_new` (
  `movie_id` integer,
  `scene_id` integer,
  `scene_index` tinyint,
  foreign key(`movie_id`) references `movies`(`id`) on delete cascade,
  foreign key(`scene_id`) references `scenes`(`id`) on delete cascade,
  PRIMARY KEY(`movie_id`, `scene_id`)
);

INSERT INTO `movies_scenes_new`
  (
    `movie_id`,
    `scene_id`,
    `scene_index`
  )
  SELECT 
    `movie_id`,
    `scene_id`,
    `scene_index`
  FROM `movies_scenes` WHERE 
  `movie_id` IS NOT NULL AND `scene_id` IS NOT NULL
  ON CONFLICT (`movie_id`, `scene_id`) DO NOTHING;

DROP TABLE `movies_scenes`;
ALTER TABLE `movies_scenes_new` rename to `movies_scenes`;

CREATE INDEX `index_movies_scenes_on_movie_id` on `movies_scenes` (`movie_id`);

-- the following index is removed in favour of primary key
-- CREATE INDEX `index_movies_scenes_on_scene_id` on `movies_scenes` (`scene_id`);


CREATE TABLE `scenes_cover_new` (
  `scene_id` integer primary key,
  `cover` blob not null,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE
);

INSERT INTO `scenes_cover_new`
  (
    `scene_id`,
    `cover`
  )
  SELECT 
    `scene_id`,
    `cover`
  FROM `scenes_cover` WHERE
  `scene_id` IS NOT NULL;

DROP TABLE `scenes_cover`;
ALTER TABLE `scenes_cover_new` rename to `scenes_cover`;

-- the following index is removed in favour of primary key
-- CREATE UNIQUE INDEX `index_scene_covers_on_scene_id` on `scenes_cover` (`scene_id`);


CREATE TABLE `performers_images_new` (
  `performer_id` integer,
  `image_id` integer,
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE,
  foreign key(`image_id`) references `images`(`id`) on delete CASCADE,
  PRIMARY KEY(`image_id`, `performer_id`)
);

INSERT INTO `performers_images_new`
  (
    `performer_id`,
    `image_id`
  )
  SELECT 
    `performer_id`,
    `image_id`
  FROM `performers_images` WHERE 
  `performer_id` IS NOT NULL AND `image_id` IS NOT NULL
  ON CONFLICT (`image_id`, `performer_id`) DO NOTHING;

DROP TABLE `performers_images`;
ALTER TABLE `performers_images_new` rename to `performers_images`;

CREATE INDEX `index_performers_images_on_performer_id` on `performers_images` (`performer_id`);

-- the following index is removed in favour of primary key
-- CREATE INDEX `index_performers_images_on_image_id` on `performers_images` (`image_id`);


CREATE TABLE `images_tags_new` (
  `image_id` integer,
  `tag_id` integer,
  foreign key(`image_id`) references `images`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE,
  PRIMARY KEY(`image_id`, `tag_id`)
);

INSERT INTO `images_tags_new`
  (
    `image_id`,
    `tag_id`
  )
  SELECT 
    `image_id`,
    `tag_id`
  FROM `images_tags` WHERE 
  `image_id` IS NOT NULL AND `tag_id` IS NOT NULL
  ON CONFLICT (`image_id`, `tag_id`) DO NOTHING;

DROP TABLE `images_tags`;
ALTER TABLE `images_tags_new` rename to `images_tags`;

CREATE INDEX `index_images_tags_on_tag_id` on `images_tags` (`tag_id`);

-- the following index is removed in favour of primary key
-- CREATE INDEX `index_images_tags_on_image_id` on `images_tags` (`image_id`);


CREATE TABLE `scene_stash_ids_new` (
  `scene_id` integer NOT NULL,
  `endpoint` varchar(255) NOT NULL,
  `stash_id` varchar(36) NOT NULL,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE,
  PRIMARY KEY(`scene_id`, `endpoint`)
);

INSERT INTO `scene_stash_ids_new`
  (
    `scene_id`,
    `endpoint`,
    `stash_id`
  )
  SELECT 
    `scene_id`,
    `endpoint`,
    `stash_id`
  FROM `scene_stash_ids` WHERE
  `scene_id` IS NOT NULL AND `endpoint` IS NOT NULL AND `stash_id` IS NOT NULL;

DROP TABLE `scene_stash_ids`;
ALTER TABLE `scene_stash_ids_new` rename to `scene_stash_ids`;

-- the following index is removed in favour of primary key
-- CREATE INDEX `index_scene_stash_ids_on_scene_id` ON `scene_stash_ids` (`scene_id`);


CREATE TABLE `scenes_galleries_new` (
  `scene_id` integer NOT NULL,
  `gallery_id` integer NOT NULL,
  foreign key(`scene_id`) references `scenes`(`id`) on delete CASCADE,
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE,
  PRIMARY KEY(`scene_id`, `gallery_id`)
);

INSERT INTO `scenes_galleries_new`
  (
    `scene_id`,
    `gallery_id`
  )
  SELECT 
    `scene_id`,
    `gallery_id`
  FROM `scenes_galleries` WHERE 
  `scene_id` IS NOT NULL AND `gallery_id` IS NOT NULL
  ON CONFLICT (`scene_id`, `gallery_id`) DO NOTHING;

DROP TABLE `scenes_galleries`;
ALTER TABLE `scenes_galleries_new` rename to `scenes_galleries`;

CREATE INDEX `index_scenes_galleries_on_gallery_id` on `scenes_galleries` (`gallery_id`);

-- the following index is removed in favour of primary key
-- CREATE INDEX `index_scenes_galleries_on_scene_id` on `scenes_galleries` (`scene_id`);


CREATE TABLE `galleries_images_new` (
  `gallery_id` integer NOT NULL,
  `image_id` integer NOT NULL,
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE,
  foreign key(`image_id`) references `images`(`id`) on delete CASCADE,
  PRIMARY KEY(`gallery_id`, `image_id`)
);

INSERT INTO `galleries_images_new`
  (
    `gallery_id`,
    `image_id`
  )
  SELECT 
    `gallery_id`,
    `image_id`
  FROM `galleries_images` WHERE 
  `image_id` IS NOT NULL AND `gallery_id` IS NOT NULL
  ON CONFLICT (`gallery_id`, `image_id`) DO NOTHING;

DROP TABLE `galleries_images`;
ALTER TABLE `galleries_images_new` rename to `galleries_images`;

CREATE INDEX `index_galleries_images_on_image_id` on `galleries_images` (`image_id`);

-- the following index is removed in favour of primary key
-- CREATE INDEX `index_galleries_images_on_gallery_id` on `galleries_images` (`gallery_id`);


CREATE TABLE `performers_galleries_new` (
  `performer_id` integer NOT NULL,
  `gallery_id` integer NOT NULL,
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE,
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE,
  PRIMARY KEY(`gallery_id`, `performer_id`)
);

INSERT INTO `performers_galleries_new`
  (
    `performer_id`,
    `gallery_id`
  )
  SELECT 
    `performer_id`,
    `gallery_id`
  FROM `performers_galleries` WHERE
  `performer_id` IS NOT NULL AND `gallery_id` IS NOT NULL
  ON CONFLICT (`gallery_id`, `performer_id`) DO NOTHING;

DROP TABLE `performers_galleries`;
ALTER TABLE `performers_galleries_new` rename to `performers_galleries`;

CREATE INDEX `index_performers_galleries_on_performer_id` on `performers_galleries` (`performer_id`);

-- the following index is removed in favour of primary key
-- CREATE INDEX `index_performers_galleries_on_gallery_id` on `performers_galleries` (`gallery_id`);


CREATE TABLE `galleries_tags_new` (
  `gallery_id` integer NOT NULL,
  `tag_id` integer NOT NULL,
  foreign key(`gallery_id`) references `galleries`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE,
  PRIMARY KEY(`gallery_id`, `tag_id`)
);

INSERT INTO `galleries_tags_new`
  (
    `gallery_id`,
    `tag_id`
  )
  SELECT 
    `gallery_id`,
    `tag_id`
  FROM `galleries_tags` WHERE 
  `tag_id` IS NOT NULL AND `gallery_id` IS NOT NULL
  ON CONFLICT (`gallery_id`, `tag_id`) DO NOTHING;

DROP TABLE `galleries_tags`;
ALTER TABLE `galleries_tags_new` rename to `galleries_tags`;

CREATE INDEX `index_galleries_tags_on_tag_id` on `galleries_tags` (`tag_id`);

-- the following index is removed in favour of primary key
-- CREATE INDEX `index_galleries_tags_on_gallery_id` on `galleries_tags` (`gallery_id`);


CREATE TABLE `performers_tags_new` (
  `performer_id` integer NOT NULL,
  `tag_id` integer NOT NULL,
  foreign key(`performer_id`) references `performers`(`id`) on delete CASCADE,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE,
  PRIMARY KEY(`performer_id`, `tag_id`)
);

INSERT INTO `performers_tags_new`
  (
    `performer_id`,
    `tag_id`
  )
  SELECT 
    `performer_id`,
    `tag_id`
  FROM `performers_tags` WHERE true
  ON CONFLICT (`performer_id`, `tag_id`) DO NOTHING;

DROP TABLE `performers_tags`;
ALTER TABLE `performers_tags_new` rename to `performers_tags`;

CREATE INDEX `index_performers_tags_on_tag_id` on `performers_tags` (`tag_id`);

-- the following index is removed in favour of primary key
-- CREATE INDEX `index_performers_tags_on_performer_id` on `performers_tags` (`performer_id`);


CREATE TABLE `tag_aliases_new` (
  `tag_id` integer NOT NULL,
  `alias` varchar(255) NOT NULL,
  foreign key(`tag_id`) references `tags`(`id`) on delete CASCADE,
  PRIMARY KEY(`tag_id`, `alias`)
);

INSERT INTO `tag_aliases_new`
  (
    `tag_id`,
    `alias`
  )
  SELECT 
    `tag_id`,
    `alias`
  FROM `tag_aliases`;

DROP TABLE `tag_aliases`;
ALTER TABLE `tag_aliases_new` rename to `tag_aliases`;

CREATE UNIQUE INDEX `tag_aliases_alias_unique` on `tag_aliases` (`alias`);


CREATE TABLE `studio_aliases_new` (
  `studio_id` integer NOT NULL,
  `alias` varchar(255) NOT NULL,
  foreign key(`studio_id`) references `studios`(`id`) on delete CASCADE,
  PRIMARY KEY(`studio_id`, `alias`)
);

INSERT INTO `studio_aliases_new`
  (
    `studio_id`,
    `alias`
  )
  SELECT 
    `studio_id`,
    `alias`
  FROM `studio_aliases`;

DROP TABLE `studio_aliases`;
ALTER TABLE `studio_aliases_new` rename to `studio_aliases`;

CREATE UNIQUE INDEX `studio_aliases_alias_unique` on `studio_aliases` (`alias`);

PRAGMA foreign_keys=ON;
