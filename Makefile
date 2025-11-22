include .env

.PHONY: generate_api
generate_api:
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app/document node_tool swagger-cli bundle -o ./bundle.gen.yml -t yaml ./openapi.yml
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app admin oapi-codegen --config ./config/models.yml ./document/bundle.gen.yml
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app admin oapi-codegen --config ./config/server.yml ./document/bundle.gen.yml
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app admin oapi-codegen --config ./config/spec.yml ./document/bundle.gen.yml

.PHONY: migrate-dry
migrate-dry:
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app/schema admin bash -c 'cat `ls -F | grep -v / | sed s/\*$$//g | grep -v foreign_key.sql` | go run github.com/k0kubun/sqldef/cmd/mysqldef@v0.15.12 --user=${DB_USER} --password=${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE} --dry-run'

.PHONY: migrate
migrate:
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app/schema admin bash -c 'cat `ls -F | grep -v / | sed s/\*$$//g | grep -v foreign_key.sql` | go run github.com/k0kubun/sqldef/cmd/mysqldef@v0.15.12 --user=${DB_USER} --password=${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE}'
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app/schema admin bash -c 'cat foreign_key.sql | mysql -u ${DB_USER} -p${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE} -vvv'

.PHONY: insert_test
insert_test:
	docker compose -f ${DOCKER_COMPOSE} run --rm -w /app/schema/test admin bash -c 'cat `ls -F | grep -v / | sed s/\*$$//g` | mysql -u ${DB_USER} -p${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE} --default-character-set=utf8mb4 -vvv'


.PHONY: up
up:
	docker compose -f ${DOCKER_COMPOSE} up

.PHONY: up-d
up-d:
	docker compose -f ${DOCKER_COMPOSE} up -d

.PHONY: logs
logs:
	docker compose -f ${DOCKER_COMPOSE} logs --since $(date +%Y-%m-%d --date '1 day ago')

.PHONY: logs-all
logs-all:
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
ls:
	docker compose -f ${DOCKER_COMPOSE} ls
