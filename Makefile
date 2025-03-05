# Define root directory
root_dir := $(shell pwd)
docker_compose_conf_file := ${root_dir}/docker/docker-compose.yaml
docker_compose_local_conf_file := ${root_dir}/docker/docker-compose.local.yaml

# Include .env if exists
ifneq ("$(wildcard .env)","")
	include .env
endif

# PostgreSQL migration targets
pg-migrations-up:
	migrate -path internal/infrastructure/pg_migrations -database "postgres://postgres:postgres@localhost:5432/checker?sslmode=disable" up

pg-migrations-down:
	migrate -path internal/infrastructure/pg_migrations -database "postgres://postgres:postgres@localhost:5432/checker?sslmode=disable" down

# Docker compose related targets
docker-compose-local-up:
	# Initialize Swagger documentation
	PROJECT_ROOT=${root_dir} swag init -g cmd/server.go

	# Start Docker Compose locally
	PROJECT_ROOT=${root_dir} docker-compose -f ${docker_compose_local_conf_file} \
		-p $(PROJECT_NAME) --env-file ${root_dir}/.env up \
		--no-deps --build --remove-orphans --pull never

docker-compose-up:
	# Start Docker Compose with main configuration
	PROJECT_ROOT=${root_dir} docker-compose -f ${docker_compose_conf_file} \
		-p $(PROJECT_NAME) --env-file ${root_dir}/.env up \
		--no-deps --build --remove-orphans --pull never

docker-compose-up-detached:
	# Start Docker Compose in detached mode
	PROJECT_ROOT=${root_dir} docker-compose -f ${docker_compose_conf_file} \
		-p $(PROJECT_NAME) --env-file ${root_dir}/.env up \
		-d --no-deps --build --remove-orphans --pull never

docker-compose-stop:
	# Stop the running containers
	PROJECT_ROOT=${root_dir} docker-compose -f ${docker_compose_conf_file} \
		-p $(PROJECT_NAME) --env-file ${root_dir}/.env stop

docker-compose-down-volumes:
	# Stop containers and remove volumes
	PROJECT_ROOT=${root_dir} docker-compose -f ${docker_compose_conf_file} \
		--env-file ${root_dir}/.env down \
		--volumes --remove-orphans

# DB migration management inside the container
db-migrations-up:
	# Apply migrations inside the Docker container
	docker exec -it $(PROJECT_NAME)-database \
		psql -U postgres -d $(POSTGRES_DB) -f /database/upgrade.sql

db-migrations-down:
	# Rollback migrations inside the Docker container
	docker exec -it $(PROJECT_NAME)-database \
		psql -U postgres -d $(POSTGRES_DB) -f /database/downgrade.sql

# Docker build and run
docker-build:
	# Build the Docker image
	docker build -t $(PROJECT_NAME) .

docker-run:
	# Run the Docker container
	docker run -d -p 8080:8080 --name $(PROJECT_NAME) $(PROJECT_NAME)

# Submodule update (if you have submodules)
submodules-update:
	# Update git submodules
	git submodule update --recursive --remote

# Deploy target
deploy:
	# Deploy the app using Docker Compose in detached mode
	$(MAKE) docker-compose-up-detached
