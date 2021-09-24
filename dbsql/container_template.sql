CREATE TABLE IF NOT EXISTS `docker_container_option` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `key` varchar(255) NOT NULL DEFAULT '',
  `value` TEXT,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;