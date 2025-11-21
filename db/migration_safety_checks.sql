-- Migration: Add safety check polling and emergency escalation
-- Run this in TiDB Cloud Chat2Query

USE test;

-- Add emergency contact fields to travellers
ALTER TABLE `travellers` 
ADD COLUMN `emergency_contact_name` VARCHAR(255) DEFAULT NULL,
ADD COLUMN `emergency_contact_phone` VARCHAR(40) DEFAULT NULL,
ADD COLUMN `hotel_whatsapp_number` VARCHAR(40) DEFAULT NULL,
ADD COLUMN `last_nudge_at` DATETIME DEFAULT NULL COMMENT 'Last time a safety check nudge was sent';

-- Table to track safety check responses
CREATE TABLE IF NOT EXISTS `safety_check_responses` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `traveller_id` BIGINT UNSIGNED NOT NULL,
  `nudge_id` VARCHAR(100) DEFAULT NULL COMMENT 'Identifier for the nudge being responded to',
  `response` VARCHAR(20) NOT NULL COMMENT 'yes_safe or no_trouble',
  `response_time` DATETIME NOT NULL,
  `response_lat` DECIMAL(10,6) DEFAULT NULL,
  `response_lng` DECIMAL(10,6) DEFAULT NULL,
  `confirmed_escalation` TINYINT(1) DEFAULT 0 COMMENT 'If user confirmed emergency escalation',
  PRIMARY KEY (`id`),
  KEY `idx_traveller_response` (`traveller_id`, `response_time`),
  CONSTRAINT `fk_safety_traveller` FOREIGN KEY (`traveller_id`) REFERENCES `travellers` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table to track emergency escalations
CREATE TABLE IF NOT EXISTS `emergency_escalations` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `traveller_id` BIGINT UNSIGNED NOT NULL,
  `response_id` BIGINT UNSIGNED DEFAULT NULL COMMENT 'Related safety_check_response',
  `escalated_at` DATETIME NOT NULL,
  `hotel_notified` TINYINT(1) DEFAULT 0,
  `emergency_contact_notified` TINYINT(1) DEFAULT 0,
  `police_notified` TINYINT(1) DEFAULT 0,
  `location_lat` DECIMAL(10,6) NOT NULL,
  `location_lng` DECIMAL(10,6) NOT NULL,
  `location_address` VARCHAR(500) DEFAULT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT 'active, resolved, cancelled',
  `notes` TEXT,
  PRIMARY KEY (`id`),
  KEY `idx_escalation_traveller` (`traveller_id`, `status`),
  CONSTRAINT `fk_escalation_traveller` FOREIGN KEY (`traveller_id`) REFERENCES `travellers` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

