DEFAULT_GOAL: up

.PHONY: up down

up:
	docker compose up -d

down:
	docker compose down
