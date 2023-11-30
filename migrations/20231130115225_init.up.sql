CREATE TABLE IF NOT EXISTS activity (
  `id` INT NOT NULL AUTO_INCREMENT,
  `email` VARCHAR(45) NULL,
  `title` VARCHAR(45) NULL,
  `created_at` DATETIME NULL DEFAULT NOW(),
  `updated_at` DATETIME NULL DEFAULT NOW(),
  `deleted_at` DATETIME NULL,
  PRIMARY KEY (`id`) 
);

CREATE TABLE IF NOT EXISTS todo (
  `id` INT NOT NULL AUTO_INCREMENT,
  `activity_group_id` INT NULL,
  `title` VARCHAR(45) NULL,
  `created_at` DATETIME NULL DEFAULT NOW(),
  `updated_at` DATETIME NULL DEFAULT NOW(),
  `deleted_at` DATETIME NULL,
  `is_active` TINYINT(1) NULL,
  `priority` VARCHAR(45) NULL,
  PRIMARY KEY (`id`),
  INDEX `todo_idx1` (`activity_group_id` ASC) INVISIBLE,
  CONSTRAINT `todo_fk1`
    FOREIGN KEY (`activity_group_id`)
    REFERENCES `activity` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);
