postgres12:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mMbzKhVc2DTye79dNfMts -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

psql:
	docker exec -it postgres12 psql -U root -d simple_bank

initdb: postgres12 sleep-2 createdb migrateup

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:mMbzKhVc2DTye79dNfMts@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:mMbzKhVc2DTye79dNfMts@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:mMbzKhVc2DTye79dNfMts@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:mMbzKhVc2DTye79dNfMts@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

server:
	go run main.go

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/xmeizh/simplebank/db/postgresql Store

sleep-%:
	sleep $(@:sleep-%=%)

.PHONY: postgres12 createdb dropdb migrateup migratedown migratedown1 migrateup1 sqlc psql server mock 
