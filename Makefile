root_dir := $(shell pwd)
docker_compose_conf_file := ${root_dir}/docker/docker-compose.yaml
docker_compose_local_conf_file := ${root_dir}/docker/docker-compose.local.yaml

ifneq ("$(wildcard .env)","")
	include .env
endif

docker-compose-local-up:
	PROJECT_ROOT=${root_dir} swag init -g cmd/main.go
	PROJECT_ROOT=${root_dir} docker compose -f ${docker_compose_local_conf_file} \
		-p $(PROJECT_NAME) --env-file ${root_dir}/.env up \
		--no-deps --build --remove-orphans --pull never

docker-compose-up:
	PROJECT_ROOT=${root_dir} docker compose -f ${docker_compose_conf_file} \
		-p $(PROJECT_NAME) --env-file ${root_dir}/.env up \
		--no-deps --build --remove-orphans --pull never

docker-compose-up-detached:
	PROJECT_ROOT=${root_dir} docker compose -f ${docker_compose_conf_file} \
		-p $(PROJECT_NAME) --env-file ${root_dir}/.env up \
		-d --no-deps --build --remove-orphans --pull never

docker-compose-stop:
	PROJECT_ROOT=${root_dir} docker compose -f ${docker_compose_conf_file} \
		-p $(PROJECT_NAME) --env-file ${root_dir}/.env stop

docker-compose-down-volumes:
	PROJECT_ROOT=${root_dir} docker compose -f ${docker_compose_conf_file} \
		--env-file ${root_dir}/.env down \
		--volumes --remove-orphans

db-migrations-up:
	docker exec -it $(PROJECT_NAME)-database \
		psql -U postgres -d $(POSTGRES_DB) -f /database/upgrade.sql

db-migrations-down:
	docker exec -it $(PROJECT_NAME)-database \
		psql -U postgres -d $(POSTGRES_DB) -f /database/downgrade.sql

deploy:
	$(MAKE) docker-compose-up-detached
