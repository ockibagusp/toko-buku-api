CREATE TABLE IF NOT EXISTS `order` (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    x_idempotency_key VARCHAR(38) NOT NULL COMMENT 'Idempotency  http POST.\n\nHTTP Header:\nx-idempotency-key=….',
    status ENUM("processed", "received") NOT NULL DEFAULT 'processed',
    count TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Min.: 1 pc.\nMax.: 15 pcs.',
    PRIMARY KEY (id)
) ENGINE = InnoDB;

CREATE UNIQUE INDEX `idempotency_key_UNIQUE` ON `toko-buku-api`.`order` (`x-idempotency-key` ASC) VISIBLE;
