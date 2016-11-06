-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `mydb` DEFAULT CHARACTER SET utf8 ;
USE `mydb` ;

-- -----------------------------------------------------
-- Table `mydb`.`dimension_department`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`dimension_department` (
  `department` VARCHAR(128) NOT NULL,
  `opaque` VARCHAR(45) NULL,
  PRIMARY KEY (`department`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`dimension_group`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`dimension_group` (
  `egroup` VARCHAR(128) NOT NULL,
  `opaque` VARCHAR(45) NULL,
  PRIMARY KEY (`egroup`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`dimension_date`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`dimension_date` (
  `ts` INT NOT NULL,
  `day` INT NOT NULL,
  `month` INT NOT NULL,
  `year` INT NOT NULL,
  PRIMARY KEY (`ts`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`dimension_company`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`dimension_company` (
  `company` VARCHAR(128) NOT NULL,
  `opaque` VARCHAR(45) NULL,
  PRIMARY KEY (`company`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`fact_shares`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`fact_shares` (
  `id` INT NOT NULL,
  `owner_login` VARCHAR(45) NOT NULL,
  `owner_uid` INT NOT NULL,
  `owner_department` VARCHAR(128) NOT NULL,
  `owner_group` VARCHAR(128) NOT NULL,
  `owner_company` VARCHAR(128) NOT NULL,
  `sharee_login` VARCHAR(45) NOT NULL,
  `sharee_uid` INT NOT NULL,
  `sharee_department` VARCHAR(128) NOT NULL,
  `sharee_group` VARCHAR(128) NOT NULL,
  `sharee_company` VARCHAR(128) NOT NULL,
  `stime` INT NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fact_shares_stime_date_ts`
    FOREIGN KEY (`stime`)
    REFERENCES `mydb`.`dimension_date` (`ts`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fact_shares_owner_department_dimension_department_department`
    FOREIGN KEY (`owner_department`)
    REFERENCES `mydb`.`dimension_department` (`department`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fact_shares_owner_group_group_group`
    FOREIGN KEY (`owner_group`)
    REFERENCES `mydb`.`dimension_group` (`egroup`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fact_shares_owner_company_company_company`
    FOREIGN KEY (`owner_company`)
    REFERENCES `mydb`.`dimension_company` (`company`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fact_shares_sahree_department_department_department`
    FOREIGN KEY (`sharee_department`)
    REFERENCES `mydb`.`dimension_department` (`department`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fact_shares_sharee_group_group_group`
    FOREIGN KEY (`sharee_group`)
    REFERENCES `mydb`.`dimension_group` (`egroup`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fact_shares_sharee_company_company_company`
    FOREIGN KEY (`sharee_company`)
    REFERENCES `mydb`.`dimension_company` (`company`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;

CREATE INDEX `fact_shares_stime_ts_idx` ON `mydb`.`fact_shares` (`stime` ASC);

CREATE INDEX `fact_shares_owner_department_dimension_department_departmen_idx` ON `mydb`.`fact_shares` (`owner_department` ASC);

CREATE INDEX `fact_shares_owner_group_group_group_idx` ON `mydb`.`fact_shares` (`owner_group` ASC);

CREATE INDEX `fact_shares_owner_company_company_company_idx` ON `mydb`.`fact_shares` (`owner_company` ASC);

CREATE INDEX `fact_shares_sahree_department_department_department_idx` ON `mydb`.`fact_shares` (`sharee_department` ASC);

CREATE INDEX `fact_shares_sharee_group_group_group_idx` ON `mydb`.`fact_shares` (`sharee_group` ASC);

CREATE INDEX `fact_shares_sharee_company_company_company_idx` ON `mydb`.`fact_shares` (`sharee_company` ASC);


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
