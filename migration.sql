-- MySQL dump 10.13  Distrib 5.5.62, for Win64 (AMD64)
--
-- Host: 192.168.13.71    Database: zoho-center
-- ------------------------------------------------------
-- Server version	8.0.21

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `items`
--

DROP TABLE IF EXISTS `items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `items` (
  `id` int NOT NULL AUTO_INCREMENT,
  `zoho_id` varchar(32) NOT NULL DEFAULT '' COMMENT '商品对应zoho的item_id',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '商品名称',
  `sku` varchar(32) NOT NULL DEFAULT '' COMMENT '商品在zoho的sku',
  `status` varchar(32) NOT NULL DEFAULT '' COMMENT '商品状态',
  `um` varchar(32) NOT NULL DEFAULT '' COMMENT '商品单位',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '商品描述',
  `rate` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '价格',
  `initial_stock` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '商品初始库存',
  `initial_rate` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '商品初始库存单价',
  `purchase_rate` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '商品采购价',
  `sales_rate` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '商品销售价',
  `stock_on_hand` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '当前库存',
  `available_stock` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '可用库存',
  `actual_available_stock` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '实际可用库存',
  `vendor_id` varchar(32) NOT NULL DEFAULT '' COMMENT '供应商zohoid',
  `source` varchar(32) NOT NULL DEFAULT '' COMMENT '商品来源',
  `zoho_updated` timestamp NOT NULL COMMENT 'zoho更新时间',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `updated_by` varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `items`
--

LOCK TABLES `items` WRITE;
/*!40000 ALTER TABLE `items` DISABLE KEYS */;
/*!40000 ALTER TABLE `items` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `salesorder_items`
--

DROP TABLE IF EXISTS `salesorder_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `salesorder_items` (
  `id` int NOT NULL AUTO_INCREMENT,
  `zoho_id` varchar(32) NOT NULL DEFAULT '' COMMENT '销售订单商品对应zoho的line_item_id',
  `salesorder_id` varchar(32) NOT NULL DEFAULT '' COMMENT '销售订单的zoho_id',
  `item_id` varchar(32) NOT NULL DEFAULT '' COMMENT '商品的zoho_id',
  `sku` varchar(255) NOT NULL DEFAULT '' COMMENT '商品SKU',
  `item_name` varchar(255) NOT NULL DEFAULT '' COMMENT '商品名称',
  `rate` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '单价',
  `quantity` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '数量',
  `discount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '折扣',
  `item_total` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '商品总价',
  `tax_total` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '税费',
  `quantity_packed` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '已打包数量',
  `quantity_shipped` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '已配送数量',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `updated_by` varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `salesorder_items`
--

LOCK TABLES `salesorder_items` WRITE;
/*!40000 ALTER TABLE `salesorder_items` DISABLE KEYS */;
/*!40000 ALTER TABLE `salesorder_items` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `salesorders`
--

DROP TABLE IF EXISTS `salesorders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `salesorders` (
  `id` int NOT NULL AUTO_INCREMENT,
  `zoho_id` varchar(32) NOT NULL DEFAULT '' COMMENT '销售订单对应zoho的item_id',
  `salesorder_number` varchar(255) NOT NULL DEFAULT '' COMMENT '销售订单编码',
  `date` date DEFAULT NULL COMMENT '销售订单日期',
  `expected_shipment_date` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '发货日期',
  `customer_id` varchar(32) NOT NULL DEFAULT '' COMMENT '顾客ID',
  `customer_name` varchar(255) NOT NULL DEFAULT '' COMMENT '顾客名称',
  `order_status` varchar(32) NOT NULL DEFAULT '' COMMENT '销售订单状态',
  `invoiced_status` varchar(32) NOT NULL DEFAULT '' COMMENT 'invoice状态',
  `paid_status` varchar(32) NOT NULL DEFAULT '' COMMENT '付款状态',
  `shipped_status` varchar(32) NOT NULL DEFAULT '' COMMENT '发货状态',
  `source` varchar(32) NOT NULL DEFAULT '' COMMENT '商品来源',
  `salesperson_id` varchar(32) NOT NULL DEFAULT '' COMMENT '销售ID',
  `salesperson_name` varchar(255) NOT NULL DEFAULT '' COMMENT '销售名称',
  `shipping_charge` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '运费',
  `sub_total` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT 'subtotal',
  `discount_total` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '折扣',
  `tax_total` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '税费',
  `total` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '总价',
  `status` varchar(32) NOT NULL DEFAULT '' COMMENT '销售订单状态',
  `zoho_updated` timestamp NOT NULL COMMENT 'zoho更新时间',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `updated_by` varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `salesorders`
--

LOCK TABLES `salesorders` WRITE;
/*!40000 ALTER TABLE `salesorders` DISABLE KEYS */;
/*!40000 ALTER TABLE `salesorders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tokens`
--

DROP TABLE IF EXISTS `tokens`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tokens` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'id',
  `code` varchar(32) NOT NULL DEFAULT '' COMMENT '用户代码',
  `access_token` varchar(128) NOT NULL DEFAULT '' COMMENT 'api token',
  `api_domain` varchar(128) NOT NULL DEFAULT '' COMMENT 'api domain',
  `token_type` varchar(32) NOT NULL DEFAULT '' COMMENT 'token类型',
  `expires_time` datetime DEFAULT NULL COMMENT '过期时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tokens`
--

LOCK TABLES `tokens` WRITE;
/*!40000 ALTER TABLE `tokens` DISABLE KEYS */;
INSERT INTO `tokens` VALUES (1,'ozmas','1000.0dae986e8cc20075311f1bb083516e2b.69d9f3f9bd049b7dd27ed6e8d4296508','https://www.zohoapis.com.au','Bearer','2022-06-16 20:35:14');
/*!40000 ALTER TABLE `tokens` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'zoho-center'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-06-17 16:50:11
