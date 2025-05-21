include .env

MIGRATIONS_PATH = ./migrations

.PHONY: generate-server
generate-server:
	@echo "generating server from OpenAPI specification"
	@go tool oapi-codegen --config cfg-api.yaml ./api/fumode_bundled_openapi.yaml

.PHONY: bundle-openapi
bundle-openapi:
	@redocly bundle ./api/fumode/openapi.yaml --output ./api/fumode_bundled_openapi.yaml --remove-unused-components

.PHONY: migrate-create
migrate-create:
	@if [-z "$(name)"]; then \
		echo "Usage: make migrate-create name=<filename>" \
		exit 1; \
	fi
	@echo "Creating migration files for $(name)"
	@migrate create -seq -ext .sql -dir ${MIGRATIONS_PATH} $(name)

.PHONY: migrate-up
migrate-up:
	@echo "Running up migrations"
	@migrate -path ${MIGRATIONS_PATH} -database ${FUMODE_DB_DSN} up

.PHONY: migrate-down
migrate-down:
	@echo "Running down migrations"
	@migrate -path ${MIGRATIONS_PATH} -database ${FUMODE_DB_DSN} down

.PHONY: migrate-drop
migrate-drop:
	@echo "Dropping all tables"
	@migrate -path ${MIGRATIONS_PATH} -database ${FUMODE_DB_DSN} drop