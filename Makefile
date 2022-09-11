DB_NAME ?= go-buckpal-db
DB_PORT ?= 3308
DB_DSN ?= root:1@tcp(localhost:${DB_PORT})/${DB_NAME}

start:
	DB_DSN="${DB_DSN}" go run ./cmd

db_run:
	docker run -d --name ${DB_NAME} -e MYSQL_DATABASE=${DB_NAME} -e MYSQL_ROOT_PASSWORD=1 -e MYSQL_ROOT_HOST=172.17.0.1 -p ${DB_PORT}:3306 mysql:8.0.30
db_start:
	docker start ${DB_NAME}
db_stop:
	docker stop ${DB_NAME}
db_remove:
	docker rm --force ${DB_NAME}
db_migrate_create:
	docker run -v $(shell pwd)/migrations:/migrations --network host migrate/migrate:v4.15.2 create -dir /migrations -ext sql -seq ${name}
db_migrate_up:
	docker run -v $(shell pwd)/migrations:/migrations --network host migrate/migrate:v4.15.2 -path=/migrations/ -database "mysql://${DB_DSN}" up

gomock_install:
	go install github.com/golang/mock/mockgen@v1.6.0
# example: make gomock_generate path=pkg/account/application/port/out file=update_account_state_port.go
gomock_generate:
	mockgen -source=${path}/${file} -destination=${path}/mock/${file} 