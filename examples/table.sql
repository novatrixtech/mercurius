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