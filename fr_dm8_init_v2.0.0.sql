CREATE TABLE "SYSDBA"."users"
(
 "id" BIGINT IDENTITY(8,1) NOT NULL,
 "name" VARCHAR(255) NOT NULL,
 "account" VARCHAR(255) NOT NULL,
 "password" VARCHAR(64) NOT NULL,
 "salt" VARCHAR(32) NOT NULL,
 "email" VARCHAR(128) NOT NULL,
 "mobile" VARCHAR(32) NOT NULL,
 "auth_method" VARCHAR(128) NOT NULL,
 "created_at" BIGINT NOT NULL,
 "updated_at" BIGINT NOT NULL,
 "deleted_at" BIGINT NULL
);
CREATE TABLE "SYSDBA"."tasks"
(
 "id" BIGINT IDENTITY(367,1) NOT NULL,
 "task_name" VARCHAR(1024) NOT NULL,
 "file_format" VARCHAR(16) NOT NULL,
 "standard_id" BIGINT NOT NULL,
 "external_info" CLOB NOT NULL,
 "async" INT NOT NULL,
 "async_status" INT DEFAULT 1
 NOT NULL,
 "status" INT NOT NULL,
 "creater_user_id" BIGINT NOT NULL,
 "standard_result" CLOB NULL,
 "files" VARCHAR(1000) DEFAULT '[]'
 NOT NULL,
 "error" VARCHAR(4096) NULL,
 "created_at" BIGINT NOT NULL,
 "updated_at" BIGINT NOT NULL
);
CREATE TABLE "SYSDBA"."standards"
(
 "id" BIGINT IDENTITY(22,1) NOT NULL,
 "name" VARCHAR(255) NOT NULL,
 "external_id" VARCHAR(255) NOT NULL,
 "is_default" INT DEFAULT 0
 NOT NULL,
 "status" INT DEFAULT 1
 NOT NULL
);
CREATE TABLE "SYSDBA"."standard_statements"
(
 "id" BIGINT IDENTITY(64,1) NOT NULL,
 "standard_id" BIGINT NOT NULL,
 "type" INT NOT NULL,
 "status" INT DEFAULT 1
 NOT NULL,
 "title_status" INT DEFAULT 1
 NOT NULL,
 "formula_status" INT DEFAULT 1
 NOT NULL,
 "create_at" BIGINT NOT NULL,
 "update_at" BIGINT NOT NULL
);
CREATE TABLE "SYSDBA"."standard_statement_titles"
(
 "id" BIGINT IDENTITY(258,1) NOT NULL,
 "statement_id" BIGINT NOT NULL,
 "name" VARCHAR(1024) NOT NULL,
 "external_id" VARCHAR(1024) NULL,
 "aliases" VARCHAR(1024) NULL,
 "status" INT DEFAULT 1
 NOT NULL,
 "order_by_id" INT NOT NULL,
 "create_at" BIGINT NOT NULL,
 "update_at" BIGINT NOT NULL
);
CREATE TABLE "SYSDBA"."standard_statement_formulas"
(
 "id" BIGINT IDENTITY(32,1) NOT NULL,
 "statement_id" BIGINT NOT NULL,
 "left" VARCHAR(2048) NOT NULL,
 "right" VARCHAR(2048) NOT NULL,
 "status" INT NOT NULL,
 "create_at" BIGINT NOT NULL,
 "update_at" BIGINT NOT NULL
);
CREATE TABLE "SYSDBA"."pages"
(
 "id" BIGINT IDENTITY(1685,1) NOT NULL,
 "task_id" BIGINT NOT NULL,
 "file_path" VARCHAR(1024) NOT NULL,
 "ocr_result" CLOB NULL,
 "status" INT NOT NULL,
 "create_at" BIGINT NOT NULL,
 "update_at" BIGINT NOT NULL
);
CREATE TABLE "SYSDBA"."logs"
(
 "id" BIGINT IDENTITY(507,1) NOT NULL,
 "task_id" BIGINT NOT NULL,
 "msg" VARCHAR(3000) NOT NULL,
 "created_at" BIGINT NOT NULL,
 "updated_at" BIGINT NOT NULL,
 "deleted_at" BIGINT NULL
);
CREATE TABLE "SYSDBA"."login_tokens"
(
 "token" VARCHAR(128) NOT NULL,
 "user_id" BIGINT NOT NULL,
 "created_at" BIGINT NOT NULL,
 "expiry" BIGINT NOT NULL
);
CREATE TABLE "SYSDBA"."formula_title_map"
(
 "formula_id" BIGINT NOT NULL,
 "title_id" BIGINT NOT NULL
);
ALTER TABLE "SYSDBA"."users" ADD CONSTRAINT  PRIMARY KEY("id") ;

ALTER TABLE "SYSDBA"."users" ADD CONSTRAINT "users_account_uindex" UNIQUE("account") ;

ALTER TABLE "SYSDBA"."tasks" ADD CONSTRAINT  PRIMARY KEY("id") ;

ALTER TABLE "SYSDBA"."standards" ADD CONSTRAINT  PRIMARY KEY("id") ;

ALTER TABLE "SYSDBA"."standard_statements" ADD CONSTRAINT  PRIMARY KEY("id") ;

ALTER TABLE "SYSDBA"."standard_statement_titles" ADD CONSTRAINT  PRIMARY KEY("id") ;

ALTER TABLE "SYSDBA"."standard_statement_formulas" ADD CONSTRAINT  PRIMARY KEY("id") ;

ALTER TABLE "SYSDBA"."pages" ADD CONSTRAINT  PRIMARY KEY("id") ;

ALTER TABLE "SYSDBA"."logs" ADD CONSTRAINT  PRIMARY KEY("id") ;

ALTER TABLE "SYSDBA"."login_tokens" ADD CONSTRAINT  PRIMARY KEY("token") ;

ALTER TABLE "SYSDBA"."tasks" ADD CHECK("standard_id" >= 0) ENABLE ;

ALTER TABLE "SYSDBA"."tasks" ADD CHECK("creater_user_id" >= 0) ENABLE ;

COMMENT ON COLUMN "SYSDBA"."tasks"."async_status" IS '1：未同步，10已同步，20同步失败';

COMMENT ON COLUMN "SYSDBA"."tasks"."status" IS '//已删除
TaskStatusDeleted = -1
//创建任务中
TaskStatusCreateTasking = 1
//任务创建完成
TaskStatusCreated = 10
//识别中
TaskStatusOcring = 20
//识别失败
TaskStatusOcrFailed = 90
//识别成功
TaskStatusOcrSuccess = 100
//已回传
TaskStatusCallbacked = 110';

COMMENT ON COLUMN "SYSDBA"."tasks"."standard_result" IS '财报结构化结果';

COMMENT ON COLUMN "SYSDBA"."tasks"."files" IS '该任务对应的原始文件存放路径["xxx/xxx.jpg","xxx/xxx/jpg"]';

COMMENT ON COLUMN "SYSDBA"."standards"."is_default" IS '1为默认 -1为非默认';

COMMENT ON COLUMN "SYSDBA"."standards"."status" IS '状态：
1：正常
-1：停 用状态 ';

ALTER TABLE "SYSDBA"."standard_statements" ADD CHECK("standard_id" >= 0) ENABLE ;

CREATE INDEX "idx-standard_id"
ON "SYSDBA"."standard_statements"("standard_id");

COMMENT ON COLUMN "SYSDBA"."standard_statements"."type" IS '财务准则报表类型：
1：资产负债表
2：利润表
3：现金流量表';

COMMENT ON COLUMN "SYSDBA"."standard_statements"."status" IS '状态：
1：启用
-1：停用';

COMMENT ON COLUMN "SYSDBA"."standard_statements"."title_status" IS '科目配置状态：
1：一配置
-1：未配置';

COMMENT ON COLUMN "SYSDBA"."standard_statements"."formula_status" IS '试算公式状态：
1：已配置
-1：未配置';

ALTER TABLE "SYSDBA"."standard_statement_titles" ADD CHECK("statement_id" >= 0) ENABLE ;

CREATE INDEX "idx-standard_statement_id"
ON "SYSDBA"."standard_statement_titles"("statement_id");

COMMENT ON COLUMN "SYSDBA"."standard_statement_titles"."aliases" IS '别名，存储为json字符串（数据库字段类型为string，暂不用json是担心其他db时兼容性问题）';

COMMENT ON COLUMN "SYSDBA"."standard_statement_titles"."status" IS '状态：
1：正常
-1：停用或删除';

COMMENT ON COLUMN "SYSDBA"."standard_statement_titles"."order_by_id" IS '排序id，越小越靠前';

ALTER TABLE "SYSDBA"."standard_statement_formulas" ADD CHECK("statement_id" >= 0) ENABLE ;

COMMENT ON COLUMN "SYSDBA"."standard_statement_formulas"."left" IS '等号左边公式';

COMMENT ON COLUMN "SYSDBA"."standard_statement_formulas"."right" IS '等号右边公式';

COMMENT ON COLUMN "SYSDBA"."standard_statement_formulas"."status" IS '状态：
1：正常
-1：删除';

ALTER TABLE "SYSDBA"."pages" ADD CHECK("task_id" >= 0) ENABLE ;

CREATE INDEX "pages_task_id_index"
ON "SYSDBA"."pages"("task_id");

COMMENT ON COLUMN "SYSDBA"."pages"."file_path" IS '文件存储路径';

COMMENT ON COLUMN "SYSDBA"."pages"."ocr_result" IS 'ocr结果';

COMMENT ON COLUMN "SYSDBA"."pages"."status" IS '状态';

ALTER TABLE "SYSDBA"."logs" ADD CHECK("task_id" >= 0) ENABLE ;

COMMENT ON COLUMN "SYSDBA"."logs"."task_id" IS 'task id';

ALTER TABLE "SYSDBA"."login_tokens" ADD CHECK("user_id" >= 0) ENABLE ;

ALTER TABLE "SYSDBA"."formula_title_map" ADD CHECK("formula_id" >= 0) ENABLE ;

ALTER TABLE "SYSDBA"."formula_title_map" ADD CHECK("title_id" >= 0) ENABLE ;

CREATE INDEX "idx-title_id"
ON "SYSDBA"."formula_title_map"("title_id");




SET IDENTITY_INSERT "SYSDBA"."logs" ON;
SET IDENTITY_INSERT "SYSDBA"."logs" OFF;
SET IDENTITY_INSERT "SYSDBA"."pages" ON;
SET IDENTITY_INSERT "SYSDBA"."pages" OFF;
SET IDENTITY_INSERT "SYSDBA"."standard_statement_formulas" ON;
SET IDENTITY_INSERT "SYSDBA"."standard_statement_formulas" OFF;
SET IDENTITY_INSERT "SYSDBA"."standard_statement_titles" ON;
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(1,52,'ceshi','','',1,3,1666002875,1666002875);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(2,52,'ceshi2','','',1,1,1666003122,1666058532);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(3,52,'ceshi6666','','["111","22222"]',-1,2,1666007254,1666059502);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(4,52,'testss','','',1,4,1666056192,1666058507);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(5,52,'test','','',1,11,1666061930,1666061930);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(6,52,'text','','',-1,10,1666061930,1666061930);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(7,52,'tuxt','','',1,9,1666061930,1666061930);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(8,52,'111','','',1,5,1666061947,1666061947);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(9,52,'222','','',1,6,1666061947,1666061947);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(10,52,'333','','["44444"]',1,7,1666061947,1666064132);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(11,52,'444','','',1,8,1666061947,1666061947);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(12,53,'ceshi','','',1,9999,1666079625,1666079625);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(13,54,'c','','',1,9999,1666080117,1666080117);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(14,46,'流动资产','','',1,1,1666232744,1666233502);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(15,46,'测试数据','','',1,2,1666232749,1666232749);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(16,55,'ceshi','','["测试科目","科技1"]',1,2,1666237608,1666237684);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(17,55,'123','','',-1,1,1666237611,1666237611);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(18,55,'321321','','["123"]',1,3,1666237614,1666237694);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(19,58,'货币资金','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(20,58,'存货','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(21,58,'流动资产','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(22,58,'交易性金融资产','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(23,58,'应收票据','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(24,58,'应收账款','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(25,58,'预付款项','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(26,58,'应收股利','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(27,58,'应收利息','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(28,58,'其他应收款','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(29,58,'其中：消耗性生物资产','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(30,58,'待摊费用','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(31,58,'一年内到期的非流动资产','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(32,58,'其他流动资产','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(33,58,'流动资产合计','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(34,58,'非流动资产','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(35,58,'可供出售金融资产','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(36,58,'持有至到期投资','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(37,58,'投资性房地产','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(38,58,'长期股权投资','','',1,9999,1666321721,1666321721);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(39,58,'长期应收款','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(40,58,'固定资产原价','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(41,58,'在建工程','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(42,58,'工程物资','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(43,58,'固定资产清理','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(44,58,'生产性生物资产','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(45,58,'油气资产','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(46,58,'无形资产','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(47,58,'开发支出','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(48,58,'商誉','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(49,58,'长期待摊费用','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(50,58,'递延所得税资产','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(51,58,'其他非流动资产','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(52,58,'非流动资产合计','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(53,58,'资产总计','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(54,58,'流动负债','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(55,58,'短期借款','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(56,58,'交易性金融负债','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(57,58,'应付票据','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(58,58,'应付账款','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(59,58,'预收款项','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(60,58,'应付职工薪酬','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(61,58,'应交税费','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(62,58,'应付利息','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(63,58,'应付股利','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(64,58,'其他应付款','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(65,58,'预提费用','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(66,58,'预计负债','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(67,58,'一年内到期的非流动负债','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(68,58,'其他流动负债','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(69,58,'流动负债合计','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(70,58,'非流动负债','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(71,58,'长期借款','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(72,58,'应付债券','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(73,58,'长期应付款','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(74,58,'专项应付款','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(75,58,'递延所得税负债','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(76,58,'其他非流动负债','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(77,58,'非流动负债合计','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(78,58,'负债合计','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(79,58,'所有者权益','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(80,58,'实收资本（股本）','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(81,58,'资本公积','','',1,9999,1666322410,1666322410);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(82,58,'盈余公积','','',1,9999,1666322633,1666322633);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(83,58,'未分配利润','','',1,9999,1666322633,1666322633);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(84,58,'减：库存股','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(85,58,'所有者权益合计','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(86,58,'负债和所有者权益合计','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(87,58,'结算备付金','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(88,58,'拆出资金','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(89,58,'以公允价值计量且其变动计入当期损益的金融资产','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(90,58,'衍生金融资产','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(91,58,'买入返售金融资产','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(92,58,'合同资产','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(93,58,'持有待售资产','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(94,58,'发放贷款和垫款','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(95,58,'债权投资','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(96,58,'其他债权投资','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(97,58,'其他权益工具投资','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(98,58,'其他非流动金融资产','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(99,58,'使用权资产','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(100,58,'拆入资金','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(101,58,'以公允价值计量且其变动计入当期损益的金融负债','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(102,58,'衍生金融负债','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(103,58,'合同负债','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(104,58,'卖出回购金融资产款','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(105,58,'吸收存款及同业存放','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(106,58,'其他综合收益','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(107,58,'专项储备','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(108,58,'一般风险准备','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(109,58,'归属于母公司股东权益合计','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(110,58,'少数股东权益','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(111,58,'其中：优先股','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(112,58,'向中央银行借款','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(113,58,'持有待售负债','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(114,58,'递延收益','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(115,58,'其他权益工具','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(116,58,'其中：原材料','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(117,58,'库存商品（产成品)','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(118,58,'减：累计折旧','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(119,58,'减：固定资产减值准备','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(120,58,'其中：土地使用权','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(121,58,'其中：特准储备物资','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(122,58,'应付权证','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(123,58,'其中：特种储备基金','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(124,58,'未确认的投资损失（以\-\号填列）','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(125,58,'外币报表折算差额','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(126,58,'固定资产净值','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(127,58,'固定资产净额','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(128,58,'永续债','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(129,58,'应收款项融资','','',1,9999,1666322634,1666322634);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(130,59,'营业收入','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(131,59,'减：营业成本','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(132,59,'税金及附加','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(133,59,'销售费用','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(134,59,'管理费用','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(135,59,'财务费用','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(136,59,'加：公允价值变动收益','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(137,59,'其中：对联营企业和合营企业的投资收益','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(138,59,'其他收益','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(139,59,'营业利润（亏损以“-”号填列）','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(140,59,'加：营业外收入','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(141,59,'其中：非流动资产处置利得','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(142,59,'减：营业外支出','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(143,59,'其中：非流动资产处置损失','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(144,59,'利润总额（亏损总额以“-”号填列）','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(145,59,'减：所得税费用','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(146,59,'净利润（净亏损以“-”号填列）','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(147,59,'持续经营净利润（净亏损以“-”号填列）','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(148,59,'终止经营净利润（净亏损以“-”号填列）','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(149,59,'其他综合收益的税后净额','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(150,59,'不能重分类进损益的其他综合收益','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(151,59,'重新计量设定受益计划净负债或净资产的变动','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(152,59,'权益法下在被投资单位不能重分类进损益的其他综合收益中享有的份额','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(153,59,'将重分类进损益的其他综合收益','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(154,59,'权益法下在被投资单位以后将重分类进损益的其他综合收益中享有的份额','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(155,59,'可供出售金融资产公允价值变动损益','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(156,59,'持有至到期投资重分类为可供出售金融资产损益','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(157,59,'现金流经套期损益的有效部分','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(158,59,'外币财务报表折算差额','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(159,59,'其他','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(160,59,'综合收益总额','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(161,59,'每股收益','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(162,59,'基本每股收益','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(163,59,'稀释每股收益','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(164,59,'加：年初未分配利润','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(165,59,'其他转入','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(166,59,'减：提取法宝盈余公积','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(167,59,'提取企业储备基金','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(168,59,'提取企业发展基金','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(169,59,'提取职工奖励及福利基金','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(170,59,'利润归还投资','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(171,59,'应付优先股股利','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(172,59,'提取任意盈余公积','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(173,59,'应付普通股股利','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(174,59,'转作资本（或股本）的普通股股利','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(175,59,'转总部利润','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(176,59,'未分配利润','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(177,59,'其中：主营业务收入','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(178,59,'其他业务收入','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(179,59,'营业总收入','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(180,59,'研发费用','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(181,59,'其中：利息费用','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(182,59,'利息收入','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(183,59,'加：其他收益','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(184,59,'投资收益（损失以“-”号填列）','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(185,59,'以摊余成本计量的金融资产终止确认收益','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(186,59,'净敝口套期收益（损失以“-”号填列）','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(187,59,'公允价值变动收益（损失以“-”号填列）','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(188,59,'信用减值损失（损失以“-”号填列）','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(189,59,'资产减值损失（损失以“-”号填列）','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(190,59,'资产处置收益（损失以“-”号填列）','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(191,59,'归属于母公司所有者的净利润（净亏损以“-”号填列））','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(192,59,'少数股东损益（净亏损以“-”号填列）','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(193,59,'其他业务成本','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(194,59,'其中：主营业务成本','','',1,9999,1666322924,1666322924);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(195,60,'经营活动产生的现金流量','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(196,60,'销售商品、提供劳务收到的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(197,60,'收到的税费返还','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(198,60,'收到其他与经营活动有关的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(199,60,'经营活动现金流入小计','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(200,60,'购买商品、接受劳务支付的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(201,60,'支付给职工以及为职工支付的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(202,60,'支付的各项税费','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(203,60,'支付其他与经营活动有关的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(204,60,'经营活动现金流出小计','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(205,60,'经营活动产生的现金流量净额','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(206,60,'投资活动产生的现金流量','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(207,60,'收回投资收到的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(208,60,'取得投资收益收到的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(209,60,'处置固定资产、无形资产和其他长期资产收回的现金净额','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(210,60,'处置子公司及其他营业单位收到的现金净额','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(211,60,'收到其他与投资活动有关的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(212,60,'投资活动现金流入小计','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(213,60,'购建固定资产、无形资产和其他长期资产支付的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(214,60,'投资支付的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(215,60,'取得子公司及其他营业单位支付的现金净额','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(216,60,'支付其他与投资活动有关的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(217,60,'投资活动现金流出小计','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(218,60,'投资活动产生的现金流量净额','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(219,60,'筹资活动产生的现金流量','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(220,60,'吸收投资收到的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(221,60,'取得借款收到的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(222,60,'收到其他与筹资活动有关的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(223,60,'筹资活动现金流入小计','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(224,60,'偿还债务支付的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(225,60,'分配股利、利润或偿付利息支付的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(226,60,'支付其他与筹资活动有关的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(227,60,'筹资活动现金流出小计','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(228,60,'筹资活动产生的现金流量净额','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(229,60,'汇率变动对现金及现金等价物的影响','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(230,60,'现金及现金等价物净增加额','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(231,60,'期初现金及现金等价物余额','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(232,60,'期末现金及现金等价物余额','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(233,60,'其中：子公司吸收少数股东投资收到的现金','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(234,60,'其中：子公司支付给少数股东的股利、利润','','',1,9999,1666323045,1666323045);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(235,59,'货币资金','','',1,9999,1666341420,1666341420);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(236,60,'货币资金','','',1,9999,1666341457,1666341457);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(237,59,'营业收入','','',1,9999,1666341849,1666341849);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(238,58,'营业收入','','',1,9999,1666341859,1666341859);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(239,60,'营业收入','','',1,9999,1666341868,1666341868);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(240,59,'其中：营业收入','','',1,9999,1666679103,1666679103);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(241,59,'利息收入','','',1,9999,1666679133,1666679133);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(242,59,'已赚保费','','',1,9999,1666679133,1666679133);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(243,59,'手续费及佣金收入','','',1,9999,1666679133,1666679133);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(244,59,'税金及附加','','',1,9999,1666700042,1666700042);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(245,59,'销售费用','','',1,9999,1666700042,1666700042);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(246,62,'销售费用','','',1,9999,1666700520,1666700520);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(247,62,'管理费用','','',1,9999,1666700520,1666700520);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(248,62,'财务费用','','',1,9999,1666700520,1666700520);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(249,61,'返','','',1,9999,1666702069,1666702069);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(250,63,'销售费用','','',1,9999,1666702328,1666702328);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(251,63,'分保费用','','',1,9999,1666703174,1666703174);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(252,63,'管理费用','','',1,9999,1666703174,1666703174);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(253,63,'财务费用','','',1,9999,1666703174,1666703174);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(254,60,'利息收入','','',1,9999,1667877748,1667877748);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(255,60,'已赚保费','','',1,9999,1667877748,1667877748);
INSERT INTO "SYSDBA"."standard_statement_titles"("id","statement_id","name","external_id","aliases","status","order_by_id","create_at","update_at") VALUES(256,60,'税金及附加','','',1,9999,1667877852,1667877852);

SET IDENTITY_INSERT "SYSDBA"."standard_statement_titles" OFF;
SET IDENTITY_INSERT "SYSDBA"."standard_statements" ON;
INSERT INTO "SYSDBA"."standard_statements"("id","standard_id","type","status","title_status","formula_status","create_at","update_at") VALUES(58,20,1,1,-1,-1,1666321480,1666321480);
INSERT INTO "SYSDBA"."standard_statements"("id","standard_id","type","status","title_status","formula_status","create_at","update_at") VALUES(59,20,2,1,-1,-1,1666321480,1666321480);
INSERT INTO "SYSDBA"."standard_statements"("id","standard_id","type","status","title_status","formula_status","create_at","update_at") VALUES(60,20,3,1,-1,-1,1666321480,1666321480);

SET IDENTITY_INSERT "SYSDBA"."standard_statements" OFF;
SET IDENTITY_INSERT "SYSDBA"."standards" ON;
INSERT INTO "SYSDBA"."standards"("id","name","external_id","is_default","status") VALUES(20,'新会计准则财报','',1,1);

SET IDENTITY_INSERT "SYSDBA"."standards" OFF;
SET IDENTITY_INSERT "SYSDBA"."tasks" ON;
SET IDENTITY_INSERT "SYSDBA"."tasks" OFF;
SET IDENTITY_INSERT "SYSDBA"."users" ON;
INSERT INTO "SYSDBA"."users"("id","name","account","password","salt","email","mobile","auth_method","created_at","updated_at","deleted_at") VALUES(1,'admin','admin','7M9t72i6fX05YysJbEr/uyvb5lEBHXbPdVTjc4AM4m0=','V7feAtVx','','','',1663835299,1663835299,null);

SET IDENTITY_INSERT "SYSDBA"."users" OFF;


commit;

