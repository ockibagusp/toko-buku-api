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
