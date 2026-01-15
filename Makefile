.PHONY: help build up down logs restart clean migrate seed dev

help:
	@echo "Gebase Platform Commands:"
	@echo "  make build    - Build all Docker images"
	@echo "  make up       - Start all services"
	@echo "  make down     - Stop all services"
	@echo "  make logs     - View all logs"
	@echo "  make restart  - Restart all services"
	@echo "  make clean    - Remove all containers and volumes"
	@echo "  make migrate  - Run database migrations"
	@echo "  make seed     - Seed database with initial data"
	@echo "  make dev      - Start development environment (postgres + redis only)"
	@echo "  make backend  - Run backend locally"

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

restart:
	docker-compose restart

clean:
	docker-compose down -v --rmi all --remove-orphans

dev:
	docker-compose up -d postgres redis

migrate:
	docker-compose exec backend ./server migrate

seed:
	docker-compose exec backend ./server seed

backend:
	cd backend && go run ./cmd/server

db-shell:
	docker-compose exec postgres psql -U postgres -d gerege_db

redis-shell:
	docker-compose exec redis redis-cli
