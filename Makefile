INFRA_COMPOSE := -f build/docker/local/docker-compose.yml

.PHONY: help setup down infra-up infra-down app-up app-down logs aws-init-logs test lint

up: 
	@echo "INFO: Starting environment..."
	@docker-compose $(INFRA_COMPOSE) up -d

down: 
	@echo "INFO: Shutting down environment..."
	@docker-compose $(INFRA_COMPOSE) down -v --remove-orphans
	@echo "INFO: Environment shut down successfully."

logs:
	@echo "INFO: Following environment logs... (Press Ctrl+C to exit)"
	@docker-compose $(APP_COMPOSE) logs -f

test: 
	@echo "INFO: Running tests..."
	@go test -v ./...

lint: ## 💅 Runs the linter to check code quality (requires golangci-lint installed).
	@echo "INFO: Running linter..."
	@golangci-lint run

lint-fix: ## 🛠️ Runs the linter and automatically fixes issues (requires golangci-lint installed).
	@echo "INFO: Running linter with auto-fix..."
	@golangci-lint run --fix
