CREATE TABLE `tenant` (
  `id` int NOT NULL AUTO_INCREMENT,
  `guid` varchar(255) NOT NULL,
  `client_id` varchar(255) NOT NULL,
  `client_secret` varchar(255) NOT NULL,
  `title` varchar(255) NOT NULL,
  `updated` datetime NOT NULL,
  `created` datetime NOT NULL,
  `active` int NOT NULL COMMENT '0-Inactive, 1- Active',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci