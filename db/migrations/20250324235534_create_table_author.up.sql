CREATE TABLE IF NOT EXISTS author (
    id MEDIUMINT UNSIGNED NOT NULL AUTO_INCREMENT,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    country_id TINYINT UNSIGNED NOT NULL,
    author VARCHAR(50) NOT NULL,
    city VARCHAR(50) NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_author_country
        FOREIGN KEY (country_id)
        REFERENCES country (id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE)
ENGINE = InnoDB;

START TRANSACTION;

INSERT INTO author (id, updated_at, country_id, author, city) VALUES (1, '2024-12-02 17:22:47', 100, 'Buya Hamka', 'Sumatera Barat, Indonesia') ON DUPLICATE KEY UPDATE author=author;
INSERT INTO author (id, updated_at, country_id, author, city) VALUES (2, '2024-12-02 17:22:47', 100, 'Pramoedya Ananta Toer', 'Jawa Timur, Indonesia') ON DUPLICATE KEY UPDATE author=author;

COMMIT;