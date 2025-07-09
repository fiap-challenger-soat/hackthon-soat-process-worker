INFRA_COMPOSE := -f build/docker/local/docker-compose.yml

.PHONY: help setup down infra-up infra-down app-up app-down logs aws-init-logs test lint

up: 
	@echo "INFO: Starting environment..."
	@docker-compose $(INFRA_COMPOSE) up -d

down: 
	@echo "INFO: Shutting down environment..."
	@docker-compose $(INFRA_COMPOSE) down -v --remove-orphans
	@echo "INFO: Environment shut down successfully."

test: 
	@echo "INFO: Running tests..."
	@go test -cover -coverprofile=coverage.out `go list ./... | grep -v mocks | grep -v cmd`

cov: ## 📊 Generates a coverage report.
	@echo "INFO: Generating coverage report..."
	@go tool cover -html=coverage.out
	@echo "INFO: Coverage report generated at coverage.html"

lint: ## 💅 Runs the linter to check code quality (requires golangci-lint installed).
	@echo "INFO: Running linter..."
	@golangci-lint run

lint-fix: ## 🛠️ Runs the linter and automatically fixes issues (requires golangci-lint installed).
	@echo "INFO: Running linter with auto-fix..."
	@golangci-lint run --fix

gen-mock: ## 🛠️ Generates mock files for testing.
	@echo "INFO: Generating mock files..."
	@go generate ./...
