postgres:
	docker run --name postgres-apecalendar -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it postgres-apecalendar createdb --username=root --owner=root apecalendar
dropdb:
	docker exec -it postgres-apecalendar dropdb --username=root --owner=root apecalendar
test:
	go test -v -cover ./...
migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/apecalendar?sslmode=disable" --verbose up
migrateup1:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/apecalendar?sslmode=disable" --verbose up 1
migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/apecalendar?sslmode=disable" --verbose down
migratedown1:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/apecalendar?sslmode=disable" --verbose down 1
sqlc:
	sqlc generate
server:
	go run main.go
air:
	air run main.go

# gomock
# mockgen --build_flags=--mod=mod -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store
migrate-create:
	migrate create -ext sql -dir db/migrations -seq
