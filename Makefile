.PHONY: build run stop clean test logs db-shell swagger watch watch-docker

# Build the application and Docker images
build:
	go mod tidy
	docker-compose build

# Run the application in containers
run:
	docker-compose up

# Run the application in background as detacheed
run-d:
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

# Generate Swagger documentation
swagger:
	~/go/bin/swag init

# Install air for hot-reloading
install-air:
	go install github.com/cosmtrek/air@v1.49.0

# Watch and reload the application locally
watch:
	$(shell go env GOPATH)/bin/air

# Default target
all: build run 