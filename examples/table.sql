CREATE TABLE `acesso` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `advertising_default_id` int(11) DEFAULT NULL,
  `advertising_by_car_id` int(11) DEFAULT NULL,
  `marca` varchar(100) COLLATE latin1_general_ci NOT NULL,
  `modelo` varchar(100) COLLATE latin1_general_ci NOT NULL,
  `ano` int(4) NOT NULL,
  `combustivel` varchar(50) COLLATE latin1_general_ci NOT NULL,
  `propriedade` varchar(50) COLLATE latin1_general_ci NOT NULL,
  `latitude` decimal(18,12) NOT NULL,
  `longitude` decimal(18,12) NOT NULL,
  `date` date DEFAULT NULL,
  `time` time DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `advertising_acesso` (`advertising_default_id`),
  KEY `advertising_by_car` (`advertising_by_car_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4417 DEFAULT CHARSET=latin1 COLLATE=latin1_general_ci;
