# Cortes Surplus Inventory Management System - Backend

This is the backend API for the Cortes Surplus Inventory Management System, built with Go and Fiber.

## Prerequisites

- Go 1.24 or higher
- MySQL 8.0 or higher

## Setup

1. Clone the repository
2. Create a `.env` file in the root directory based on `.env.example`
3. Set up the database:
   ```bash
   mysql -u your_username -p your_database < schema.sql
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```

## Running the Application

```bash
go run cmd/web/main.go
```

The server will start on port 8080 by default.

## API Endpoints

### User Management

- `POST /api/users/register` - Register a new user
- `POST /api/users/login` - Login
- `GET /api/users` - Get all users (requires authentication)
- `GET /api/users/:id` - Get a specific user (requires authentication)
- `PUT /api/users/:id` - Update a user (requires authentication)
- `DELETE /api/users/:id` - Delete a user (requires authentication)
- `PUT /api/users/:id/activate` - Activate a user (requires authentication)
- `PUT /api/users/:id/deactivate` - Deactivate a user (requires authentication)
- `PUT /api/users/:id/password` - Update a user's password (requires authentication)

## Default Admin User

Email: admin@example.com
Password: admin123

## Development

### Project Structure

- `cmd/web/` - Application entry point
- `internal/` - Internal packages
  - `config/` - Configuration
  - `handlers/` - HTTP handlers
  - `models/` - Data models
  - `repositories/` - Database operations

### Running with Docker

Build the Docker image:

```bash
docker build -f Backend/Dockerfile -t backend-dev .
```

Run the container (host-gateway for macOS/Windows):

```bash
docker run -it --rm \
  --add-host host.docker.internal:host-gateway \
  -p 8080:8080 \
  --env-file Backend/.env \
  -e PORT=8080 \
  -v "$(pwd)/Backend":/app \
  backend-dev
```

Or, on Linux with host networking:

```bash
docker run -it --rm \
  --network host \
  -e DB_HOST=127.0.0.1 \
  -e DB_PORT=3306 \
  -e DB_USERNAME=root \
  -e DB_PASSWORD="YOUR_DB_PASSWORD" \
  -e DB_NAME=oop \
  -e JWT_SECRET="your_jwt_secret" \
  -e FRONTEND_URL="http://localhost:9000" \
  -e PORT=8080 \
  -v "$(pwd)/Backend":/app \
  backend-dev
```
