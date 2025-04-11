ENV_FILE = .env
EXAMPLE_ENV_FILE = .env.example

env:
	@echo "Создание .env из .env.example"
	@cp $(EXAMPLE_ENV_FILE) $(ENV_FILE)
	@echo "Если нужно, отредактируйте файл .env перед запуском."

up:
	@echo "Запуск контейнеров в фоновом режиме..."
	@docker-compose -f deployments/docker-compose.yaml --env-file .env up -d

run:
	@echo "Запуск приложения в контейнере..."
	@docker-compose -f deployments/docker-compose.yaml exec social-app ./social_binary

del:
	@echo "Удаляем контейнеры..."
	@docker-compose -f deployments/docker-compose.yaml down -v

dev-up:
	@docker-compose -f deployments/docker-compose.yaml --env-file .env up -d postgres migrate

dev-down:
	@docker-compose -f deployments/docker-compose.yaml down

dev-run:
	@go run cmd/social/main.go