DB_URL=postgresql://root:secret@localhost:5432/bank_server?sslmode=disable

network:
	docker network create bank-network
postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root bank_server
dropdb:
	docker exec -it postgres12 dropdb --username=root bank_server
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1
new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)
db_docs:
	dbdocs build doc/db.dbml
db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml 
sqlc:
	sqlc generate
test:
	go test -v -cover -coverprofile=coverage.out -short ./...
test_json:
	go test -v -json > test-report.json ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/umarhadi/bank-server/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/umarhadi/bank-server/worker TaskDistributor
	mockgen -package mockemail -destination mail/mock/sender.go github.com/umarhadi/bank-server/mail EmailSender
proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=bank-server proto/*.proto
	statik -src=./doc/swagger -dest=./doc
evans:
	evans --host localhost --port 9090 -r repl
redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine
.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration sqlc test server mock db_docs db_schema proto evans redis