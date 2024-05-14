CREATE TABLE `request_status` (
  `id` int NOT NULL AUTO_INCREMENT,
  `guid` varchar(255) NOT NULL,
  `membership_request_id` int NOT NULL,
  `payment_status` int NOT NULL COMMENT '0: Free, 1; Pending, 2 Success, 3 Failed',
  `status` tinyint NOT NULL COMMENT '1:Pending, 2: Approve, 3: Reject',
  `status_updated` datetime NOT NULL,
  `staff_id` int NOT NULL,
  `outlet_id` int NOT NULL,
  `updated` datetime NOT NULL,
  `created` datetime NOT NULL,
  `active` tinyint NOT NULL DEFAULT '1' COMMENT '0: Inactive, 1:Inactive',
  PRIMARY KEY (`id`),
  KEY `membership_request_id` (`membership_request_id`),
  KEY `payment_status` (`payment_status`),
  KEY `staff_id` (`staff_id`),
  KEY `outlet_id` (`outlet_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
