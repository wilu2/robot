DROP TABLE IF EXISTS `captchas`;
CREATE TABLE `captchas`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `captcha` varchar(128)  NOT NULL COMMENT '验证码',
  `captcha_id` varchar(32) NOT NULL COMMENT '会话id',
  `created_at` bigint(20) NOT NULL COMMENT '创建时间',
  `expiry` bigint(20) NOT NULL COMMENT '过期时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

alter table `users` add column `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '用户状态 0未使用，1正常，2锁定，3禁用' after `auth_method`;
alter table `users` add column `last_login_time` bigint(20) NOT NULL COMMENT '上次登录时间' after `auth_method`;
alter table `users` add column `last_login_fail_time` bigint(20) NOT NULL COMMENT '上次登录失败时间' after `auth_method`;
alter table `users` add column `failed_logins` int(11) NOT NULL DEFAULT 0 COMMENT '连续登录失败次数' after `auth_method`;
alter table `users` add column `expiry_time` bigint(20) NOT NULL DEFAULT 0 COMMENT '过期时间' after `auth_method`;

UPDATE `users` SET `last_login_time` = unix_timestamp(), `last_login_fail_time` = unix_timestamp(), `expiry_time` = unix_timestamp()+(30 * 24 * 60 * 60);


DROP TABLE IF EXISTS `settings`;
CREATE TABLE `settings`  (
  `key` varchar(128) NOT NULL COMMENT '配置key',
  `setting` int(11) NOT NULL COMMENT '配置值',
  `remark` varchar(128) NOT NULL COMMENT '配置说明'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

INSERT INTO `settings` ( `key`, `setting`,`remark` ) VALUES ( "failed_logins" , 5 , "限制连续登录失败次数" );
INSERT INTO `settings` ( `key`, `setting`,`remark` ) VALUES ( "locked_time" , 5 , "锁定时间，连续输入密码错误后锁定账户时间，单位m" );
INSERT INTO `settings` ( `key`, `setting`,`remark` ) VALUES ( "valid_time" , 30 , "账户默认有效期，单位d" );
INSERT INTO `settings` ( `key`, `setting`,`remark` ) VALUES ( "session_time" , 72 , "会话有效期，单位h" );
INSERT INTO `settings` ( `key`, `setting`,`remark` ) VALUES ( "not_login_time" , 72 , "长时间未登录锁定多久禁用，单位d" );
