CREATE TABLE `playback_defaults` (
  `id` integer not null primary key autoincrement,
  `user_agent_pattern` varchar(255) NOT NULL,
  `priority` integer NOT NULL DEFAULT 100,
  `stream_type` varchar(50) NOT NULL,
  `quality` varchar(50),
  `created_at` datetime not null,
  `updated_at` datetime not null
);

CREATE UNIQUE INDEX `playback_defaults_pattern_unique` ON `playback_defaults` (`user_agent_pattern`);
