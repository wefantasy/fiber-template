CREATE TABLE `user` (
  `id` int PRIMARY KEY AUTO_INCREMENT COMMENT '编号',
  `username` varchar(256) NOT NULL UNIQUE COMMENT '用户账户',
  `password` varchar(256) NOT NULL COMMENT '用户密码',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期 默认为当前时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新日期',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除日期'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';