
CREATE DATABASE If Not Exists `financial_statement` CHARACTER SET 'utf8mb4';

USE financial_statement;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for formula_title_map
-- ----------------------------
DROP TABLE IF EXISTS `formula_title_map`;
CREATE TABLE `formula_title_map`  (
  `formula_id` int(11) UNSIGNED NOT NULL,
  `title_id` int(11) UNSIGNED NOT NULL,
  INDEX `idx-title_id`(`title_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for login_tokens
-- ----------------------------
DROP TABLE IF EXISTS `login_tokens`;
CREATE TABLE `login_tokens`  (
  `token` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `user_id` int(11) UNSIGNED NOT NULL,
  `created_at` bigint(20) NOT NULL,
  `expiry` bigint(20) NOT NULL,
  PRIMARY KEY (`token`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for logs
-- ----------------------------
DROP TABLE IF EXISTS `logs`;
CREATE TABLE `logs`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `task_id` int(11) UNSIGNED NOT NULL COMMENT 'task id',
  `msg` varchar(3000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` bigint(20) NOT NULL,
  `updated_at` bigint(20) NOT NULL,
  `deleted_at` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for pages
-- ----------------------------
DROP TABLE IF EXISTS `pages`;
CREATE TABLE `pages`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `task_id` int(11) UNSIGNED NOT NULL,
  `file_path` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '文件存储路径',
  `ocr_result` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT 'ocr结果',
  `status` int(11) NOT NULL COMMENT '状态',
  `create_at` bigint(20) NOT NULL,
  `update_at` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `pages_task_id_index`(`task_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for standard_statement_formulas
-- ----------------------------
DROP TABLE IF EXISTS `standard_statement_formulas`;
CREATE TABLE `standard_statement_formulas`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `statement_id` int(11) UNSIGNED NOT NULL,
  `left` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '等号左边公式',
  `right` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '等号右边公式',
  `status` int(11) NOT NULL COMMENT '状态：\r\n1：正常\r\n-1：删除',
  `create_at` bigint(20) NOT NULL,
  `update_at` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx-standard_statement_id`(`statement_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for standard_statement_titles
-- ----------------------------
DROP TABLE IF EXISTS `standard_statement_titles`;
CREATE TABLE `standard_statement_titles`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `statement_id` int(11) UNSIGNED NOT NULL,
  `name` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `external_id` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `aliases` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '别名，存储为json字符串（数据库字段类型为string，暂不用json是担心其他db时兼容性问题）',
  `status` int(11) NOT NULL DEFAULT 1 COMMENT '状态：\r\n1：正常\r\n-1：停用或删除',
  `order_by_id` int(11) NOT NULL COMMENT '排序id，越小越靠前',
  `create_at` bigint(20) NOT NULL,
  `update_at` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx-standard_statement_id`(`statement_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for standard_statements
-- ----------------------------
DROP TABLE IF EXISTS `standard_statements`;
CREATE TABLE `standard_statements`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `standard_id` int(11) UNSIGNED NOT NULL,
  `type` int(11) NOT NULL COMMENT '财务准则报表类型：\r\n1：资产负债表\r\n2：利润表\r\n3：现金流量表',
  `status` int(11) NOT NULL DEFAULT 1 COMMENT '状态：\r\n1：启用\r\n-1：停用',
  `title_status` int(11) NOT NULL DEFAULT 1 COMMENT '科目配置状态：\r\n1：一配置\r\n-1：未配置',
  `formula_status` int(11) NOT NULL DEFAULT 1 COMMENT '试算公式状态：\r\n1：已配置\r\n-1：未配置',
  `create_at` bigint(20) NOT NULL,
  `update_at` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx-standard_id`(`standard_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for standards
-- ----------------------------
DROP TABLE IF EXISTS `standards`;
CREATE TABLE `standards`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `external_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `is_default` int(11) NOT NULL DEFAULT 0 COMMENT '1为默认 -1为非默认',
  `status` int(11) NOT NULL DEFAULT 1 COMMENT '状态：\r\n1：正常\r\n-1：停用',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for tasks
-- ----------------------------
DROP TABLE IF EXISTS `tasks`;
CREATE TABLE `tasks`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `task_name` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `file_format` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `standard_id` int(11) UNSIGNED NOT NULL,
  `external_info` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `async` int(11) NOT NULL,
  `async_status` int(11) NOT NULL DEFAULT 1 COMMENT '1：未同步，10已同步，20同步失败',
  `status` int(11) NOT NULL COMMENT '//已删除\r\nTaskStatusDeleted = -1\r\n//创建任务中\r\nTaskStatusCreateTasking = 1\r\n//任务创建完成\r\nTaskStatusCreated = 10\r\n//识别中\r\nTaskStatusOcring = 20\r\n//识别失败\r\nTaskStatusOcrFailed = 90\r\n//识别成功\r\nTaskStatusOcrSuccess = 100\r\n//已回传\r\nTaskStatusCallbacked = 110',
  `creater_user_id` int(11) UNSIGNED NOT NULL,
  `standard_result` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '财报结构化结果',
  `files` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '[]' COMMENT '该任务对应的原始文件存放路径[\"xxx/xxx.jpg\",\"xxx/xxx/jpg\"]',
  `error` varchar(4096) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `created_at` bigint(20) NOT NULL,
  `updated_at` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `account` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `password` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `salt` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `email` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `mobile` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `auth_method` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` bigint(20) NOT NULL,
  `updated_at` bigint(20) NOT NULL,
  `deleted_at` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `users_account_uindex`(`account`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ALTER TABLE `tasks` 
-- MODIFY COLUMN `external_info` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL AFTER `standard_id`;

ALTER TABLE `tasks` 
MODIFY COLUMN `files` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '该任务对应的原始文件存放路径[\"xxx/xxx.jpg\",\"xxx/xxx/jpg\"]' AFTER `standard_result`;

INSERT INTO `financial_statement`.`users`(`name`, `account`, `password`, `salt`, `email`, `mobile`, `auth_method`, `created_at`, `updated_at`, `deleted_at`) VALUES ('admin', 'admin', '7M9t72i6fX05YysJbEr/uyvb5lEBHXbPdVTjc4AM4m0=', 'V7feAtVx', '', '', '', 1663835299, 1663835299, NULL);
INSERT INTO `financial_statement`.`standards`( `name`, `external_id`, `is_default`, `status`) VALUES ('新会计准则财报', '', 1, 1);
INSERT INTO `financial_statement`.`standard_statements`( `standard_id`, `type`, `status`, `title_status`, `formula_status`, `create_at`, `update_at`) VALUES ( 1, 1, 1, -1, -1, 1666321480, 1666321480);
INSERT INTO `financial_statement`.`standard_statements`( `standard_id`, `type`, `status`, `title_status`, `formula_status`, `create_at`, `update_at`) VALUES ( 1, 2, 1, -1, -1, 1666321480, 1666321480);
INSERT INTO `financial_statement`.`standard_statements`( `standard_id`, `type`, `status`, `title_status`, `formula_status`, `create_at`, `update_at`) VALUES ( 1, 3, 1, -1, -1, 1666321480, 1666321480);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '货币资金', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '存货', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '流动资产', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '交易性金融资产', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '应收票据', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '应收账款', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '预付款项', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '应收股利', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '应收利息', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其他应收款', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其中：消耗性生物资产', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '待摊费用', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '一年内到期的非流动资产', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其他流动资产', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '流动资产合计', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '非流动资产', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '可供出售金融资产', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '持有至到期投资', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '投资性房地产', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '长期股权投资', '', '', 1, 9999, 1666321721, 1666321721);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '长期应收款', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '固定资产原价', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '在建工程', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '工程物资', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '固定资产清理', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '生产性生物资产', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '油气资产', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '无形资产', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '开发支出', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '商誉', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '长期待摊费用', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '递延所得税资产', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其他非流动资产', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '非流动资产合计', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '资产总计', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '流动负债', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '短期借款', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '交易性金融负债', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '应付票据', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '应付账款', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '预收款项', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '应付职工薪酬', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '应交税费', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '应付利息', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '应付股利', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其他应付款', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '预提费用', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '预计负债', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '一年内到期的非流动负债', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其他流动负债', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '流动负债合计', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '非流动负债', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '长期借款', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '应付债券', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '长期应付款', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '专项应付款', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '递延所得税负债', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其他非流动负债', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '非流动负债合计', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '负债合计', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '所有者权益', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '实收资本（股本）', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '资本公积', '', '', 1, 9999, 1666322410, 1666322410);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '盈余公积', '', '', 1, 9999, 1666322633, 1666322633);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '未分配利润', '', '', 1, 9999, 1666322633, 1666322633);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '减：库存股', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '所有者权益合计', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '负债和所有者权益合计', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '结算备付金', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '拆出资金', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '以公允价值计量且其变动计入当期损益的金融资产', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '衍生金融资产', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '买入返售金融资产', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '合同资产', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '持有待售资产', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '发放贷款和垫款', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '债权投资', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其他债权投资', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其他权益工具投资', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其他非流动金融资产', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '使用权资产', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '拆入资金', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '以公允价值计量且其变动计入当期损益的金融负债', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '衍生金融负债', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '合同负债', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '卖出回购金融资产款', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '吸收存款及同业存放', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其他综合收益', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '专项储备', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '一般风险准备', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '归属于母公司股东权益合计', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '少数股东权益', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其中：优先股', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '向中央银行借款', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '持有待售负债', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '递延收益', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其他权益工具', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其中：原材料', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '库存商品（产成品)', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '减：累计折旧', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '减：固定资产减值准备', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其中：土地使用权', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其中：特准储备物资', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '应付权证', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '其中：特种储备基金', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '未确认的投资损失（以\\-\\号填列）', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '外币报表折算差额', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '固定资产净值', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '固定资产净额', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '永续债', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '应收款项融资', '', '', 1, 9999, 1666322634, 1666322634);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '营业收入', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '减：营业成本', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '税金及附加', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '销售费用', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '管理费用', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '财务费用', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '加：公允价值变动收益', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '其中：对联营企业和合营企业的投资收益', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '其他收益', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '营业利润（亏损以“-”号填列）', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '加：营业外收入', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '其中：非流动资产处置利得', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '减：营业外支出', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '其中：非流动资产处置损失', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '利润总额（亏损总额以“-”号填列）', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '减：所得税费用', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '净利润（净亏损以“-”号填列）', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '持续经营净利润（净亏损以“-”号填列）', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '终止经营净利润（净亏损以“-”号填列）', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '其他综合收益的税后净额', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '不能重分类进损益的其他综合收益', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '重新计量设定受益计划净负债或净资产的变动', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '权益法下在被投资单位不能重分类进损益的其他综合收益中享有的份额', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '将重分类进损益的其他综合收益', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '权益法下在被投资单位以后将重分类进损益的其他综合收益中享有的份额', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '可供出售金融资产公允价值变动损益', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '持有至到期投资重分类为可供出售金融资产损益', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '现金流经套期损益的有效部分', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '外币财务报表折算差额', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '其他', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '综合收益总额', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '每股收益', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '基本每股收益', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '稀释每股收益', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '加：年初未分配利润', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '其他转入', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '减：提取法宝盈余公积', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '提取企业储备基金', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '提取企业发展基金', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '提取职工奖励及福利基金', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '利润归还投资', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '应付优先股股利', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '提取任意盈余公积', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '应付普通股股利', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '转作资本（或股本）的普通股股利', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '转总部利润', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '未分配利润', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '其中：主营业务收入', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '其他业务收入', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '营业总收入', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '研发费用', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '其中：利息费用', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '利息收入', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '加：其他收益', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '投资收益（损失以“-”号填列）', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '以摊余成本计量的金融资产终止确认收益', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '净敝口套期收益（损失以“-”号填列）', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '公允价值变动收益（损失以“-”号填列）', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '信用减值损失（损失以“-”号填列）', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '资产减值损失（损失以“-”号填列）', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '资产处置收益（损失以“-”号填列）', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '归属于母公司所有者的净利润（净亏损以“-”号填列））', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '少数股东损益（净亏损以“-”号填列）', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '其他业务成本', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '其中：主营业务成本', '', '', 1, 9999, 1666322924, 1666322924);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '经营活动产生的现金流量', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '销售商品、提供劳务收到的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '收到的税费返还', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '收到其他与经营活动有关的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '经营活动现金流入小计', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '购买商品、接受劳务支付的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '支付给职工以及为职工支付的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '支付的各项税费', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '支付其他与经营活动有关的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '经营活动现金流出小计', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '经营活动产生的现金流量净额', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '投资活动产生的现金流量', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '收回投资收到的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '取得投资收益收到的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '处置固定资产、无形资产和其他长期资产收回的现金净额', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '处置子公司及其他营业单位收到的现金净额', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '收到其他与投资活动有关的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '投资活动现金流入小计', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '购建固定资产、无形资产和其他长期资产支付的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '投资支付的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '取得子公司及其他营业单位支付的现金净额', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '支付其他与投资活动有关的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '投资活动现金流出小计', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '投资活动产生的现金流量净额', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '筹资活动产生的现金流量', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '吸收投资收到的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '取得借款收到的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '收到其他与筹资活动有关的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '筹资活动现金流入小计', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '偿还债务支付的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '分配股利、利润或偿付利息支付的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '支付其他与筹资活动有关的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '筹资活动现金流出小计', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '筹资活动产生的现金流量净额', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '汇率变动对现金及现金等价物的影响', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '现金及现金等价物净增加额', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '期初现金及现金等价物余额', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '期末现金及现金等价物余额', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '其中：子公司吸收少数股东投资收到的现金', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '其中：子公司支付给少数股东的股利、利润', '', '', 1, 9999, 1666323045, 1666323045);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '货币资金', '', '', 1, 9999, 1666341420, 1666341420);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '货币资金', '', '', 1, 9999, 1666341457, 1666341457);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '营业收入', '', '', 1, 9999, 1666341849, 1666341849);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (1, '营业收入', '', '', 1, 9999, 1666341859, 1666341859);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (3, '营业收入', '', '', 1, 9999, 1666341868, 1666341868);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '其中：营业收入', '', '', 1, 9999, 1666679103, 1666679103);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '利息收入', '', '', 1, 9999, 1666679133, 1666679133);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '已赚保费', '', '', 1, 9999, 1666679133, 1666679133);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '手续费及佣金收入', '', '', 1, 9999, 1666679133, 1666679133);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '税金及附加', '', '', 1, 9999, 1666700042, 1666700042);
INSERT INTO `financial_statement`.`standard_statement_titles`( `statement_id`, `name`, `external_id`, `aliases`, `status`, `order_by_id`, `create_at`, `update_at`) VALUES (2, '销售费用', '', '', 1, 9999, 1666700042, 1666700042);



SET FOREIGN_KEY_CHECKS = 1;
