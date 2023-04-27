/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80027
 Source Host           : localhost:3306
 Source Schema         : regulus_main

 Target Server Type    : MySQL
 Target Server Version : 80027
 File Encoding         : 65001

 Date: 26/04/2023 20:36:24
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for mdao_block_tx_error
-- ----------------------------
DROP TABLE IF EXISTS `mdao_block_tx_error`;
CREATE TABLE `mdao_block_tx_error` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `block` bigint DEFAULT NULL,
  `txhash` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `add_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for mdao_data
-- ----------------------------
DROP TABLE IF EXISTS `mdao_data`;
CREATE TABLE `mdao_data` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `block` bigint DEFAULT NULL,
  `txhash` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `wallet` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `minter` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `term` int DEFAULT NULL,
  `rewards` decimal(64,18) DEFAULT NULL,
  `loss` decimal(64,18) DEFAULT NULL,
  `ts` timestamp NULL DEFAULT NULL,
  `add_time` timestamp NULL DEFAULT NULL,
  `round` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=453 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

SET FOREIGN_KEY_CHECKS = 1;
