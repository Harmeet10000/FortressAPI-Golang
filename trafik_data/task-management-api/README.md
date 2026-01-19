# Task Management API

A professional, production-ready Task Management REST API built with Go, featuring a clean architecture with dependency injection, comprehensive error handling, and industry best practices.

## 🏗️ Architecture

This project follows a **feature-based multilayered architecture**:

```
cmd/api/               # Application entry point
internal/
  ├── app/            # Server setup and initialization
  ├── config/         # Configuration management (Koanf)
  ├── connections/    # Database and Redis connections
  ├── database/       # SQL schema and queries (SQLC)
  ├── errs/           # Custom error types
  ├── feature/        # Feature modules
  │   ├── category/   # Category feature
  │   ├── todo/       # Todo feature
  │   └── comment/    # Comment feature
  ├── logger/         # Structured logging (Zerolog)
  ├── middlewares/    # HTTP middlewares
  └── validation/     # Request validation utilities
```

Each feature module contains:
- **Model**: Domain entities
- **DTO**: Data Transfer Objects (request/response)
- **Repository**: Data access layer
- **Service**: Business logic layer
- **Handler**: HTTP request handlers
- **Routes**: Route registration

## 🚀 Tech Stack

- **Web Framework**: Echo v4
- **Configuration**: Koanf
- **Database**: PostgreSQL with pgx/v5
- **Migrations**: Goose
- **Query Builder**: SQLC (type-safe SQL)
- **Caching**: Redis
- **Logging**: Zerolog
- **Background Jobs**: Asynq
- **Email**: Resend
- **Documentation**: Swagger/OpenAPI

## 📋 Prerequisites

- Go 1.22+
- PostgreSQL 14+
- Redis 7+
- Docker & Docker Compose (optional)

## 🛠️ Setup

### 1. Clone the repository

```bash
git clone <repository-url>
cd task-management-api
```

### 2. Install dependencies

```bash
make deps
```

### 3. Install required tools

```bash
make install-tools
```

### 4. Start infrastructure (Docker)

```bash
make docker-up
```

Or manually:
```bash
docker-compose up -d
```

### 5. Run migrations

```bash
make migrate-up
```

### 6. Generate SQLC code

```bash
make sqlc
```

### 7. Generate Swagger documentation

```bash
make swagger
```

### 8. Run the application

```bash
make run
```

Or build and run:
```bash
make build
./bin/api
```

## 🌐 API Endpoints

### Categories

- `POST /api/v1/categories` - Create a category
- `GET /api/v1/categories` - List categories
- `GET /api/v1/categories/:id` - Get category by ID
- `PUT /api/v1/categories/:id` - Update category
- `DELETE /api/v1/categories/:id` - Delete category

### Todos

- `POST /api/v1/todos` - Create a todo
- `GET /api/v1/todos` - List todos
- `GET /api/v1/todos/:id` - Get todo by ID
- `PUT /api/v1/todos/:id` - Update todo
- `PATCH /api/v1/todos/:id/status` - Update todo status
- `DELETE /api/v1/todos/:id` - Delete todo

### Comments

- `POST /api/v1/comments` - Create a comment
- `GET /api/v1/todos/:todoId/comments` - List comments for a todo
- `GET /api/v1/comments/:id` - Get comment by ID
- `PUT /api/v1/comments/:id` - Update comment
- `DELETE /api/v1/comments/:id` - Delete comment

### Health

- `GET /health` - Health check endpoint

### Documentation

- `GET /swagger/*` - Swagger UI

## 📝 Configuration

Configuration is managed through `config.yaml` and can be overridden with environment variables:

```yaml
server:
  host: 0.0.0.0
  port: 8080

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: task_management

redis:
  host: localhost
  port: 6379
```

Environment variables use the `APP_` prefix:
```bash
export APP_SERVER_PORT=9000
export APP_DATABASE_HOST=db.example.com
```

## 🧪 Testing

```bash
make test
```

## 📦 Building

```bash
make build
```

The binary will be created at `bin/api`.

## 🔄 Database Migrations

```bash
# Run migrations
make migrate-up

# Rollback last migration
make migrate-down

# Reset all migrations
make migrate-reset
```

## 📚 API Documentation

After running `make swagger`, access the Swagger UI at:
```
http://localhost:8080/swagger/index.html
```

## 🏃 Quick Start Development

```bash
make setup  # Complete setup
make dev    # Start development environment
```

## 🛡️ Features

- ✅ Clean Architecture with Dependency Injection
- ✅ Feature-based organization
- ✅ Type-safe SQL queries with SQLC
- ✅ Comprehensive error handling
- ✅ Request validation
- ✅ Structured logging
- ✅ Database migrations
- ✅ Graceful shutdown
- ✅ Health checks
- ✅ API documentation
- ✅ Docker support
- ✅ Makefile for common tasks

## 📄 License

MIT License

## 👥 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
