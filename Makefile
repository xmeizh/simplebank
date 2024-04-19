postgres12:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mMbzKhVc2DTye79dNfMts -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:mMbzKhVc2DTye79dNfMts@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:mMbzKhVc2DTye79dNfMts@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose down

.PHONY: postgres12 createdb dropdb migrateup migratedown
