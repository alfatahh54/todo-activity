include /${PWD}/.env
start:
		@go run .
dev:
		@gin --appPort 3020 --port 3030  --immediate run .
build:
		@go build .
run: build
		./fulfillment
migrateup:
		@migrate -database "mysql://${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" -path ${PWD}/migrations -verbose up ${n}
migratedown:
		@migrate -database "mysql://${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" -path ${PWD}/migrations -verbose down ${n}
migratecreate:
		@migrate create -ext sql -dir ${PWD}/migrations ${name}