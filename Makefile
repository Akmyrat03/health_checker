

pg-migrations-up:
	migrate -path internal/infrastructure/database -database "postgres://postgres:postgres@localhost:5432/checker?sslmode=disable" up