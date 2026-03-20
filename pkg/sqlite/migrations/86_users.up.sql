CREATE TABLE `users` (
    `id` INTEGER PRIMARY KEY,
    `username` TEXT NOT NULL,
    `api_key` TEXT,
    `password_hash` TEXT,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL
);

CREATE UNIQUE INDEX `users_username_unique` on `users` (`username`);

CREATE TABLE `user_roles` (
    `user_id` INTEGER NOT NULL,
    `role` TEXT NOT NULL,
    PRIMARY KEY (`user_id`, `role`),
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);