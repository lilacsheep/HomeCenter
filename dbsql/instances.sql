CREATE TABLE IF NOT EXISTS `instances` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `protocol` int(11) NOT NULL DEFAULT 0,
  `address` varchar(64) NOT NULL DEFAULT '',
  `username` varchar(64) NOT NULL DEFAULT '',
  `password` varchar(255) NOT NULL DEFAULT '',
  `private_key` TEXT,
  `status` boolean NOT NULL DEFAULT TRUE,
  `force_country` boolean NOT NULL DEFAULT FALSE,
  `country_code` varchar(32) NOT NULL DEFAULT '',
  `country` varchar(32) NOT NULL DEFAULT '',
  `delay` int(11) NOT NULL DEFAULT 0,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `address` (`address`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4