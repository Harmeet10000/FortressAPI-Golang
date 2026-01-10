# LLM Instructions for Fortress API

## Overview

You are an AI Pair Programming Assistant specializing in Go backend development. This project uses modern Go best practices with a focus on scalability, modular architecture, and separation of concerns. You are expected to provide expert-level guidance on Go-specific patterns, performance optimization, and architectural decisions.

## Tech Stack

- **Framework**: Echo (HTTP router)
- **Database**: PostgreSQL with pgx driver
- **Query Generation**: SQLc (type-safe SQL code generation)
- **Database Migrations**: Goose (migration management)
- **Validation**: go-playground/validator
- **Configuration**: Viper
- **Dependency Injection**: Uber FX
- **Logging**: Uber Zap
- **Message Queue**: RabbitMQ
- **Caching**: Redis
- **Environment**: godotenv
- **gRPC**: Protocol Buffers
- **Go Version**: 1.25.x

## Architecture Principles

### 1. **Separation of Concerns**

The project follows a layered architecture:

```
internal/
├── controllers/    # HTTP handlers, request/response mapping
├── services/       # Business logic layer
├── repository/     # Data access layer (GORM models)
├── middlewares/    # HTTP middleware (auth, logging, etc.)
├── routes/         # Route definitions
├── helpers/        # Utility functions
└── utils/          # Common utilities
```

Each layer has clear responsibilities:
- **Controllers**: Parse requests, validate input, call services, return responses
- **Services**: Implement business logic, orchestrate repository calls, handle errors
- **Repository**: Data persistence, GORM models, database queries
- **Middlewares**: Cross-cutting concerns (auth, logging, rate limiting)

### 2. **Modular Design**

When implementing features:
- Create self-contained modules with clear interfaces
- Use dependency injection (FX) to wire dependencies
- Avoid circular dependencies
- Keep packages small and focused (max 1-2 responsibilities)

### 3. **Error Handling**

- Always wrap errors with context using standard library `fmt.Errorf`
- Define custom error types for specific business logic failures
- Never ignore errors (avoid `_ = err`)
- Use Zap for structured logging of errors with context

Example:
```go
if err != nil {
    s.logger.Error("failed to fetch user", zap.Error(err), zap.Int("user_id", userID))
    return nil, fmt.Errorf("fetch user: %w", err)
}
```

### 4. **Validation**

Use `go-playground/validator` for input validation:
- Define validation tags on struct fields
- Validate in controllers before passing to services
- Return user-friendly error messages for validation failures

Example:
```go
type CreateUserRequest struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required,min=2,max=50"`
}
```

## Code Organization Guidelines

### Package Naming
- Use lowercase, single-word package names when possible
- Use plural names for collections (e.g., `users`, not `user`)
- Avoid generic names like `common`, `helper` (use `helpers`)

### File Organization
- Keep files under 500 lines when possible
- One interface per file if it's a key abstraction
- Group related functions in same file

### Dependency Injection (FX)

All dependencies should be injected through FX providers:

```go
func NewUserService(repo repository.User, logger *zap.Logger) *UserService {
    return &UserService{
        repo:   repo,
        logger: logger,
    }
}
```

Provide constructors to FX:
```go
fx.Provide(
    NewUserService,
    repository.NewUserRepository,
)
```

### Configuration (Viper)

- Load configuration from environment, `.env` files, and config files
- Define a config struct with validation
- Use Viper to unmarshal into the struct
- Never hardcode secrets or config values

Example:
```go
type Config struct {
    Database DatabaseConfig `mapstructure:"database"`
    Server   ServerConfig   `mapstructure:"server"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host" validate:"required"`
    Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`
    User     string `mapstructure:"user" validate:"required"`
    Password string `mapstructure:"password" validate:"required"`
}
```

### Logging (Zap)

Use structured logging with Zap:
- Use appropriate log levels: `Debug`, `Info`, `Warn`, `Error`, `Fatal`
- Always add relevant context fields (IDs, user info, etc.)
- Never log sensitive information (passwords, tokens)

Example:
```go
s.logger.Info("user created successfully", 
    zap.String("user_id", user.ID),
    zap.String("email", user.Email),
)

s.logger.Error("database connection failed",
    zap.Error(err),
    zap.String("host", config.DB.Host),
)
```

## HTTP Handler Pattern

When implementing HTTP handlers with Echo:

```go
// In controllers/user.go
type UserController struct {
    service services.User
    logger  *zap.Logger
    validate *validator.Validate
}

func NewUserController(service services.User, logger *zap.Logger, validate *validator.Validate) *UserController {
    return &UserController{
        service: service,
        logger: logger,
        validate: validate,
    }
}

func (c *UserController) Create(ctx echo.Context) error {
    // 1. Parse request body
    var req CreateUserRequest
    if err := ctx.BindJSON(&req); err != nil {
        c.logger.Error("invalid request", zap.Error(err))
        return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
    }

    // 2. Validate input
    if err := c.validate.Struct(req); err != nil {
        c.logger.Warn("validation failed", zap.Error(err))
        return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "validation error"})
    }

    // 3. Call service
    user, err := c.service.Create(ctx.Request().Context(), req)
    if err != nil {
        c.logger.Error("failed to create user", zap.Error(err))
        return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
    }

    // 4. Return response
    return ctx.JSON(http.StatusCreated, user)
}
```

## Database Patterns (SQLc + PostgreSQL)

### SQL Queries and Code Generation

Define SQL queries in `.sql` files and let SQLc generate type-safe Go code:

```sql
-- queries/users.sql

-- name: CreateUser :one
INSERT INTO users (id, email, name, created_at)
VALUES ($1, $2, $3, $4)
RETURNING id, email, name, created_at;

-- name: GetUserByID :one
SELECT id, email, name, created_at FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, email, name, created_at FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET email = $1, name = $2
WHERE id = $3
RETURNING id, email, name, created_at;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
```

### Goose Migrations

Create database migrations with Goose:

```sql
-- migrations/001_create_users_table.sql
-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

-- +goose Down
DROP TABLE IF EXISTS users;
```

Run migrations:
```bash
goose -dir migrations postgres "postgres://user:password@localhost/dbname" up
```

### Repository Pattern with SQLc

Define interfaces and implement with generated SQLc code:

```go
type UserRepository interface {
    CreateUser(ctx context.Context, id string, email string, name string) (*User, error)
    GetUserByID(ctx context.Context, id string) (*User, error)
    ListUsers(ctx context.Context, limit int32, offset int32) ([]User, error)
    UpdateUser(ctx context.Context, email string, name string, id string) (*User, error)
    DeleteUser(ctx context.Context, id string) error
}

type userRepository struct {
    db *pgx.Conn
    q  *db.Queries  // Generated by SQLc
}

func NewUserRepository(db *pgx.Conn) UserRepository {
    return &userRepository{
        db: db,
        q:  db.Queries(db),
    }
}

func (r *userRepository) CreateUser(ctx context.Context, id, email, name string) (*User, error) {
    row, err := r.q.CreateUser(ctx, db.CreateUserParams{
        ID:        id,
        Email:     email,
        Name:      name,
        CreatedAt: time.Now(),
    })
    if err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }
    return &User{
        ID:        row.ID,
        Email:     row.Email,
        Name:      row.Name,
        CreatedAt: row.CreatedAt,
    }, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
    row, err := r.q.GetUserByID(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("get user: %w", err)
    }
    return &User{
        ID:        row.ID,
        Email:     row.Email,
        Name:      row.Name,
        CreatedAt: row.CreatedAt,
    }, nil
}
```

### Why SQLc?

- **Type-Safe**: Generated code with compile-time guarantees
- **Performance**: No ORM overhead, direct SQL control
- **Simplicity**: Write standard SQL, get Go functions
- **Maintainability**: SQL and Go are kept in sync
- **No Runtime Magic**: Straightforward and debuggable

## Service Layer Pattern

Services implement business logic:

```go
type UserService struct {
    repo       repository.User
    logger     *zap.Logger
    cache      *redis.Client  // optional
    rabbitmq   *amqp.Channel  // optional
}

func (s *UserService) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
    // 1. Business logic validation
    // 2. Call repository methods
    // 3. Cache invalidation
    // 4. Publish events to message queue
    // 5. Return result
}
```

## Middleware Pattern (Echo)

Define middleware for cross-cutting concerns:

```go
func LoggingMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            logger.Info("incoming request",
                zap.String("method", c.Request().Method),
                zap.String("path", c.Request().URL.Path),
            )
            return next(c)
        }
    }
}

// Register middleware
echo := echo.New()
echo.Use(LoggingMiddleware(logger))
```

## Testing Guidelines

- Write tests for service layer (business logic)
- Use table-driven tests for multiple scenarios
- Mock repositories and external dependencies
- Aim for >80% coverage on critical paths

## Performance Considerations

1. **Database Optimization**:
   - Use appropriate indexes in GORM models
   - Implement pagination for list endpoints
   - Use lazy loading/preloading strategically
   - Connection pooling configured through database/sql settings

2. **Caching**:
   - Cache frequently accessed data in Redis
   - Implement cache invalidation strategy
   - Set appropriate TTLs

3. **Async Processing**:
   - Use RabbitMQ for long-running tasks
   - Implement circuit breakers for external calls
   - Use context timeouts appropriately

4. **Concurrency**:
   - Use goroutines appropriately
   - Protect shared state with sync primitives
   - Use `context.Context` for cancellation and deadlines

## Security Best Practices

1. **Authentication & Authorization**:
   - Validate all incoming requests
   - Use strong encryption for sensitive data (crypto package)
   - Never hardcode secrets

2. **Input Validation**:
   - Validate all user inputs in controllers
   - Use SQL parameterized queries (GORM handles this)
   - Sanitize outputs

3. **Error Messages**:
   - Don't expose internal error details to clients
   - Log full errors internally
   - Return generic error messages to clients

4. **Environment Secrets**:
   - Use `.env` files (with godotenv) for local development
   - Use CI/CD secrets for production
   - Never commit `.env` or secret files

## Response Structure

Use consistent response structures:

```go
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
    Page    int         `json:"page"`
    Total   int64       `json:"total"`
}
```

## Naming Conventions

- **Interfaces**: Use `Reader`, `Writer`, `UserService` (not `IUserService`)
- **Receiver variables**: Use short names (`s` for service, `r` for repository, `c` for controller)
- **Methods**: Use active verbs (`Create`, `Fetch`, `Update`, not `GetUser`, `UserGet`)
- **Variables**: Be descriptive and concise (`userID`, `createdAt`, not `u`, `ca`)

## Code Review Checklist

When reviewing code or implementing features, ensure:
- [ ] Proper error handling with context
- [ ] No circular dependencies
- [ ] Dependency injection through constructors
- [ ] Structured logging with Zap
- [ ] Input validation before processing
- [ ] Separation of concerns (controller/service/repo)
- [ ] No hardcoded values or secrets
- [ ] Appropriate use of context and timeouts
- [ ] Tests for critical paths
- [ ] Consistent response structures

## Common Patterns

### Creating New Feature

1. Create models in `internal/repository/models.go`
2. Create repository interface and implementation
3. Define request/response structs
4. Create service with business logic
5. Create controller handlers
6. Define routes
7. Wire dependencies in FX
8. Write tests

### Adding New External Service

1. Create client package (`internal/clients/` or `pkg/clients/`)
2. Define client interface
3. Implement client with proper timeout and error handling
4. Inject into service layer
5. Add retry logic if needed
6. Add proper logging

### Handling Long-Running Operations

1. Accept request in controller
2. Publish event to RabbitMQ
3. Return immediate response with task ID
4. Implement worker to process event
5. Store result in database or cache
6. Provide status endpoint for polling

## Go 1.25 Modern Patterns

### Range Over Integers (Go 1.22+)

Simplified iteration:

```go
// Old way
for i := 0; i < 10; i++ {
    // ...
}

// Go 1.22+ way
for i := range 10 {
    // i is the iteration count
}

// Combining with iterating over maps
for key := range myMap {
    // ...
}
```

### Structured Concurrency with Range Over Channels (Go 1.22+)

```go
func fetchUsers(ctx context.Context) ([]User, error) {
    usersChan := make(chan User, 10)
    errsChan := make(chan error, 1)

    go func() {
        defer close(usersChan)
        users, err := r.ListUsers(ctx, 100, 0)
        if err != nil {
            errsChan <- err
            return
        }
        for _, user := range users {
            usersChan <- user
        }
    }()

    var users []User
    for user := range usersChan {
        users = append(users, user)
    }

    select {
    case err := <-errsChan:
        return nil, err
    default:
        return users, nil
    }
}
```

### Iterator Protocol (Go 1.22+)

Define custom iterators for cleaner code:

```go
// Define an iterator that yields users
func (r *userRepository) IterUsers(ctx context.Context) iter.Seq2[User, error] {
    return func(yield func(User, error) bool) {
        users, err := r.ListUsers(ctx, 1000, 0)
        if err != nil {
            yield(User{}, err)
            return
        }
        for _, user := range users {
            if !yield(user, nil) {
                return
            }
        }
    }
}

// Usage
for user, err := range r.IterUsers(ctx) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(user.Email)
}
```

### Error Handling with errors.Join (Go 1.20+)

```go
func validateUser(user *User) error {
    var errs []error
    
    if user.Email == "" {
        errs = append(errs, errors.New("email is required"))
    }
    if user.Name == "" {
        errs = append(errs, errors.New("name is required"))
    }
    
    if len(errs) > 0 {
        return errors.Join(errs...)
    }
    return nil
}
```

### Slices Package (Go 1.21+)

Use stdlib slices for cleaner list operations:

```go
import "slices"

// Check if user exists in slice
if slices.Contains(userIDs, requestedID) {
    // ...
}

// Find user by condition
user, found := slices.BinarySearchFunc(users, targetID, func(u User, id string) int {
    return strings.Compare(u.ID, id)
})

// Sort users by email
slices.SortFunc(users, func(a, b User) int {
    return strings.Compare(a.Email, b.Email)
})

// Filter users
activeUsers := slices.DeleteFunc(users, func(u User) bool {
    return u.DeletedAt != nil
})
```

### Maps Package (Go 1.21+)

Functional operations on maps:

```go
import "maps"

// Copy a map
mapCopy := maps.Clone(originalMap)

// Check equality
if maps.Equal(map1, map2) {
    // ...
}
```

### Clear Built-in (Go 1.21+)

Clear maps and slices:

```go
// Clear a map
clear(myMap)

// Clear a slice
clear(mySlice)
```

### Modern Error Handling Pattern

```go
// Go 1.25 approach: Use error types effectively
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// Usage with errors.As (Go 1.13+)
if err := validateUser(user); err != nil {
    var valErr ValidationError
    if errors.As(err, &valErr) {
        log.Printf("Field validation failed: %v", valErr.Field)
    }
}
```

## Additional Resources

- [SQLc Documentation](https://sqlc.dev)
- [Goose Migration Tool](https://github.com/pressly/goose)
- [Go 1.25 Release Notes](https://golang.org/doc/go1.25)
- [Effective Go](https://golang.org/doc/effective_go)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Echo Framework Documentation](https://echo.labstack.com)
- [Uber FX Documentation](https://pkg.go.dev/go.uber.org/fx)
- [Zap Logger Documentation](https://pkg.go.dev/go.uber.org/zap)
- [pgx PostgreSQL Driver](https://github.com/jackc/pgx)

## When to Ask for Help

Ask for assistance when:
- Architectural decisions affect multiple layers
- Performance optimization is needed
- Complex error scenarios arise
- Adding new external dependencies
- Implementing sophisticated caching or async patterns
