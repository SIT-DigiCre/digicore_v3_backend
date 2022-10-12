.PHONY: generate_api
generate_api:
	docker compose run --rm -w /app/document node_tool swagger-cli bundle -o ./bundle.yml -t yaml ./openapi.yml
	docker compose run --rm -w /app admin oapi-codegen --config ./config/models.yml ./document/bundle.yml
	docker compose run --rm -w /app admin oapi-codegen --config ./config/server.yml ./document/bundle.yml
	docker compose run --rm -w /app admin oapi-codegen --config ./config/spec.yml ./document/bundle.yml
