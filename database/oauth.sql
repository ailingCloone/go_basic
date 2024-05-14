CREATE TABLE `oauth` (
  `id` int NOT NULL AUTO_INCREMENT,
  `customer_id` int DEFAULT NULL,
  `staff_id` int DEFAULT NULL,
  `tenant_id` int DEFAULT NULL,
  `access_token` varchar(255) NOT NULL,
  `refresh_token` varchar(255) NOT NULL,
  `expires_in` int NOT NULL,
  `scope` int NOT NULL,
  `token_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `updated` datetime NOT NULL,
  `created` datetime NOT NULL,
  `active` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `customer_id` (`customer_id`),
  KEY `tenant_id` (`tenant_id`),
  KEY `staff_id` (`staff_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci