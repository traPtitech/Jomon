-- Create "user_subjects" table
CREATE TABLE `user_subjects` (`subject` varchar(36) NOT NULL, `user_id` uuid NOT NULL, PRIMARY KEY (`subject`), INDEX `user_subjects_users_user` (`user_id`), CONSTRAINT `user_subjects_users_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE RESTRICT ON DELETE NO ACTION) CHARSET utf8mb4 COLLATE utf8mb4_bin;
