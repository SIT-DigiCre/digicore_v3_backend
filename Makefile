include .env

.PHONY: generate_api
generate_api:
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app/document node_tool swagger-cli bundle -o ./bundle.yml -t yaml ./openapi.yml
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app/document node_tool swagger-cli bundle -o ./bundle-develop.yml -t yaml ./openapi-develop.yml
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app admin oapi-codegen --config ./config/models.yml ./document/bundle.yml
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app admin oapi-codegen --config ./config/server.yml ./document/bundle.yml
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app admin oapi-codegen --config ./config/spec.yml ./document/bundle.yml

.PHONY: migrate-dry
migrate-dry:
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app/schema admin bash -c 'cat `ls -F | grep -v / | sed s/\*$$//g | grep -v foreign_key.sql` | go run github.com/k0kubun/sqldef/cmd/mysqldef@v0.15.12 --user=${DB_USER} --password=${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE} --dry-run'

.PHONY: migrate
migrate:
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app/schema admin bash -c 'cat `ls -F | grep -v / | sed s/\*$$//g | grep -v foreign_key.sql` | go run github.com/k0kubun/sqldef/cmd/mysqldef@v0.15.12 --user=${DB_USER} --password=${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE}'
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app/schema admin bash -c 'cat foreign_key.sql | mysql -u ${DB_USER} -p${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE} -vvv'

.PHONY: migrate
insert_test:
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app/schema/test admin bash -c 'cat `ls -F | grep -v / | sed s/\*$$//g | mysql -u ${DB_USER} -p${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE} -vvv'


.PHONY: up
up:
	docker compose -f ${DOCKER_COMPOSE} up

.PHONY: up-d
up-d:
	docker compose -f ${DOCKER_COMPOSE} up -d

.PHONY: logs
up-d:
	docker compose -f ${DOCKER_COMPOSE} logs --since $(date +%Y-%m-%d --date '1 day ago')

.PHONY: logs-all
up-d:
	docker compose -f ${DOCKER_COMPOSE} logs

.PHONY: down
down:
	docker compose -f ${DOCKER_COMPOSE} down

.PHONY: pull
pull:
	docker compose -f ${DOCKER_COMPOSE} pull

.PHONY: build
build:
	docker compose -f ${DOCKER_COMPOSE} build

.PHONY: ls
build:
	docker compose -f ${DOCKER_COMPOSE} ls
