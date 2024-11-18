# main
main := main.go

sql-start:
	brew services start mysql

sql-stop:
	brew services stop mysql

sql-info:
	brew services info mysql

run:
	go run $(main)

build:
	go build $(main)