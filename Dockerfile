FROM golang:1.22.1-alpine

WORKDIR /app

# Install air for hot-reloading
RUN go install github.com/cosmtrek/air@v1.49.0

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Expose the port
EXPOSE 8081

# Use air for hot-reloading
CMD ["air"] 