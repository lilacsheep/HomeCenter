CREATE TABLE IF NOT EXISTS `object_token` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  `effective` int(11) NOT NULL DEFAULT 0,
  `secret_key` varchar(255) NOT NULL,
  `upload` boolean NOT NULL DEFAULT TRUE,
  `download` boolean NOT NULL DEFAULT TRUE,
  `delete` boolean NOT NULL DEFAULT TRUE,
  `list` boolean NOT NULL DEFAULT TRUE,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `secret_key` (`secret_key`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;