# main
main := cmd/main.go

start/sql:
	brew services start mysql

stop/sql:
	brew services stop mysql

info/sql:
	brew services info mysql

run-sql:
	@echo "------"
	@@brew services start mysql
	@echo "------"
# 1. golang run next or, (.. || ...)
# 2. mysql stop
	@@go run $(main) || brew services stop mysql

run:
	@go run $(main)

build:
	@go build $(main)

all: run-sql