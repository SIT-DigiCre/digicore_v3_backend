.PHONY: generate_api
generate_api:
	docker compose run --rm -w /app/document node_tool swagger-cli bundle -o ./bundle.yml -t yaml ./openapi.yml
	docker compose run --rm -w /app admin oapi-codegen --config ./config/models.yml ./document/bundle.yml
	docker compose run --rm -w /app admin oapi-codegen --config ./config/server.yml ./document/bundle.yml
	docker compose run --rm -w /app admin oapi-codegen --config ./config/spec.yml ./document/bundle.yml

migrate-dry:
	docker compose run --rm -w /app/schema admin bash -c 'cat `ls | grep -v foreign_key.sql` | go run github.com/k0kubun/sqldef/cmd/mysqldef@v0.13.9 --user=$${DB_USER} --password=$${DB_PASSWORD} --host=$${DB_HOST} $${DB_DATABASE} --dry-run'

migrate:
	docker compose run --rm -w /app/schema admin bash -c 'cat `ls | grep -v foreign_key.sql` | go run github.com/k0kubun/sqldef/cmd/mysqldef@v0.13.9 --user=$${DB_USER} --password=$${DB_PASSWORD} --host=$${DB_HOST} $${DB_DATABASE}'
	docker compose run --rm -w /app/schema admin bash -c 'mysql -u $${DB_USER} -p$${DB_PASSWORD} --host=$${DB_HOST} $${DB_DATABASE} -vvv < foreign_key.sql'
