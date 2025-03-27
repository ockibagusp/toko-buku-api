CREATE TABLE IF NOT EXISTS book (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    author_id MEDIUMINT UNSIGNED NOT NULL,
    `type_id` SMALLINT UNSIGNED NULL,
    title VARCHAR(50) NOT NULL,
    sku VARCHAR(30) NOT NULL,
    price DECIMAL(6,2) NOT NULL DEFAULT 0.0,
    stock MEDIUMINT UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    CONSTRAINT author_id
        FOREIGN KEY (author_id)
        REFERENCES author (id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
    CONSTRAINT `type_id`
        FOREIGN KEY (`type_id`)
        REFERENCES `type` (id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE)
ENGINE = InnoDB;

START TRANSACTION;

INSERT INTO book (id, created_at, updated_at, deleted_at, author_id, type_id, title, sku, price, stock) VALUES (1, '2024-12-02 19:01:00', NULL, NULL, 1, 1, 'Tenggelamnya Kapal van der Wijck', 'kapal-van-der_1', 6.45, 100) ON DUPLICATE KEY UPDATE id=id;
INSERT INTO book (id, created_at, updated_at, deleted_at, author_id, type_id, title, sku, price, stock) VALUES (2, '2024-12-02 19:02:00', NULL, NULL, 2, 1, 'Bumi Manusia', 'bumi-manusia_1', 8.20, 100) ON DUPLICATE KEY UPDATE id=id;

COMMIT;