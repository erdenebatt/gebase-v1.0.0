# Docker Configuration

## docker-compose.yml

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: gebase_postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: gerege_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./docker/init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - gebase_network

  redis:
    image: redis:7-alpine
    container_name: gebase_redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - gebase_network

  backend:
    build:
      context: ./backend
      dockerfile: docker/Dockerfile
    container_name: gebase_backend
    ports:
      - "8000:8000"
    environment:
      - SERVER_PORT=8000
      - SERVER_MODE=development
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=gerege_db
      - DB_SCHEMA=gerege_base
      - DB_SSLMODE=disable
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - JWT_SECRET=your-super-secret-key-change-in-production
      - CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - gebase_network
    restart: unless-stopped

  admin:
    build:
      context: ./frontend/admin
      dockerfile: Dockerfile
    container_name: gebase_admin
    ports:
      - "3001:3001"
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8000
      - NEXT_PUBLIC_APP_NAME=Gebase Admin
    depends_on:
      - backend
    networks:
      - gebase_network
    restart: unless-stopped

  portal:
    build:
      context: ./frontend/portal
      dockerfile: Dockerfile
    container_name: gebase_portal
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8000
      - NEXT_PUBLIC_APP_NAME=Gebase Portal
    depends_on:
      - backend
    networks:
      - gebase_network
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:

networks:
  gebase_network:
    driver: bridge
```

## docker/init.sql

```sql
-- Create schema
CREATE SCHEMA IF NOT EXISTS gerege_base;

-- Set default search path
ALTER DATABASE gerege_db SET search_path TO gerege_base, public;

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
```

## Backend Dockerfile (backend/docker/Dockerfile)

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/server

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/server .
COPY --from=builder /app/docs ./docs

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8000

ENTRYPOINT ["./server"]
```

## Admin Frontend Dockerfile (frontend/admin/Dockerfile)

```dockerfile
# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci

# Copy source code
COPY . .

# Build
ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build

# Runtime stage
FROM node:20-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# Create non-root user
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copy built files
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

USER nextjs

EXPOSE 3001

ENV PORT=3001
ENV HOSTNAME="0.0.0.0"

CMD ["node", "server.js"]
```

## Portal Frontend Dockerfile (frontend/portal/Dockerfile)

```dockerfile
# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci

# Copy source code
COPY . .

# Build
ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build

# Runtime stage
FROM node:20-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# Create non-root user
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copy built files
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

CMD ["node", "server.js"]
```

## .dockerignore (for all services)

```
# Dependencies
node_modules
vendor

# Build outputs
.next
out
dist
build

# Development
.env*.local
.git
.gitignore

# IDE
.idea
.vscode
*.swp
*.swo

# Testing
coverage
*.test.js

# Misc
README.md
*.md
Makefile
```

## Makefile (Project root)

```makefile
.PHONY: help build up down logs restart clean

# Default target
help:
	@echo "Gebase Platform Commands:"
	@echo "  make build    - Build all Docker images"
	@echo "  make up       - Start all services"
	@echo "  make down     - Stop all services"
	@echo "  make logs     - View all logs"
	@echo "  make restart  - Restart all services"
	@echo "  make clean    - Remove all containers and volumes"
	@echo "  make dev      - Start development environment"
	@echo "  make migrate  - Run database migrations"
	@echo "  make seed     - Seed database with initial data"

# Build all images
build:
	docker-compose build

# Start all services
up:
	docker-compose up -d

# Stop all services
down:
	docker-compose down

# View logs
logs:
	docker-compose logs -f

# Restart services
restart:
	docker-compose restart

# Clean everything
clean:
	docker-compose down -v --rmi all --remove-orphans

# Development environment (with hot reload)
dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

# Run migrations
migrate:
	docker-compose exec backend ./server migrate

# Seed database
seed:
	docker-compose exec backend ./server seed

# Backend specific
backend-shell:
	docker-compose exec backend sh

# Database shell
db-shell:
	docker-compose exec postgres psql -U postgres -d gerege_db

# Redis shell
redis-shell:
	docker-compose exec redis redis-cli
```

## docker-compose.dev.yml (Development overrides)

```yaml
version: '3.8'

services:
  backend:
    build:
      context: ./backend
      dockerfile: docker/Dockerfile.dev
    volumes:
      - ./backend:/app
      - /app/vendor
    environment:
      - SERVER_MODE=development
      - GIN_MODE=debug

  admin:
    build:
      context: ./frontend/admin
      dockerfile: Dockerfile.dev
    volumes:
      - ./frontend/admin:/app
      - /app/node_modules
      - /app/.next
    environment:
      - NODE_ENV=development

  portal:
    build:
      context: ./frontend/portal
      dockerfile: Dockerfile.dev
    volumes:
      - ./frontend/portal:/app
      - /app/node_modules
      - /app/.next
    environment:
      - NODE_ENV=development
```

## Backend Development Dockerfile (backend/docker/Dockerfile.dev)

```dockerfile
FROM golang:1.22-alpine

WORKDIR /app

# Install air for hot reload
RUN go install github.com/cosmtrek/air@latest

# Install dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8000

CMD ["air", "-c", ".air.toml"]
```

## Air config (.air.toml) for backend hot reload

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/server"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "docs"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_error = true

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
```

## Frontend Development Dockerfile (frontend/admin/Dockerfile.dev)

```dockerfile
FROM node:20-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

EXPOSE 3001

ENV PORT=3001

CMD ["npm", "run", "dev"]
```

## GitHub Actions CI/CD (.github/workflows/ci.yml)

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Run tests
        working-directory: ./backend
        run: |
          go mod download
          go test -v ./...
      
      - name: Run linter
        uses: golangci/golangci-lint-action@v4
        with:
          working-directory: ./backend

  test-frontend:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        app: [admin, portal]
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Node
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: frontend/${{ matrix.app }}/package-lock.json
      
      - name: Install dependencies
        working-directory: ./frontend/${{ matrix.app }}
        run: npm ci
      
      - name: Run linter
        working-directory: ./frontend/${{ matrix.app }}
        run: npm run lint
      
      - name: Build
        working-directory: ./frontend/${{ matrix.app }}
        run: npm run build

  build-and-push:
    needs: [test-backend, test-frontend]
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      
      - name: Login to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Build and push backend
        uses: docker/build-push-action@v5
        with:
          context: ./backend
          file: ./backend/docker/Dockerfile
          push: true
          tags: ghcr.io/${{ github.repository }}/backend:latest
      
      - name: Build and push admin
        uses: docker/build-push-action@v5
        with:
          context: ./frontend/admin
          push: true
          tags: ghcr.io/${{ github.repository }}/admin:latest
      
      - name: Build and push portal
        uses: docker/build-push-action@v5
        with:
          context: ./frontend/portal
          push: true
          tags: ghcr.io/${{ github.repository }}/portal:latest
```
