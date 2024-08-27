DB_URL=pgx://root:secret@localhost:5432/galore_db?sslmode=disable

network:
	docker network create loits-network
postgres:
	docker run --name galore-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine
mysql:
	docker run --name mysql8 -p 3306:3306  -e MYSQL_ROOT_PASSWORD=secret -d mysql:8
createdb:
	docker exec -it galore-postgres createdb --username=root --owner=root loits_db
dropdb:
	docker exec -it galore-postgres dropdb loits_db
migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up
migrateup1:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up 1
migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down
migratedown1:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down 1
new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)
db_docs:
	dbdocs build doc/db.dbml
db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml
sqlc_init:
		docker run --rm -v "F:\Projects\simple_bank:/src" -w /src sqlc/sqlc init
sqlc:
	sqlc generate
test:
	go test -v -cover -short ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go m1thranidr/loits/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go m1thranidr/loits/worker TaskDistributor
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
	evans --host localhost --port 9090 -r repl
redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine
.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration db_docs db_schema sqlc test server mock proto evans redis
