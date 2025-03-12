Upgrading from MySQL <8.4 to MySQL >9.0 requires running MySQL 8.4 first:

- brew services stop mysql
- brew install mysql@8.4
- brew services start mysql@8.4
- brew services stop mysql@8.4
- brew services start mysql

We've installed your MySQL database without a root password. To secure it run:
mysql_secure_installation

MySQL is configured to only allow connections from localhost by default

To connect run:
mysql -u root

db: toko-buku-api
pass: rahasia!

POST /api/products

sku: novels-name-{no}

- set novels-name-{no} baru, sukses
- set novels-name-{no} ada, gagal

## id, name, sku

1, bumi manusia, novels-foo-1

DB

- author
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `city` varchar(255) DEFAULT NULL,

- book
  `id` int NOT NULL AUTO_INCREMENT,
  `author_id` int NOT NULL,
  `type_id` int DEFAULT NULL,
  `name` varchar(38) NOT NULL,
  `price` decimal(8,2) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `FK_authorBook` (`author_id`),
  CONSTRAINT `FK_authorBook` FOREIGN KEY (`author_id`) REFERENCES `author` (`id`)

- type
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,

https://github.com/swaggo/swag
$ go install github.com/swaggo/swag/cmd/swag@latest
$ swag init
$ go get -u github.com/swaggo/http-swagger

https://reqbin.com/
