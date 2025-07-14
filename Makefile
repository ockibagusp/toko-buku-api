# main
main := cmd/main.go

SLEEP_COUNT = 1 # <- 3, e.g.

start/sql:
	brew services start mysql

stop/sql:
	brew services stop mysql

info/sql:
	brew services info mysql

run:
	@go run $(main)

all: start/sql
	@@sleep $(SLEEP_COUNT)
	@echo "------"
# 1. golang run next (or), (.. || ...)
# 2. print "------"
# 2. stop mysql
	@@go run $(main) \
		|| echo "------" \
		&& brew services stop mysql

air:
	@make start/sql 
	@@sleep $(SLEEP_COUNT)
	@air \
		|| @echo "------" \
		&& brew services stop mysql	

build:
	@go build $(main)
