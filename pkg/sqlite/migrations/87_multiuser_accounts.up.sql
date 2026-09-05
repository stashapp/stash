CREATE TABLE `users` (
  `id` integer not null primary key autoincrement,
  `username` varchar(255) not null,
  `normalized_username` varchar(255) not null,
  `password_hash` varchar(255),
  `role` varchar(20) not null check (`role` in ('USER', 'MODERATOR', 'ADMIN')),
  `status` varchar(40) not null check (`status` in ('ACTIVE', 'DISABLED', 'PASSWORD_CHANGE_REQUIRED')),
  `created_at` datetime not null,
  `updated_at` datetime not null
);

CREATE UNIQUE INDEX `users_normalized_username_unique` ON `users` (`normalized_username`);

CREATE TABLE `user_sessions` (
  `id` varchar(64) not null primary key,
  `user_id` integer not null references `users` (`id`) on delete cascade,
  `secret_hash` varchar(255) not null,
  `csrf_hash` varchar(255) not null,
  `created_at` datetime not null,
  `last_seen_at` datetime not null,
  `idle_expires_at` datetime not null,
  `absolute_expires_at` datetime not null,
  `revoked_at` datetime,
  `user_agent` text,
  `remote_address` text
);

CREATE UNIQUE INDEX `user_sessions_secret_hash_unique` ON `user_sessions` (`secret_hash`);
CREATE INDEX `user_sessions_user_id` ON `user_sessions` (`user_id`);
CREATE INDEX `user_sessions_expiry` ON `user_sessions` (`idle_expires_at`, `absolute_expires_at`);

CREATE TABLE `user_audit_events` (
  `id` integer not null primary key autoincrement,
  `occurred_at` datetime not null,
  `actor_user_id` integer references `users` (`id`) on delete set null,
  `event_type` varchar(100) not null,
  `target_type` varchar(100),
  `target_id` varchar(255),
  `result` varchar(40) not null,
  `details_json` text
);

CREATE INDEX `user_audit_events_occurred_at` ON `user_audit_events` (`occurred_at`);
CREATE INDEX `user_audit_events_actor` ON `user_audit_events` (`actor_user_id`, `occurred_at`);
