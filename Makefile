include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d todoapp-postgres

env-cleanup:
	@read -p "Do you want to delete all volume files? This will delete the database! [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgres port_forwarder && \
		echo "Files successfully cleaned"; \
	else \
		echo "Clean files canceled"; \
	fi

env-port-forward:
	@docker compose up -d port_forwarder

env-port-close:
	@docker compose down port_forwarder


migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Missing required parameter: seq. Example: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	mkdir -p migrations; \
	docker compose run --rm \
		--user "$(shell id -u):$(shell id -g)" \
		todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Missing required parameter action. Example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
		"$(action)"

logs-cleanup:
	@read -p "Do you want to delete all logs files? This will delete all logs! [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Logs files successfully cleaned"; \
	else \
		echo "Clean logs files canceled"; \
	fi
	
todoapp-run:
	@export LOGGER_FOLDER="${PROJECT_ROOT}/out/logs" && \
	go mod tidy && \
	go run cmd/todoapp/main.go