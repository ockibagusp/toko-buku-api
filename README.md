# toko-buku-api

## Air: cloud Live reload for Go apps

```sh
$ air
```

[Air](https://github.com/cosmtrek/air) is yet another live-reloading command line utility for Go applications in development.

## Golang Migrate

[Golang Migrate](https://github.com/golang-migrate/migrate/blob/master/README.md)

### Install

```sh
$ go install -tags "mysql" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Checking database

```sh
$ migrate -help
```

Table Country

```sh
$ migrate create -ext sql -dir db/migrations create_table_country
```

```sh
$ migrate -path db/migrations -database "mysql://root:rahasia\!@tcp(localhost:3306)/toko-buku-api?charset=utf8mb4&parseTime=True" up
```

=> password: `rahasia!` menjadi `rahasia\!`

Table Author

```sh
$ migrate create -ext sql -dir db/migrations create_table_author
```

### Force database

(Force database version)[https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md#forcing-your-database-version]

```sh
$ migrate -path db/migrations -database "mysql://root:rahasia\!@tcp(localhost:3306)/toko-buku-api?charset=utf8mb4&parseTime=True" force 1
```
