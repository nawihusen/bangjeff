CREATE TABLE `user` (
    `id`                    INT                     NOT NULL AUTO_INCREMENT,
    `username`              VARCHAR(100)            NOT NULL,
    `password`              VARCHAR(100)            NOT NULL,
    `name`                  VARCHAR(100)            NOT NULL,
    `phone`                 VARCHAR(100)            NOT NULL,
    `email`                 VARCHAR(100)            NOT NULL,
    `address`               VARCHAR(100)            NOT NULL,
    `dtm_crt`               TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `token`                 VARCHAR(100),
    `active`                INT,
    PRIMARY KEY (`id`)
);
