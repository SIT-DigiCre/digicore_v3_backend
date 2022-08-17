.PHONY: generate_api
generate_api: ./document/openapi.yml
	oapi-codegen --config ./config/models.yml ./document/openapi.yml
	oapi-codegen --config ./config/server.yml ./document/openapi.yml
	oapi-codegen --config ./config/spec.yml ./document/openapi.yml
