CREATE TABLE `users` (
  `id` varchar(36) NOT NULL,
  `user_name` varchar(45) NOT NULL,
  `password` varchar(45) NOT NULL,
  `real_name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `rank` varchar(45) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_name_UNIQUE` (`user_name`),
  UNIQUE KEY `email_UNIQUE` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `tracks` (
  `id` varchar(36) NOT NULL,
  `name` varchar(45) NOT NULL,
  `short_name` varchar(10) NOT NULL,
  `type` varchar(45) NOT NULL,
  `logo` varchar(45) NOT NULL,
  `description` varchar(255) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `track_name_UNIQUE` (`name`),
  UNIQUE KEY `short_track_name_UNIQUE` (`short_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `cars` (
  `id` varchar(36) NOT NULL,
  `name` varchar(45) NOT NULL,
  `short_name` varchar(10) NOT NULL,
  `logo` varchar(45) NOT NULL,
  `type` varchar(45) NOT NULL,
  `num_seats` int NOT NULL,
  `description` varchar(45) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `track` varchar(36) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`),
  UNIQUE KEY `short_name_UNIQUE` (`short_name`),
  KEY `idx_cars_for_track` (`track`,`type`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


ALTER TABLE `gorrc`.`cars` 
ADD INDEX `idx_cars_for_track` (`track` ASC, `type` ASC, `name` ASC) VISIBLE;

CREATE TABLE `track_x_owner` (
  `track` varchar(36) NOT NULL,
  `owner` varchar(36) NOT NULL,
  PRIMARY KEY (`track`,`owner`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
