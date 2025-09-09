-- MySQL 8 schema for keti3
-- Database: kt3 (created by docker-compose)
-- Charset/Collation: utf8mb4 / utf8mb4_0900_ai_ci

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 1) Student Info
CREATE TABLE IF NOT EXISTS `keti3_student_info` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `uid` BIGINT UNSIGNED NOT NULL COMMENT 'logical uid',
  `school_id` BIGINT UNSIGNED NOT NULL,
  `student_name` VARCHAR(64) NOT NULL,
  `student_num` VARCHAR(64) NOT NULL,
  `grade_id` BIGINT UNSIGNED NOT NULL,
  `class_name` VARCHAR(64) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_student_unique` (`school_id`, `student_num`, `grade_id`, `class_name`),
  KEY `idx_uid` (`uid`),
  KEY `idx_school_grade_class` (`school_id`, `grade_id`, `class_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 2) Student Submit
CREATE TABLE IF NOT EXISTS `keti3_student_submit` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `submit_id` BIGINT UNSIGNED DEFAULT NULL COMMENT 'group/attempt id if applicable',
  `uid` BIGINT UNSIGNED NOT NULL,
  `school_id` BIGINT UNSIGNED NOT NULL,
  `student_name` VARCHAR(64) NOT NULL,
  `student_num` VARCHAR(64) NOT NULL,
  `grade_id` BIGINT UNSIGNED NOT NULL,
  `class_name` VARCHAR(64) NOT NULL,
  `voice_url` VARCHAR(255) DEFAULT NULL,
  `oplog_url` VARCHAR(255) DEFAULT NULL,
  `screenshot_url` VARCHAR(255) DEFAULT NULL,
  `date` DATE DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_uid_date` (`uid`, `date`),
  KEY `idx_school_grade_class_date` (`school_id`, `grade_id`, `class_name`, `date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 3) Operation Log
CREATE TABLE IF NOT EXISTS `keti3_op_log` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `submit_id` BIGINT UNSIGNED DEFAULT NULL,
  `uid` BIGINT UNSIGNED NOT NULL,
  `op_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `op_type` VARCHAR(32) NOT NULL,
  `op_object` VARCHAR(64) DEFAULT NULL,
  `object_no` VARCHAR(64) DEFAULT NULL,
  `object_name` VARCHAR(128) DEFAULT NULL,
  `data_before` JSON DEFAULT NULL,
  `data_after` JSON DEFAULT NULL,
  `voice_url` VARCHAR(255) DEFAULT NULL,
  `screenshot_url` VARCHAR(255) DEFAULT NULL,
  `deleted_at` TIMESTAMP NULL DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_uid_time` (`uid`, `op_time`),
  KEY `idx_submit_id` (`submit_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;
