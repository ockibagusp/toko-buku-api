CREATE TABLE IF NOT EXISTS order_has_book (
    order_id BIGINT UNSIGNED NOT NULL,
    book_id INT UNSIGNED NOT NULL,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (order_id, book_id),
    CONSTRAINT fk_order_has_book_order1
        FOREIGN KEY (order_id)
        REFERENCES `order` (id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,
    CONSTRAINT fk_order_has_book_book1
        FOREIGN KEY (book_id)
        REFERENCES book (id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION)
ENGINE = InnoDB;