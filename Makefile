pg-migrations-up:
	migrate -path internal/infrastructure/pg_migrations -database "postgres://postgres:postgres@localhost:5432/checker?sslmode=disable" up