-- MySQL dump 10.13  Distrib 8.0.27, for Linux (x86_64)
--
-- Host: localhost    Database: microdb
-- ------------------------------------------------------
-- Server version	8.0.27

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
--
-- Database microdb
--
DROP DATABASE IF EXISTS `microdb`; 
CREATE DATABASE `microdb` default character set=utf8;
USE `microdb`;
GRANT ALL PRIVILEGES ON microdb.* TO 'micro'@'%' ;
FLUSH PRIVILEGES;

--
-- Table structure for table `article`
--

DROP TABLE IF EXISTS `article`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `article` (
  `articleId` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '类型',
  `subject` varchar(500) NOT NULL DEFAULT '' COMMENT '标题',
  `addTime` datetime NOT NULL DEFAULT '2020-11-19 00:00:00' COMMENT '添加时间',
  `isHot` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否热榜,0:非热，1，热',
  `url` varchar(500) NOT NULL DEFAULT '' COMMENT '链接地址',
  `domain` varchar(100) NOT NULL DEFAULT '' COMMENT '域',
  `userId` bigint unsigned NOT NULL DEFAULT '0' COMMENT '推荐用户',
  `approvalStaffId` int unsigned NOT NULL DEFAULT '0' COMMENT '批准员工',
  `digSum` int unsigned NOT NULL DEFAULT '0' COMMENT '被顶的次数',
  `commentSum` int unsigned NOT NULL DEFAULT '0' COMMENT '被评论的次数',
  `isPublish` tinyint NOT NULL DEFAULT '0' COMMENT '0,未发布，1，已发布',
  PRIMARY KEY (`articleId`),
  KEY `isPublish` (`isPublish`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='文章表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `article`
--

LOCK TABLES `article` WRITE;
/*!40000 ALTER TABLE `article` DISABLE KEYS */;
/*!40000 ALTER TABLE `article` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `userId` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `role` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '角色: 0 普通用户，1 管理员',
  `type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '类型: 0 未付费 ，1 付费',
  `name` varchar(500) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(20) NOT NULL DEFAULT '' COMMENT '密码',
  `email` varchar(20) NOT NULL DEFAULT '' COMMENT '邮箱',
  `phone` varchar(20) NOT NULL DEFAULT '' COMMENT '电话号码',
  `person_id` varchar(18) NOT NULL DEFAULT '' COMMENT '生份证号',
  `emp_id` varchar(100) NOT NULL DEFAULT '' COMMENT '员工号',
  `gender` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '性别：0 女，1 男',
  `age` int unsigned NOT NULL DEFAULT '0' COMMENT '年龄',
  `introduce` varchar(100) NOT NULL DEFAULT '' COMMENT '自我介绍',
  `hobby` varchar(100) NOT NULL DEFAULT '' COMMENT '个人爱好',
  `buDomain` varchar(100) NOT NULL DEFAULT '' COMMENT '部门',
  `typeDevice` varchar(100) NOT NULL DEFAULT '' COMMENT '设备类型',
  `createTime` datetime NOT NULL DEFAULT '2020-11-19 00:00:00' COMMENT '创建时间',
  `isHot` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否热榜,0:非热，1，热',
  `approvalStaffId` int unsigned NOT NULL DEFAULT '0' COMMENT '审批人',
  `restSum` int unsigned NOT NULL DEFAULT '0' COMMENT '重置密码次数',
  `failSum` int unsigned NOT NULL DEFAULT '0' COMMENT '登陆失败次数',
  `isLogin` tinyint NOT NULL DEFAULT '0' COMMENT '状态:0,离线，1，在线',
  PRIMARY KEY (`userId`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-01-06  2:57:43
