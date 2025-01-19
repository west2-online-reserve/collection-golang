CREATE TABLE casaos.`user`
(
    `id` bigint NOT NULL AUTO_INCREMENT,
    `username` varchar(255) NOT NULL,
    `password` varchar(255) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT current_timestamp,
    `updated_at` timestamp NOT NULL ON UPDATE current_timestamp DEFAULT current_timestamp,
    `deleted_at` timestamp NULL DEFAULT NULL,
    CONSTRAINT `id` PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4;

CREATE TABLE casaos.`task`
(
    `id` bigint NOT NULL AUTO_INCREMENT,
    `user_id` bigint NOT NULL,
    `title` varchar(255) NOT NULL,
    `content` varchar(255) NOT NULL,
    `status` bigint NOT NULL,
    `start_at` timestamp NOT NULL,
    `end_at` timestamp NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT current_timestamp,
    `updated_at` timestamp NOT NULL ON UPDATE current_timestamp DEFAULT current_timestamp,
    `deleted_at` timestamp NULL DEFAULT NULL,
    CONSTRAINT `id` PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4;
