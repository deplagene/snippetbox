run:
	@go run ./cmd/

migrate:
	@go run ./cmd/ --migrate

up:
	@docker compose up -d

down:
	@docker compose down