postgres:
	docker run --name receiptdb -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:17-alpine

console:
	docker exec -it receiptdb psql -U root

createdb:
	docker exec -it receiptdb createdb --username=root --owner=root receipt_processor

dropdb:
	docker exec -it receiptdb dropdb receipt_processor

migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/receipt_processor?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/receipt_processor?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres console createdb dropdb migrateup migratedown sqlc test server