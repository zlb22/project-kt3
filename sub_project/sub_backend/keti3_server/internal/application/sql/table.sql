-- 用户操作日志
drop table if exists keti3_op_log;
CREATE TABLE `keti3_op_log` (
  `id` int(10) unsigned NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'ID',
  `uid` int NOT NULL DEFAULT 0 COMMENT '用户ID',
  `op_time` int NOT NULL DEFAULT 0 COMMENT '操作时间',
  `op_type` varchar(32) NOT NULL DEFAULT '' COMMENT '操作类型',
  `op_object` varchar(20) NOT NULL DEFAULT '' COMMENT '操作对象',
  `object_name` varchar(32) NOT NULL DEFAULT '' COMMENT '操作对象名称',
  `object_no` int NOT NULL DEFAULT '0' COMMENT '操作对象编号',
  `data_before` text NOT NULL  COMMENT '更改之前的数据',
  `data_after` text NOT NULL  COMMENT '更改之后的数据',
  `voice_url` varchar(255) NOT NULL DEFAULT '' COMMENT '音频地址',
  `screenshot_url` varchar(255) NOT NULL DEFAULT '' COMMENT '截图地址',
  `deleted_at` int(10) NOT NULL DEFAULT '0' COMMENT '删除时间',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  KEY `idx_uid` (`uid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户操作日志表';





-- op_type操作类型
-- 新增 add
-- 移动 move
-- 删除 delete
-- 缩放 scale
-- 宽高 aspect
-- 旋转 rotate
-- 翻转 flip
-- 扭曲 distortion
-- 撤销 undo
-- 重做 redo
-- 下一步 save 
