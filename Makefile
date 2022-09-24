.PHONY: generate_api
generate_api: ./document/bundle.yml
	oapi-codegen --config ./config/models.yml ./document/bundle.yml
	oapi-codegen --config ./config/server.yml ./document/bundle.yml
	oapi-codegen --config ./config/spec.yml ./document/bundle.yml
