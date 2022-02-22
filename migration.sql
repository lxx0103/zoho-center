/*token表*/
CREATE TABLE `tokens` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'id',
  `code` varchar(32) not null default '' comment '用户代码',
  `access_token` varchar(128) not null default '' comment 'api token',
  `api_domain` varchar(128) not null default '' comment 'api domain',
  `token_type` varchar(32) not null default '' comment 'token类型',
  `expires_time` datetime DEFAULT NULL comment '过期时间',
  PRIMARY KEY (`id`),
  Unique Key (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

/*token 初始化*/
INSERT INTO tokens (`code`, `access_token`, `api_domain`, `token_type`, `expires_time`) VALUES ('ozmas', '', '', '', '1970-01-01');

/*商品表*/
CREATE TABLE items (
	`id` int not null auto_increment,
	`zoho_id` varchar(32) not null default '' comment '商品对应zoho的item_id',
	`name` varchar(255) not null default '' comment '商品名称',
	`sku` varchar(32) not null default '' comment '商品在zoho的sku',
	`status` varchar (32) not null default '' comment  '商品状态',
	`um` varchar(32) not null default '' comment '商品单位',
	`desc` varchar(255) not null default '' comment '商品描述',
	`initial_stock` decimal(10,2) not null default 0.00 comment '商品初始库存',
	`initial_rate` decimal(10,2) not null default 0.00 comment '商品初始库存单价',
	`purchae_rate` decimal(10,2) not null default 0.00 comment '商品采购价',
	`sales_rate` decimal(10,2) not null default 0.00 comment '商品销售价',
	`vendor_id` varchar(32) not null default '' comment '供应商zohoid',
	`source` varchar(32) not null default '' comment '商品来源',
	`zoho_updated` timestamp not null comment 'zoho更新时间',	
	`created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  	`created_by` varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
  	`updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  	`updated_by` varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci