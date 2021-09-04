CREATE TABLE IF NOT EXISTS `ddns_operation_settings` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `provider` varchar(255) NOT NULL DEFAULT '',
  `domain` varchar(255) NOT NULL DEFAULT '',
  `sub_domain` varchar(255) NOT NULL DEFAULT '',
  `time_interval` varchar(255) NOT NULL DEFAULT '',
  `use_public_ip` boolean NOT NULL DEFAULT TRUE,
  `net_card` int(11) NOT NULL,
  `record_id` varchar(255) NOT NULL DEFAULT '',
  `dnspod_id` varchar(255) NOT NULL DEFAULT '',
  `dnspod_token` varchar(255) NOT NULL DEFAULT '',
  `updated_on` datetime,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'update time',
  PRIMARY KEY (`id`),
  KEY `sub_domain` (`sub_domain`, `domain`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;