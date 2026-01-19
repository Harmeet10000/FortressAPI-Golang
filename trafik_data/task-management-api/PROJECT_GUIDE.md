# Task Management API - Complete Implementation Guide

## 🎉 Project Overview

This is a **production-ready** Go backend application for task management with:

- ✅ **Clean Architecture** with feature-based organization
- ✅ **Dependency Injection** throughout all layers
- ✅ **Type-safe SQL** with SQLC
- ✅ **Graceful shutdown** handling
- ✅ **Comprehensive error handling**
- ✅ **Request validation**
- ✅ **Structured logging** with Zerolog
- ✅ **API documentation** with Swagger
- ✅ **Database migrations** with Goose
- ✅ **Docker support** for local development

## 📁 Project Structure

```
task-management-api/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point with DI
├── internal/
│   ├── app/
│   │   └── server.go              # Server initialization & graceful shutdown
│   ├── config/
│   │   └── config.go              # Configuration management (Koanf)
│   ├── connections/
│   │   ├── database.go            # PostgreSQL connection pool
│   │   └── redis.go               # Redis client
│   ├── database/
│   │   ├── schema/                # Goose migrations
│   │   │   ├── 001_create_categories.sql
│   │   │   ├── 002_create_todos.sql
│   │   │   └── 003_create_comments.sql
│   │   ├── queries/               # SQL queries for SQLC
│   │   │   ├── categories.sql
│   │   │   ├── todos.sql
│   │   │   └── comments.sql
│   │   └── db/                    # Generated SQLC code (after running sqlc generate)
│   ├── errs/
│   │   └── errors.go              # Custom error types
│   ├── feature/
│   │   ├── category/              # Category feature module
│   │   │   ├── model.go           # Domain model
│   │   │   ├── dto.go             # Data Transfer Objects
│   │   │   ├── repository.go     # Data access layer
│   │   │   ├── service.go         # Business logic layer
│   │   │   ├── handler.go         # HTTP handlers
│   │   │   └── routes.go          # Route registration
│   │   ├── todo/                  # Todo feature module
│   │   │   ├── model.go
│   │   │   ├── dto.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   ├── handler.go
│   │   │   └── routes.go
│   │   └── comment/               # Comment feature module
│   │       ├── model.go
│   │       ├── dto.go
│   │       ├── repository.go
│   │       ├── service.go
│   │       ├── handler.go
│   │       └── routes.go
│   ├── logger/
│   │   └── logger.go              # Zerolog wrapper
│   ├── middlewares/
│   │   ├── error_handler.go       # Global error handler
│   │   ├── logger.go              # Request logger
│   │   └── setup.go               # Middleware initialization
│   └── validation/
│       └── validator.go           # Request validation utilities
├── scripts/
│   └── migrate.sh                 # Database migration helper
├── config.yaml                    # Configuration file
├── docker-compose.yml             # Docker services
├── Makefile                       # Build automation
├── sqlc.yaml                      # SQLC configuration
├── go.mod                         # Go modules
└── README.md                      # Documentation
```

## 🏗️ Architecture Layers

### 1. **Route Layer**
- Registers HTTP routes
- Maps URLs to handlers
- Example: `category/routes.go`

### 2. **Handler Layer**
- Processes HTTP requests/responses
- Binds request data
- Calls service layer
- Returns JSON responses
- Example: `category/handler.go`

### 3. **Service Layer**
- Contains business logic
- Validates requests
- Orchestrates repository calls
- Transforms models to DTOs
- Example: `category/service.go`

### 4. **Repository Layer**
- Data access layer
- Executes database queries
- Converts database models to domain models
- Example: `category/repository.go`

### 5. **DTO Layer**
- Request/Response objects
- Validation logic
- Data transformation
- Example: `category/dto.go`

### 6. **Model Layer**
- Domain entities
- Core business objects
- Example: `category/model.go`

## 🚀 Quick Start

### Prerequisites
```bash
# Install required tools
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

### Setup Steps

1. **Start infrastructure:**
```bash
make docker-up
```

2. **Run migrations:**
```bash
make migrate-up
```

3. **Generate SQLC code:**
```bash
make sqlc
```

4. **Generate Swagger docs:**
```bash
make swagger
```

5. **Run application:**
```bash
make run
```

Or use the all-in-one command:
```bash
make setup && make dev
```

## 📊 Dependency Injection Flow

```
main.go
  ├── Load Config
  ├── Initialize Logger
  ├── Create Server
  ├── Initialize Connections (DB, Redis)
  │
  ├── Create Repositories (injecting DB, Logger)
  │   ├── CategoryRepository(db, logger)
  │   ├── TodoRepository(db, logger)
  │   └── CommentRepository(db, logger)
  │
  ├── Create Services (injecting Repositories, Logger)
  │   ├── CategoryService(categoryRepo, logger)
  │   ├── TodoService(todoRepo, categoryRepo, logger)
  │   └── CommentService(commentRepo, todoRepo, logger)
  │
  ├── Create Handlers (injecting Services, Logger)
  │   ├── CategoryHandler(categoryService, logger)
  │   ├── TodoHandler(todoService, logger)
  │   └── CommentHandler(commentService, logger)
  │
  └── Setup HTTP Server (handlers → routes → middlewares)
```

## 🔧 Configuration

### File-based (config.yaml)
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
```

### Environment Variables (override config.yaml)
```bash
export APP_SERVER_PORT=9000
export APP_DATABASE_HOST=prod-db.example.com
```

## 🛠️ Key Features Implemented

### 1. Global Error Handler
- Custom error types
- Proper HTTP status codes
- Detailed error responses
- Error logging

### 2. Request Validation
- Field-level validation
- Custom validation rules
- Detailed validation errors

### 3. Database Management
- Connection pooling
- Health checks
- Graceful shutdown
- Type-safe queries (SQLC)

### 4. Logging
- Structured logging
- Request/response logging
- Error logging
- Performance metrics

### 5. Graceful Shutdown
- Signal handling (SIGTERM, SIGINT)
- Connection cleanup
- Pending request completion

## 📝 API Endpoints

### Categories
- `POST /api/v1/categories` - Create category
- `GET /api/v1/categories` - List categories
- `GET /api/v1/categories/:id` - Get category
- `PUT /api/v1/categories/:id` - Update category
- `DELETE /api/v1/categories/:id` - Delete category

### Todos
- `POST /api/v1/todos` - Create todo
- `GET /api/v1/todos` - List todos
- `GET /api/v1/todos/:id` - Get todo

### Comments
- `POST /api/v1/comments` - Create comment
- `GET /api/v1/comments/:id` - Get comment
- `GET /api/v1/todos/:todoId/comments` - List todo comments

### System
- `GET /health` - Health check
- `GET /swagger/*` - API documentation

## 🧪 Example API Requests

### Create Category
```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Work",
    "description": "Work-related tasks",
    "color": "#FF5733"
  }'
```

### Create Todo
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete project documentation",
    "description": "Write comprehensive API docs",
    "status": "pending",
    "priority": "high"
  }'
```

### List Todos
```bash
curl http://localhost:8080/api/v1/todos?limit=10&offset=0
```

## 🔍 SQLC Code Generation

SQLC generates type-safe Go code from SQL queries:

**Input (queries/categories.sql):**
```sql
-- name: GetCategoryByID :one
SELECT * FROM categories WHERE id = $1;
```

**Output (db/categories.sql.go):**
```go
func (q *Queries) GetCategoryByID(ctx context.Context, id uuid.UUID) (Category, error) {
    // Generated type-safe code
}
```

## 🔐 Best Practices Implemented

1. **Separation of Concerns** - Each layer has a single responsibility
2. **Dependency Injection** - Loose coupling between components
3. **Error Handling** - Comprehensive error types and handling
4. **Validation** - Input validation at service layer
5. **Logging** - Structured logging throughout
6. **Configuration** - Centralized, environment-aware config
7. **Database** - Connection pooling, health checks
8. **HTTP** - Proper status codes, middleware chain
9. **Security** - CORS, rate limiting, security headers
10. **Documentation** - Swagger/OpenAPI specs

## 📚 Make Commands

```bash
make help           # Show all available commands
make build          # Build the application
make run            # Run the application
make test           # Run tests
make clean          # Clean build artifacts
make migrate-up     # Run migrations
make migrate-down   # Rollback migration
make sqlc           # Generate SQLC code
make swagger        # Generate Swagger docs
make docker-up      # Start Docker services
make docker-down    # Stop Docker services
make setup          # Complete setup
make dev            # Start development environment
```

## 🎯 Next Steps for Enhancement

1. **Authentication & Authorization**
   - JWT tokens
   - Role-based access control
   - User management

2. **Testing**
   - Unit tests for all layers
   - Integration tests
   - Mock repositories

3. **Background Jobs**
   - Implement Asynq for async tasks
   - Email notifications with Resend
   - Scheduled tasks

4. **Caching**
   - Redis caching layer
   - Cache invalidation strategies

5. **Monitoring**
   - Prometheus metrics
   - Distributed tracing
   - Performance monitoring

6. **Deployment**
   - Kubernetes manifests
   - CI/CD pipelines
   - Production configuration

## 🤝 Contributing

This codebase follows professional Go standards and is structured for easy maintenance and scalability. Each feature is self-contained, making it easy to add new features or modify existing ones.

## 📄 License

MIT License

---

**Built with ❤️ using Go and modern backend practices**
