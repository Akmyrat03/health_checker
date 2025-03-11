root_dir:= $(shell pwd)
config_file := ${root_dir}/config.json
docker_compose_conf_file:= $(root_dir)/docker/docker-compose.yaml
dockerfile:= $(root_dir)/docker/Dockerfile
image_name:= health_checker

docker-build:
	docker build -t $(image_name) -f $(dockerfile) .

docker-compose-up:
	docker-compose -f $(docker_compose_conf_file) up -d

docker-compose-down:
	docker-compose -f $(docker_compose_conf_file) down

pg-migrations-up:
	migrate -path internal/infrastructure/pg_migrations -database "postgres://postgres:postgres@localhost:5432/checker?sslmode=disable" up

pg-migrations-down:
	migrate -path internal/infrastructure/pg_migrations -database "postgres://postgres:postgres@localhost:5432/checker?sslmode=disable" down
