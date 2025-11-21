CREATE DATABASE IF NOT EXISTS `go-mysql-crud`;
USE `go-mysql-crud`;

CREATE TABLE IF NOT EXISTS `posts` (
  `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(255) NOT NULL,
  `content` TEXT NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `posts` (`title`, `content`)
SELECT 'Hello world', 'This is the seed post.'
WHERE NOT EXISTS (SELECT 1 FROM `posts` LIMIT 1);

CREATE TABLE IF NOT EXISTS `test` (
  `id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `age` INT NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `test` (`name`, `age`)
SELECT * FROM (
  SELECT 'Alice' AS `name`, 30 AS `age`
  UNION ALL
  SELECT 'Bob', 27
) AS seeds
WHERE NOT EXISTS (SELECT 1 FROM `test`);

CREATE TABLE IF NOT EXISTS `travellers` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `phone` VARCHAR(40) NOT NULL,
  `hotel_name` VARCHAR(255) NOT NULL,
  `hotel_address` VARCHAR(255) NOT NULL,
  `hotel_language_prompt` VARCHAR(255) NOT NULL,
  `location_permission` TINYINT(1) DEFAULT 0,
  `last_lat` DECIMAL(10,6) DEFAULT NULL,
  `last_lng` DECIMAL(10,6) DEFAULT NULL,
  `last_location_source` VARCHAR(50) DEFAULT NULL,
  `last_location_at` DATETIME DEFAULT NULL,
  `network_lost` TINYINT(1) DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `sos_events` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `traveller_id` BIGINT UNSIGNED NOT NULL,
  `channel` VARCHAR(50) NOT NULL,
  `notes` TEXT,
  `status` VARCHAR(20) NOT NULL DEFAULT 'active',
  `triggered_at` DATETIME NOT NULL,
  `resolved_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_sos_traveller` (`traveller_id`),
  CONSTRAINT `fk_sos_traveller` FOREIGN KEY (`traveller_id`) REFERENCES `travellers` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `hazard_zones` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `latitude` DECIMAL(10,6) NOT NULL,
  `longitude` DECIMAL(10,6) NOT NULL,
  `radius_km` DECIMAL(5,2) NOT NULL,
  `severity` VARCHAR(20) NOT NULL,
  `safe_rating` DECIMAL(2,1) NOT NULL DEFAULT 5.0,
  `description` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `weather_alerts` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `latitude` DECIMAL(10,6) NOT NULL,
  `longitude` DECIMAL(10,6) NOT NULL,
  `radius_km` DECIMAL(5,2) NOT NULL,
  `type` VARCHAR(50) NOT NULL,
  `description` VARCHAR(255) NOT NULL,
  `severity` VARCHAR(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `solo_guides` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `language` VARCHAR(50) NOT NULL,
  `contact` VARCHAR(100) NOT NULL,
  `rating` DECIMAL(2,1) NOT NULL,
  `rate_per_hr` DECIMAL(6,2) NOT NULL,
  `specialty` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `hazard_zones` (`name`, `latitude`, `longitude`, `radius_km`, `severity`, `safe_rating`, `description`)
SELECT * FROM (
  SELECT 'Market backlane', 28.612000, 77.230000, 0.25, 'high', 3.2, 'Avoid the poorly lit backlane near the market after 9pm. Take the main road.' UNION ALL
  SELECT 'Closed metro construction', 28.610500, 77.235500, 0.30, 'medium', 3.5, 'Construction zone ahead. Use the northern footpath for better lighting.'
) AS hazard_seed
WHERE NOT EXISTS (SELECT 1 FROM `hazard_zones`);

INSERT INTO `weather_alerts` (`latitude`, `longitude`, `radius_km`, `type`, `description`, `severity`)
SELECT * FROM (
  SELECT 28.613900, 77.209000, 5.00, 'rain', 'Heavy rain predicted in 45 mins. Carry a jacket or wait indoors.', 'medium'
) AS weather_seed
WHERE NOT EXISTS (SELECT 1 FROM `weather_alerts`);

INSERT INTO `solo_guides` (`name`, `language`, `contact`, `rating`, `rate_per_hr`, `specialty`)
SELECT * FROM (
  SELECT 'Anika Rao', 'English/Hindi', '+91-90000-11111', 4.9, 20.00, 'Hidden cafes & street food' UNION ALL
  SELECT 'Felipe Morales', 'Spanish/English', '+34-600-200-300', 4.7, 18.50, 'Night heritage walks'
) AS guide_seed
WHERE NOT EXISTS (SELECT 1 FROM `solo_guides`);

