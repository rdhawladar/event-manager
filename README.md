# Event Manager API

A RESTful API built with Go for managing events, featuring PostgreSQL database integration, Docker support, and Swagger documentation.

## 🛠 Technical Stack

- **Backend**: Go (Golang)
- **Database**: PostgreSQL
- **Containerization**: Docker & Docker Compose
- **Documentation**: Swagger/OpenAPI
- **Development Tools**:
  - Air (for hot-reloading)
  - Make (for build automation)

## 📋 Prerequisites

- Go 1.22 or higher
- Docker and Docker Compose
- Make (optional, for using Makefile commands)

## 🚀 Quick Start

### Using Docker (Recommended)

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd event-manager
   ```

2. Build and run using Docker Compose:
   ```bash
   make build
   make run
   ```

   Or run in detached mode:
   ```bash
   make run-detached
   ```

3. The API will be available at `http://localhost:8081`
4. Swagger documentation is available at `http://localhost:8081/swagger/`

### Local Development

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Install Air for hot-reloading:
   ```bash
   make install-air
   ```

3. Run the application with hot-reloading:
   ```bash
   make watch-local
   ```

## 📚 API Documentation

The API documentation is available through Swagger UI at `http://localhost:8081/swagger/`. You can also generate the documentation locally:

```bash
make swagger
```

### Available Endpoints

- `GET /events` - List all events (paginated)
- `POST /events` - Create a new event
- `GET /events/{id}` - Get a specific event
- `PUT /events/{id}` - Update an event
- `DELETE /events/{id}` - Delete an event

## 🗄 Database Schema

```sql
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    location VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

## 🔧 Development Commands

- `make build` - Build the application and Docker images
- `make run` - Run the application in containers
- `make run-detached` - Run the application in background
- `make stop` - Stop all containers
- `make clean` - Clean up containers, images, and volumes
- `make logs` - Show logs from all containers
- `make service-logs service=<service>` - Show logs from specific service
- `make db-shell` - Access PostgreSQL shell
- `make restart service=<service>` - Rebuild and restart a specific service
- `make swagger` - Generate Swagger documentation
- `make watch-local` - Run with hot-reloading locally

## 📝 API Request Examples

### Create Event
```bash
curl -X POST http://localhost:8081/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Team Meeting",
    "description": "Weekly team sync",
    "location": "Conference Room A",
    "start_time": "2024-03-20T10:00:00Z",
    "end_time": "2024-03-20T11:00:00Z",
    "created_by": "john.doe@example.com"
  }'
```

### List Events
```bash
curl http://localhost:8081/events?page=1&page_size=10
```

### Get Event
```bash
curl http://localhost:8081/events/1
```

### Update Event
```bash
curl -X PUT http://localhost:8081/events/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Meeting",
    "description": "Updated description",
    "location": "New Location",
    "start_time": "2024-03-21T10:00:00Z",
    "end_time": "2024-03-21T11:00:00Z",
    "created_by": "john.doe@example.com"
  }'
```

### Delete Event
```bash
curl -X DELETE http://localhost:8081/events/1
```

## 🔍 Validation Rules

The API includes validation for event creation and updates:

- Title: Required, non-empty string
- Location: Required, non-empty string
- Created By: Required, valid email address
- Start Time: Required, must be in the future
- End Time: Required, must be after start time

## 🐳 Docker Configuration

The application uses Docker Compose with two services:

1. **App Service**:
   - Go application
   - Hot-reloading with Air
   - Port 8081 exposed
   - Volume mounts for development

2. **Database Service**:
   - PostgreSQL 15 (Alpine)
   - Port 5432 exposed
   - Persistent volume for data
   - Health checks configured

## 🔐 Environment Variables

The following environment variables can be configured:

- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (default: postgres)
- `DB_PASSWORD`: Database password (default: postgres)
- `DB_NAME`: Database name (default: event_management)

## 📦 Project Structure

```
event-manager/
├── main.go           # Main application file
├── Dockerfile        # Docker configuration
├── docker-compose.yml # Docker Compose configuration
├── go.mod           # Go module file
├── go.sum           # Go module checksum
├── Makefile         # Build automation
└── docs/            # Swagger documentation
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.