CREATE TABLE `acesso` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `advertising_default_id` int(11) DEFAULT NULL,
  `advertising_by_car_id` int(11) DEFAULT NULL,
  `marca` varchar(100) NOT NULL,
  `modelo` varchar(100) NOT NULL,
  `ano` int(4) NOT NULL,
  `combustivel` varchar(50) NOT NULL,
  `propriedade` varchar(50) NOT NULL,
  `latitude` decimal(18,12) NOT NULL,
  `longitude` decimal(18,12) NOT NULL,
  `date` date DEFAULT NULL,
  `time` time DEFAULT NULL,
  PRIMARY KEY (`id`)
);
INSERT INTO `acesso` VALUES (273,3,NULL,'teste','teste',2001,'flex','dono',-7.999898000000,10.984729472398,'2016-05-05','20:48:17'),(274,3,NULL,'teste','teste',2001,'flex','dono',-7.999898000000,10.984729472398,'2016-05-05','20:48:19'),(275,1,NULL,'Chevrolet','Spin',2016,'flex','dono',-23.531531531532,-46.654543655696,'2016-05-05','20:50:20'),(276,1,NULL,'Chevrolet','Spin',2016,'flex','dono',-23.531531531532,-46.654543655696,'2016-05-05','20:50:24');

CREATE TABLE `logac_accesstokenacessos` (
  `logac_id` INT NOT NULL,
  `logac_accesstoken` TEXT NULL,
  `logac_function` VARCHAR(45) NULL,
  `logac_when` DATETIME NULL,
  PRIMARY KEY (`logac_id`));

  CREATE TABLE `logacr_logaccesstokenrequest` (
  `logacr_id` INT NOT NULL,
  `logacr_accesstoken` TEXT NULL,
  `sys_id` INT NULL,
  `logacr_when` DATETIME NULL,
  PRIMARY KEY (`logacr_id`));

  CREATE TABLE `rol_roles` (
  `rol_id` INT NOT NULL,
  `rol_description` VARCHAR(45) NULL,
  PRIMARY KEY (`rol_id`));

  CREATE TABLE `rsys_rolesystems` (
  `rsys_id` INT NOT NULL,
  `rol_id` INT NULL,
  `sys_id` INT NULL,
  PRIMARY KEY (`rsys_id`));