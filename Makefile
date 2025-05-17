.PHONY: generate-server
generate-server:
	@echo "generating server from OpenAPI specification"
	@go tool oapi-codegen --config cfg-api.yaml ./api/fumode/openapi.yaml

.PHONY: bundle-openapi
bundle-openapi:
	@redocly bundle ./api/fumode/openapi.yaml --output ./api/fumode_bundled_openapi.yaml --remove-unused-components