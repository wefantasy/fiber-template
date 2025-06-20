CREATE TABLE IF NOT EXISTS `user`
(
    `id`         INT PRIMARY KEY AUTO_INCREMENT COMMENT '编号',
    `username`   VARCHAR(255) NOT NULL UNIQUE COMMENT '用户账户',
    `password`   VARCHAR(255) NOT NULL COMMENT '用户密码',
    `created_at` TIMESTAMP    NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期 默认为当前时间',
    `updated_at` TIMESTAMP         DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
    `deleted_at` TIMESTAMP         DEFAULT NULL COMMENT '删除日期'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='用户表';