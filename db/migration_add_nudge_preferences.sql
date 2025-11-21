-- Migration: Add nudge preferences to travellers table
-- Run this in TiDB Cloud Chat2Query

USE test;

ALTER TABLE `travellers` 
ADD COLUMN `nudge_frequency_minutes` INT DEFAULT NULL COMMENT 'NULL = conditional nudges, otherwise fixed frequency in minutes',
ADD COLUMN `nudges_paused_until` DATETIME DEFAULT NULL COMMENT 'NULL = not paused, otherwise pause until this time',
ADD COLUMN `nudges_disabled` TINYINT(1) DEFAULT 0 COMMENT '1 = paused indefinitely, 0 = active';

