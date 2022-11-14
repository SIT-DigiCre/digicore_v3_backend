include .env

.PHONY: generate_api
generate_api:
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app/document node_tool swagger-cli bundle -o ./bundle.yml -t yaml ./openapi.yml
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app admin oapi-codegen --config ./config/models.yml ./document/bundle.yml
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app admin oapi-codegen --config ./config/server.yml ./document/bundle.yml
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app admin oapi-codegen --config ./config/spec.yml ./document/bundle.yml

.PHONY: migrate-dry
migrate-dry:
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app/schema admin bash -c 'cat `ls | grep -v foreign_key.sql` | go run github.com/k0kubun/sqldef/cmd/mysqldef@v0.13.19 --user=${DB_USER} --password=${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE} --dry-run'

.PHONY: migrate
migrate:
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app/schema admin bash -c 'cat `ls | grep -v foreign_key.sql` | go run github.com/k0kubun/sqldef/cmd/mysqldef@v0.13.19 --user=${DB_USER} --password=${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE}'
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app/schema admin bash -c 'mysql -u ${DB_USER} -p${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE} -vvv < foreign_key.sql'

.PHONY: up
up:
	docker compose -f ${DOCKER_COMPOSE} up

.PHONY: up-d
up-d:
	docker compose -f ${DOCKER_COMPOSE} up -d

.PHONY: down
down:
	docker compose -f ${DOCKER_COMPOSE} down

.PHONY: pull
pull:
	docker compose -f ${DOCKER_COMPOSE} pull

.PHONY: build
build:
	docker compose -f ${DOCKER_COMPOSE} build
