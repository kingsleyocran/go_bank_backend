postgresinit:
	docker run --name postgres142 -p 5432:5432 -e POSTGRES_USER=root  -e POSTGRES_PASSWORD=secret -d postgres:14.2-alpine

postgresinitdb:
	docker run --name postgres142 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_DB=simple_bank  -e POSTGRES_PASSWORD=secret -d postgres:14.2-alpine

dockerstart: 
	docker start postgres142

dockerstop:
	docker stop postgres142

psql: 
	docker exec -it postgres142 psql -d simple_bank -U root

createdb: 
	docker exec -it postgres142 createdb -U root -O root simple_bank

dropdb: 
	docker exec -it postgres142 dropdb simple_bank

migratecreate:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlcinit:
	sqlc init

sqlcgen:
	sqlc generate


.PHONY: createdb dropdb postgresinit postgresinitdb dockerstart dockerstop psql migrateup migratedown migratecreate sqlcgen sqlcinit