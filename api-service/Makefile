DB_URL=pgx://root:secret@localhost:5432/galore_db?sslmode=disable

network:
	docker network create galore-network
postgres:
	docker run --name galore-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine
createdb:
	docker exec -it galore-postgres createdb --username=root --owner=root galore_db
dropdb:
	docker exec -it galore-postgres dropdb galore_db
migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up
migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down
#evans:
#	evans --host localhost --port 9090 -r repl
#redis:
#	docker run --name redis -p 6379:6379 -d redis:7-alpine
new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)
#db_docs:
#	dbdocs build doc/db.dbml
#db_schema:
#	dbml2sql --postgres -o doc/schema.sql doc/db.dbml
sqlc:
	sqlc generate
test:
	go test -v -cover -short ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/m1thrandir225/galore-services/db/sqlc Store
	mockgen -package mockstorage -destination storage/mock/storage.go github.com/m1thrandir225/galore-services/storage FileService
	mockgen -package mockcache -destination cache/mock/cache.go github.com/m1thrandir225/galore-services/cache KeyValueStore
	mockgen -package mockembedding -destination embedding_service/mock/embedding.go github.com/m1thrandir225/galore-services/embedding_service EmbeddingService
	mockgen -package mockcocktailgen -destination cocktail_gen/mock/cocktail_gen.go github.com/m1thrandir225/galore-services/cocktail_gen CocktailGenerator
	mockgen -package mockimagegen -destination image_gen/mock/image_gen.go github.com/m1thrandir225/galore-services/image_gen ImageGenerator
	mockgen -package mockcategorize -destination categorizer_service/mock/categorizer.go github.com/m1thrandir225/galore-services/categorizer_service CategorizerService
#proto:
#	rm -f pb/*.go
#	rm -f doc/swagger/*.swagger.json
#	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
#	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
#	--grpc-gateway_out=pb --gr

.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration db_docs db_schema sqlc test server mock proto evans redisURL)" -verbose down 1
