.PHONY: build run stop clean test logs db-shell

# Build the application and Docker images
build:
	go mod tidy
	docker-compose build

# Run the application in containers
run:
	docker-compose up

# Run the application in background
run-detached:
	docker-compose up -d

# Stop all containers
stop:
	docker-compose down

# Clean up containers, images, and volumes
clean:
	docker-compose down -v
	docker system prune -f

# Show logs from all containers
logs:
	docker-compose logs -f

# Show logs from specific service (usage: make service-logs service=app)
service-logs:
	docker-compose logs -f $(service)

# Access PostgreSQL shell
db-shell:
	docker-compose exec db psql -U postgres -d event_management

# Rebuild and restart a specific service (usage: make restart service=app)
restart:
	docker-compose up -d --build $(service)

# Default target
all: build run 