include .env
export

run:
	go run cmd/api/main.go
doc:
	go tool swag init -g cmd/api/main.go
migrate:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate -source file://database/migrations -database "$(DATABASE_URL)" up
seed:
	go run cmd/seeder/main.go