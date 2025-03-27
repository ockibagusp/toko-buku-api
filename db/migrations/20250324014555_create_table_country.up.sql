CREATE TABLE IF NOT EXISTS country (
    id TINYINT UNSIGNED NOT NULL AUTO_INCREMENT,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    iso3 VARCHAR(3) NOT NULL,
    country VARCHAR(50) NOT NULL,
    nice_country VARCHAR(50) NOT NULL,
    currency VARCHAR(8) NOT NULL,
    PRIMARY KEY (id)
) ENGINE = InnoDB;

START TRANSACTION;

INSERT INTO country (id, updated_at, iso3, country, nice_country, currency) VALUES (100, NULL, 'IDN', 'INDONESIA', 'Indonesia', 'Rp') ON DUPLICATE KEY UPDATE id=id;
INSERT INTO country (id, updated_at, iso3, country, nice_country, currency) VALUES (225, NULL, 'GBR', 'UNITED KINGDOM', 'United Kingdom', '£') ON DUPLICATE KEY UPDATE id=id;
INSERT INTO country (id, updated_at, iso3, country, nice_country, currency) VALUES (226, NULL, 'USA', 'UNITED STATES', 'United States', '$') ON DUPLICATE KEY UPDATE id=id;

COMMIT;