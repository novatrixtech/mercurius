CREATE TABLE `logcli_loginclient` (
  `logcli_id` INT NOT NULL AUTO_INCREMENT,
  `logcli_clientname` VARCHAR(255) NULL,
  `logcli_clientlegacyid` VARCHAR(255) NULL,
  `logcli_clientid` VARCHAR(255) NULL,
  `logcli_secret` VARCHAR(255) NULL,
  `logcli_role` VARCHAR(255) NOT NULL,
  `logcli_lastupdate` DATETIME NOT NULL DEFAULT now(),
  PRIMARY KEY (`logcli_id`));

CREATE TABLE `logac_accesstokenacessos` (
  `logac_id` INT NOT NULL AUTO_INCREMENT,
  `logac_quando` DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
  `logac_accesstoken` VARCHAR(240) NULL,
  `logac_funcao` VARCHAR(240) NULL,
  PRIMARY KEY (`logac_id`));
   

--Adicionar a tabela para armazenar os access token gerados quando e por quem

CREATE TABLE `logacr_logaccesstokenrequest` (
  `logacr_id` INT NOT NULL AUTO_INCREMENT,
  `logacr_accesstoken` VARCHAR(200) NOT NULL,
  `contato_id` INT NOT NULL,
  `logacr_when` DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`logacr_id`));
