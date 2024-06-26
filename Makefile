DB_SOURCE=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

network: 
	docker network create bank-network

postgres:
	docker run --name postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

psql:
	docker exec -it postgres psql -U root -d simple_bank

initdb: postgres sleep-2 createdb migrateup

dropdb:
	docker exec -it postgres dropdb simple_bank

create_db_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down 1

sqlc:
	sqlc generate

server:
	go run main.go

test:
	go test -v -cover -short ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/xmeizh/simplebank/db/postgresql Store
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/xmeizh/simplebank/worker TaskDistributor

sleep-%:
	sleep $(@:sleep-%=%)

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml  

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --port 9090 -r repl --package pb --service SimpleBank

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: postgres createdb dropdb migrateup migratedown migratedown1 migrateup1 sqlc psql server network mock db_docs db_schema proto evans redis create_db_migration
