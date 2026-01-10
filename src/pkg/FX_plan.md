# Complete Fx Learning Plan - Zero to Hero

This is a structured 4-week plan to master Fx from fundamentals to advanced patterns.

---

## 📚 **Week 1: Foundations**

### Day 1-2: Understanding Dependency Injection
**Goal**: Understand WHY DI exists before learning Fx

#### Theory
- What is Dependency Injection?
- Constructor injection vs Field injection vs Method injection
- Benefits: testability, modularity, loose coupling
- Problems DI solves in large codebases

#### Practice
```go
// Exercise 1: Manual DI without any framework
// Create a simple app with Logger -> Database -> Repository -> Service -> Handler
// Wire them manually in main()

type Logger interface {
    Info(msg string)
}

type Database struct {
    logger Logger
}

type UserRepository struct {
    db     *Database
    logger Logger
}

type UserService struct {
    repo   *UserRepository
    logger Logger
}

// Task: Wire all dependencies manually
func main() {
    logger := NewConsoleLogger()
    db := NewDatabase(logger)
    repo := NewUserRepository(db, logger)
    service := NewUserService(repo, logger)
    
    // Notice the pain? This is what Fx solves.
}
```

**Resources**:
- Read: [Martin Fowler's DI article](https://martinfowler.com/articles/injection.html)
- Watch: "Dependency Injection" basics on YouTube

---

### Day 3-4: Fx Basics - Provide & Invoke
**Goal**: Create your first Fx application

#### Concepts
- `fx.New()` - Application container
- `fx.Provide()` - Register constructors
- `fx.Invoke()` - Execute functions after wiring
- How Fx resolves dependencies automatically

#### Practice
```go
// Exercise 2: Convert Day 1 manual wiring to Fx
package main

import (
    "fmt"
    "go.uber.org/fx"
)

type Logger struct{}

func (l *Logger) Info(msg string) {
    fmt.Println("[INFO]", msg)
}

func NewLogger() *Logger {
    return &Logger{}
}

type Database struct {
    logger *Logger
}

func NewDatabase(logger *Logger) *Database {
    logger.Info("Database initialized")
    return &Database{logger: logger}
}

type UserService struct {
    db     *Database
    logger *Logger
}

func NewUserService(db *Database, logger *Logger) *UserService {
    logger.Info("UserService initialized")
    return &UserService{db: db, logger: logger}
}

func PrintServices(service *UserService, logger *Logger) {
    logger.Info("All services ready!")
}

func main() {
    fx.New(
        fx.Provide(
            NewLogger,
            NewDatabase,
            NewUserService,
        ),
        fx.Invoke(PrintServices),
    ).Run()
}
```

**Tasks**:
1. Run the code, observe initialization order
2. Add a new `CacheService` that depends on Logger
3. Make Database depend on CacheService
4. Try creating a circular dependency (Logger -> Database -> Logger) and see the error

**Resources**:
- Official Fx docs: https://uber-go.github.io/fx/
- Read the "Getting Started" section

---

### Day 5-6: Lifecycle Management
**Goal**: Master graceful startup and shutdown

#### Concepts
- `fx.Lifecycle` interface
- `fx.Hook` with OnStart and OnStop
- Shutdown order (reverse of startup)
- Context timeouts in lifecycle hooks

#### Practice
```go
// Exercise 3: Build a service with proper lifecycle
package main

import (
    "context"
    "database/sql"
    "fmt"
    "net/http"
    "time"

    _ "github.com/mattn/go-sqlite3"
    "go.uber.org/fx"
    "go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
    return zap.NewDevelopment()
}

func NewDatabase(lc fx.Lifecycle, logger *zap.Logger) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        return nil, err
    }

    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            logger.Info("📦 Starting database connection")
            if err := db.PingContext(ctx); err != nil {
                return err
            }
            logger.Info("✅ Database connected")
            return nil
        },
        OnStop: func(ctx context.Context) error {
            logger.Info("🛑 Closing database connection")
            return db.Close()
        },
    })

    return db, nil
}

func NewHTTPServer(lc fx.Lifecycle, logger *zap.Logger) *http.Server {
    mux := http.NewServeMux()
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            logger.Info("🚀 Starting HTTP server on :8080")
            go func() {
                if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                    logger.Fatal("Server error", zap.Error(err))
                }
            }()
            return nil
        },
        OnStop: func(ctx context.Context) error {
            logger.Info("🛑 Shutting down HTTP server")
            return server.Shutdown(ctx)
        },
    })

    return server
}

func main() {
    fx.New(
        fx.Provide(
            NewLogger,
            NewDatabase,
            NewHTTPServer,
        ),
    ).Run()
}
```

**Tasks**:
1. Run the app, hit Ctrl+C, observe shutdown order
2. Add a Redis client with lifecycle hooks
3. Make HTTP server wait 5 seconds before shutting down (simulate draining)
4. Add a background worker that stops gracefully

**Challenge**: 
Create a service that fails to start (return error from OnStart). See how Fx handles it.

---

### Day 7: Week 1 Project
**Goal**: Build a complete mini-application

**Project: TODO API with SQLite**

Requirements:
- HTTP server with routes: POST /todos, GET /todos, GET /todos/:id
- SQLite database with lifecycle management
- Logger (zap)
- Proper graceful shutdown
- All wired with Fx

**Structure**:
```
todo-app/
├── main.go           # Fx application setup
├── handler.go        # HTTP handlers
├── service.go        # Business logic
├── repository.go     # Database operations
└── models.go         # Data structures
```

**Success Criteria**:
- App starts cleanly
- All dependencies auto-wired
- Graceful shutdown works (Ctrl+C)
- No manual dependency wiring in main()

---

## 📚 **Week 2: Intermediate Patterns**

### Day 8-9: Error Handling & Multiple Constructors
**Goal**: Handle initialization errors properly

#### Concepts
- Returning errors from constructors
- `fx.Error` type
- Providing multiple instances of same type
- Using `fx.Annotate` for disambiguation

#### Practice
```go
// Exercise 4: Error handling
package main

import (
    "context"
    "errors"
    "fmt"
    "os"

    "go.uber.org/fx"
    "go.uber.org/zap"
)

type Config struct {
    DatabaseURL string
    APIKey      string
}

// Constructor that can fail
func NewConfig() (*Config, error) {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return nil, errors.New("DATABASE_URL environment variable required")
    }

    apiKey := os.Getenv("API_KEY")
    if apiKey == "" {
        return nil, errors.New("API_KEY environment variable required")
    }

    return &Config{
        DatabaseURL: dbURL,
        APIKey:      apiKey,
    }, nil
}

// Multiple loggers example
func NewAppLogger() (*zap.Logger, error) {
    logger, err := zap.NewProduction()
    if err != nil {
        return nil, err
    }
    return logger.Named("app"), nil
}

func NewAccessLogger() (*zap.Logger, error) {
    logger, err := zap.NewProduction()
    if err != nil {
        return nil, err
    }
    return logger.Named("access"), nil
}

// Using fx.Annotate to provide multiple loggers
func main() {
    fx.New(
        fx.Provide(
            NewConfig,
            fx.Annotate(
                NewAppLogger,
                fx.ResultTags(`name:"appLogger"`),
            ),
            fx.Annotate(
                NewAccessLogger,
                fx.ResultTags(`name:"accessLogger"`),
            ),
        ),
        fx.Invoke(func(
            cfg *Config,
            appLogger *zap.Logger `name:"appLogger"`,
            accessLogger *zap.Logger `name:"accessLogger"`,
        ) {
            appLogger.Info("App started", zap.String("db", cfg.DatabaseURL))
            accessLogger.Info("Access logging enabled")
        }),
    ).Run()
}
```

**Tasks**:
1. Run without env vars, see Fx handle the error gracefully
2. Create two database connections (read replica + write primary)
3. Use `fx.Annotate` to distinguish them
4. Inject both into a service

---

### Day 10-11: Modules & Organization
**Goal**: Structure large applications

#### Concepts
- `fx.Module()` for grouping related dependencies
- Module composition
- Private vs exported providers
- Decorate pattern

#### Practice
```go
// Exercise 5: Modular architecture
package main

import (
    "go.uber.org/fx"
    "go.uber.org/zap"
)

// Database module
var DatabaseModule = fx.Module("database",
    fx.Provide(
        NewDatabaseConfig,
        NewDatabase,
        NewUserRepository,
        NewProductRepository,
    ),
)

// HTTP module
var HTTPModule = fx.Module("http",
    fx.Provide(
        NewRouter,
        NewUserHandler,
        NewProductHandler,
        NewHTTPServer,
    ),
)

// Cache module
var CacheModule = fx.Module("cache",
    fx.Provide(
        NewRedisClient,
        NewCacheService,
    ),
)

// Background jobs module
var JobsModule = fx.Module("jobs",
    fx.Provide(
        NewJobScheduler,
        NewEmailJob,
        NewCleanupJob,
    ),
)

func main() {
    fx.New(
        // Global dependencies
        fx.Provide(NewLogger),
        fx.Provide(NewConfig),
        
        // Feature modules
        DatabaseModule,
        CacheModule,
        HTTPModule,
        JobsModule,
        
        // Optional: disable modules based on config
        fx.Options(
            fx.Decorate(func(cfg *Config, logger *zap.Logger) *Config {
                if !cfg.EnableJobs {
                    logger.Info("Background jobs disabled")
                }
                return cfg
            }),
        ),
    ).Run()
}
```

**Tasks**:
1. Create an `AuthModule` with JWT service, auth middleware, and auth handler
2. Create a `MetricsModule` with Prometheus metrics
3. Make modules conditionally load based on config
4. Create module dependencies (HTTPModule depends on DatabaseModule)

---

### Day 12-13: Testing with Fx
**Goal**: Write testable code using Fx

#### Concepts
- Testing individual constructors
- Using `fx.New` in tests
- `fx.NopLogger` for test silence
- Mocking dependencies
- `fx.Populate` for dependency extraction

#### Practice
```go
// Exercise 6: Testing Fx applications
package main

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "go.uber.org/fx"
    "go.uber.org/fx/fxtest"
)

// Mock repository
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetUser(id int) (*User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

// Test using fxtest
func TestUserService(t *testing.T) {
    var service *UserService
    
    app := fxtest.New(t,
        fx.Provide(
            func() UserRepository {
                mock := &MockUserRepository{}
                mock.On("GetUser", 1).Return(&User{ID: 1, Name: "Alice"}, nil)
                return mock
            },
            NewUserService,
        ),
        fx.Populate(&service),
    )
    
    app.RequireStart()
    defer app.RequireStop()
    
    user, err := service.GetUser(1)
    assert.NoError(t, err)
    assert.Equal(t, "Alice", user.Name)
}

// Integration test with real dependencies
func TestIntegration(t *testing.T) {
    app := fxtest.New(t,
        fx.Provide(
            NewInMemoryDatabase,
            NewUserRepository,
            NewUserService,
        ),
        fx.Invoke(func(service *UserService) {
            // Run integration tests
            user, err := service.CreateUser("Bob")
            assert.NoError(t, err)
            assert.NotNil(t, user)
        }),
    )
    
    app.RequireStart()
    app.RequireStop()
}
```

**Tasks**:
1. Write unit tests for all services from Week 1 project
2. Write integration tests with real SQLite
3. Create test helpers for common Fx setups
4. Benchmark Fx initialization time

---

### Day 14: Week 2 Project
**Goal**: Build a modular e-commerce API

**Project: E-commerce API Backend**

Requirements:
- Multiple modules: Auth, Products, Orders, Payments
- PostgreSQL database (use Docker)
- Redis cache
- JWT authentication
- Full test coverage
- Structured with Fx modules

**Structure**:
```
ecommerce-api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── auth/
│   │   ├── module.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   └── repository.go
│   ├── products/
│   │   └── ...
│   ├── orders/
│   │   └── ...
│   └── payments/
│       └── ...
├── pkg/
│   ├── database/
│   ├── cache/
│   └── middleware/
└── go.mod
```

**Success Criteria**:
- Clean module separation
- All modules independently testable
- Can disable modules via config
- Graceful shutdown of all resources

---

## 📚 **Week 3: Advanced Patterns**

### Day 15-16: Value Groups & Result Objects
**Goal**: Handle collections of dependencies

#### Concepts
- `fx.Supply()` for providing values
- Value groups with `group:` tags
- Result objects for multiple returns
- Param objects for multiple inputs

#### Practice
```go
// Exercise 7: Value groups for plugins
package main

import (
    "fmt"
    "net/http"

    "go.uber.org/fx"
)

// Middleware interface
type Middleware interface {
    Handle(http.Handler) http.Handler
}

// Multiple middleware implementations
type LoggingMiddleware struct{}

func (m *LoggingMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Request:", r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

type AuthMiddleware struct{}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Auth logic...
        next.ServeHTTP(w, r)
    })
}

type CORSMiddleware struct{}

func (m *CORSMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        next.ServeHTTP(w, r)
    })
}

// Provide all middleware as a group
func NewLoggingMiddleware() Middleware {
    return &LoggingMiddleware{}
}

func NewAuthMiddleware() Middleware {
    return &AuthMiddleware{}
}

func NewCORSMiddleware() Middleware {
    return &CORSMiddleware{}
}

// Server that consumes all middleware
type ServerParams struct {
    fx.In

    Middlewares []Middleware `group:"middleware"`
}

func NewServer(params ServerParams) *http.Server {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello!"))
    })

    // Apply all middleware
    for i := len(params.Middlewares) - 1; i >= 0; i-- {
        handler = params.Middlewares[i].Handle(handler)
    }

    return &http.Server{
        Addr:    ":8080",
        Handler: handler,
    }
}

func main() {
    fx.New(
        fx.Provide(
            fx.Annotate(
                NewLoggingMiddleware,
                fx.ResultTags(`group:"middleware"`),
            ),
            fx.Annotate(
                NewAuthMiddleware,
                fx.ResultTags(`group:"middleware"`),
            ),
            fx.Annotate(
                NewCORSMiddleware,
                fx.ResultTags(`group:"middleware"`),
            ),
            NewServer,
        ),
    ).Run()
}
```

**Tasks**:
1. Create a plugin system where handlers are registered via value groups
2. Build a metrics collector that aggregates multiple metric sources
3. Create a notification system with multiple providers (email, SMS, push)
4. Use result objects to return multiple related dependencies

---

### Day 17-18: Dynamic Configuration & Feature Flags
**Goal**: Build configurable applications

#### Concepts
- Config-driven module loading
- Feature flags integration
- Environment-specific providers
- `fx.Decorate` for config transformation

#### Practice
```go
// Exercise 8: Feature flags and dynamic loading
package main

import (
    "go.uber.org/fx"
    "go.uber.org/zap"
)

type FeatureFlags struct {
    EnableCache      bool
    EnableAuth       bool
    EnableMetrics    bool
    EnableBackgroundJobs bool
}

func NewFeatureFlags() *FeatureFlags {
    return &FeatureFlags{
        EnableCache:      getEnvBool("ENABLE_CACHE", true),
        EnableAuth:       getEnvBool("ENABLE_AUTH", true),
        EnableMetrics:    getEnvBool("ENABLE_METRICS", false),
        EnableBackgroundJobs: getEnvBool("ENABLE_JOBS", true),
    }
}

func getEnvBool(key string, defaultVal bool) bool {
    // Implementation...
    return defaultVal
}

func ProvideOptionalCache(flags *FeatureFlags, logger *zap.Logger) fx.Option {
    if !flags.EnableCache {
        logger.Info("Cache disabled")
        return fx.Options()
    }

    return fx.Module("cache",
        fx.Provide(NewRedisClient),
        fx.Provide(NewCacheService),
    )
}

func ProvideOptionalAuth(flags *FeatureFlags, logger *zap.Logger) fx.Option {
    if !flags.EnableAuth {
        logger.Info("Auth disabled")
        return fx.Options()
    }

    return fx.Module("auth",
        fx.Provide(NewJWTService),
        fx.Provide(NewAuthMiddleware),
    )
}

func main() {
    app := fx.New(
        fx.Provide(
            NewLogger,
            NewFeatureFlags,
        ),
        
        // Dynamically load modules based on flags
        fx.Options(
            fx.Invoke(ProvideOptionalCache),
            fx.Invoke(ProvideOptionalAuth),
        ),
    )
    
    app.Run()
}
```

**Tasks**:
1. Implement multi-environment configs (dev, staging, prod)
2. Create A/B testing framework with feature flags
3. Build hot-reload for configuration changes
4. Implement circuit breaker pattern with config thresholds

---

### Day 19-20: Performance & Observability
**Goal**: Production-ready applications

#### Concepts
- Fx initialization metrics
- Tracing dependency resolution
- Memory profiling
- Custom Fx logger
- Startup time optimization

#### Practice
```go
// Exercise 9: Observability and performance
package main

import (
    "context"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "go.uber.org/fx"
    "go.uber.org/zap"
)

var (
    fxStartupDuration = promauto.NewHistogram(prometheus.HistogramOpts{
        Name: "fx_startup_duration_seconds",
        Help: "Fx application startup duration",
    })

    componentInitDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name: "component_init_duration_seconds",
        Help: "Individual component initialization duration",
    }, []string{"component"})
)

// Instrumented constructor
func NewInstrumentedDatabase(lc fx.Lifecycle, logger *zap.Logger) (*Database, error) {
    start := time.Now()
    defer func() {
        duration := time.Since(start).Seconds()
        componentInitDuration.WithLabelValues("database").Observe(duration)
        logger.Info("Database init took", zap.Duration("duration", time.Since(start)))
    }()

    db := &Database{}
    
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            startTime := time.Now()
            // Connect...
            logger.Info("Database connection took", zap.Duration("duration", time.Since(startTime)))
            return nil
        },
        OnStop: func(ctx context.Context) error {
            return db.Close()
        },
    })

    return db, nil
}

// Custom Fx logger for structured logging
type FxLogger struct {
    logger *zap.Logger
}

func (l *FxLogger) Printf(format string, args ...interface{}) {
    l.logger.Sugar().Infof(format, args...)
}

func NewFxLogger(logger *zap.Logger) fxevent.Logger {
    return &fxevent.ZapLogger{Logger: logger}
}

func main() {
    startTime := time.Now()
    
    app := fx.New(
        fx.WithLogger(NewFxLogger),
        fx.Provide(
            NewLogger,
            NewInstrumentedDatabase,
            NewInstrumentedCache,
            NewInstrumentedHTTPServer,
        ),
        fx.Invoke(func(lc fx.Lifecycle) {
            lc.Append(fx.Hook{
                OnStart: func(ctx context.Context) error {
                    fxStartupDuration.Observe(time.Since(startTime).Seconds())
                    return nil
                },
            })
        }),
    )
    
    app.Run()
}
```

**Tasks**:
1. Add distributed tracing (OpenTelemetry)
2. Implement health checks for all components
3. Add startup time budget (fail if > 10s)
4. Create dashboard showing initialization metrics

---

### Day 21: Week 3 Project
**Goal**: Production-grade microservice

**Project: User Service with Full Observability**

Requirements:
- Complete CRUD API
- PostgreSQL + Redis
- JWT auth
- Prometheus metrics
- Distributed tracing
- Health checks
- Feature flags
- Structured logging
- Graceful degradation (works without Redis)
- Full test coverage

**Success Criteria**:
- Startup time < 5 seconds
- All errors traced
- Can toggle features without restart
- 90%+ code coverage
- Load tested (1000 req/s)

---

## 📚 **Week 4: Real-World Patterns & Mastery**

### Day 22-23: Microservices Patterns
**Goal**: Build distributed systems with Fx

#### Concepts
- gRPC server integration
- Service discovery
- Message queue consumers
- Multi-protocol servers (HTTP + gRPC)

#### Practice
```go
// Exercise 10: Multi-protocol microservice
package main

import (
    "context"
    "net"
    "net/http"

    "go.uber.org/fx"
    "go.uber.org/zap"
    "google.golang.org/grpc"
)

// gRPC server
type GRPCServer struct {
    server *grpc.Server
}

func NewGRPCServer(lc fx.Lifecycle, logger *zap.Logger) *GRPCServer {
    server := grpc.NewServer()
    
    // Register services...
    
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            listener, err := net.Listen("tcp", ":9090")
            if err != nil {
                return err
            }
            
            logger.Info("gRPC server starting on :9090")
            go server.Serve(listener)
            return nil
        },
        OnStop: func(ctx context.Context) error {
            logger.Info("Stopping gRPC server")
            server.GracefulStop()
            return nil
        },
    })
    
    return &GRPCServer{server: server}
}

// HTTP server for REST API
func NewHTTPServer(lc fx.Lifecycle, logger *zap.Logger) *http.Server {
    server := &http.Server{
        Addr: ":8080",
    }
    
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            logger.Info("HTTP server starting on :8080")
            go server.ListenAndServe()
            return nil
        },
        OnStop: func(ctx context.Context) error {
            logger.Info("Stopping HTTP server")
            return server.Shutdown(ctx)
        },
    })
    
    return server
}

// Kafka consumer
type KafkaConsumer struct{}

func NewKafkaConsumer(lc fx.Lifecycle, logger *zap.Logger) *KafkaConsumer {
    consumer := &KafkaConsumer{}
    
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            logger.Info("Starting Kafka consumer")
            go consumer.Start()
            return nil
        },
        OnStop: func(ctx context.Context) error {
            logger.Info("Stopping Kafka consumer")
            return consumer.Stop()
        },
    })
    
    return consumer
}

func main() {
    fx.New(
        fx.Provide(
            NewLogger,
            NewDatabase,
            NewGRPCServer,
            NewHTTPServer,
            NewKafkaConsumer,
        ),
    ).Run()
}
```

**Tasks**:
1. Build a microservice with HTTP + gRPC + WebSocket
2. Add RabbitMQ consumer for async processing
3. Implement circuit breaker between services
4. Add service mesh integration (Istio/Linkerd)

---

### Day 24-25: Advanced Testing Strategies
**Goal**: Master testing complex Fx apps

#### Topics
- Testing lifecycle hooks
- Testing async operations
- Contract testing
- Chaos engineering with Fx
- Performance testing

#### Practice
```go
// Exercise 11: Advanced testing patterns
package main

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "go.uber.org/fx"
    "go.uber.org/fx/fxtest"
)

// Test lifecycle timing
func TestLifecycleOrdering(t *testing.T) {
    var order []string
    
    newService := func(name string) func(fx.Lifecycle) {
        return func(lc fx.Lifecycle) {
            lc.Append(fx.Hook{
                OnStart: func(ctx context.Context) error {
                    order = append(order, name+":start")
                    return nil
                },
                OnStop: func(ctx context.Context) error {
                    order = append(order, name+":stop")
                    return nil
                },
            })
        }
    }
    
    app := fxtest.New(t,
        fx.Invoke(newService("A")),
        fx.Invoke(newService("B")),
        fx.Invoke(newService("C")),
    )
    
    app.RequireStart()
    assert.Equal(t, []string{"A:start", "B:start", "C:start"}, order)
    
    order = nil
    app.RequireStop()
    assert.Equal(t, []string{"C:stop", "B:stop", "A:stop"}, order)
}

// Test timeout handling
func TestLifecycleTimeout(t *testing.T) {
    app := fxtest.New(t,
        fx.StartTimeout(100*time.Millisecond),
        fx.Invoke(func(lc fx.Lifecycle) {
            lc.Append(fx.Hook{
                OnStart: func(ctx context.Context) error {
                    time.Sleep(200 * time.Millisecond) // Exceeds timeout
                    return nil
                },
            })
        }),
    )
    
    err := app.Start(context.Background())
    assert.Error(t, err)
}

// Benchmark Fx initialization
func BenchmarkFxInitialization(b *testing.B) {
    for i := 0; i < b.N; i++ {
        app := fx.New(
            fx.NopLogger,
            fx.Provide(
                NewLogger,
                NewDatabase,
                NewCache,
                NewService,
            ),
        )
        app.Start(context.Background())
        app.Stop(context.Background())
    }
}
```

**Tasks**:
1. Write chaos tests (randomly fail dependencies)
2. Test concurrent startup/shutdown
3. Benchmark with 100+ dependencies
4. Create test doubles for all external services

---

### Day 26: Production Deployment
**Goal**: Deploy Fx apps to production

#### Topics
- Docker containerization
- Kubernetes deployment
- Configuration management (Consul, etcd)
- Secrets management
- Blue-green deployments

#### Practice
```dockerfile
# Dockerfile for Fx app
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /server .

# Fx handles SIGTERM gracefully
CMD ["./server"]
```

```yaml
# kubernetes/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fx-service
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: fx-service
        image: fx-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
```yaml
            secretKeyRef:
              name: db-credentials
              key: url
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            Port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        lifecycle:
          preStop:
            exec:
              # Give Fx time to gracefully shutdown
              command: ["/bin/sh", "-c", "sleep 15"]
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

**Tasks**:
1. Create CI/CD pipeline (GitHub Actions)
2. Set up health checks endpoint
3. Implement readiness probe
4. Configure graceful shutdown timeout
5. Add Helm charts for deployment

---

### Day 27: Common Pitfalls & Debugging
**Goal**: Learn to troubleshoot Fx issues

#### Common Issues & Solutions

```go
// Exercise 12: Debugging Fx applications

// Problem 1: Circular Dependency
// BAD:
func NewServiceA(b *ServiceB) *ServiceA { return &ServiceA{} }
func NewServiceB(a *ServiceA) *ServiceB { return &ServiceB{} }
// Error: cycle detected in dependency graph

// SOLUTION: Extract interface
type EventBus interface {
    Publish(event Event)
}

func NewServiceA(bus EventBus) *ServiceA { return &ServiceA{} }
func NewServiceB(bus EventBus) *ServiceB { return &ServiceB{} }
func NewEventBus() EventBus { return &EventBusImpl{} }

// Problem 2: Missing Dependency
func main() {
    fx.New(
        fx.Provide(NewUserService), // Needs UserRepository!
    ).Run()
}
// Error: missing type: *UserRepository

// SOLUTION: Provide all dependencies
func main() {
    fx.New(
        fx.Provide(
            NewDatabase,
            NewUserRepository, // Add this
            NewUserService,
        ),
    ).Run()
}

// Problem 3: Goroutine leaks in lifecycle
// BAD:
func NewWorker(lc fx.Lifecycle) *Worker {
    w := &Worker{}
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            go w.Run() // No way to stop!
            return nil
        },
    })
    return w
}

// GOOD:
func NewWorker(lc fx.Lifecycle) *Worker {
    w := &Worker{
        stopChan: make(chan struct{}),
    }
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            go w.Run()
            return nil
        },
        OnStop: func(ctx context.Context) error {
            close(w.stopChan)
            return w.Wait() // Wait for goroutine to finish
        },
    })
    return w
}

// Problem 4: Blocking OnStart
// BAD:
OnStart: func(ctx context.Context) error {
    return server.ListenAndServe() // Blocks forever!
}

// GOOD:
OnStart: func(ctx context.Context) error {
    go server.ListenAndServe() // Non-blocking
    return nil
}

// Problem 5: Not using context timeout
// BAD:
OnStop: func(ctx context.Context) error {
    db.Close() // Ignores context!
    return nil
}

// GOOD:
OnStop: func(ctx context.Context) error {
    // Respect context deadline
    done := make(chan error, 1)
    go func() {
        done <- db.Close()
    }()
    
    select {
    case err := <-done:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}

// Problem 6: Order-dependent initialization
// BAD: Assuming database is ready before cache
func NewCache(db *sql.DB) *Cache {
    // db might not be connected yet!
    db.Ping() // Could fail
    return &Cache{}
}

// GOOD: Use lifecycle hooks
func NewCache(db *sql.DB, lc fx.Lifecycle) *Cache {
    cache := &Cache{}
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            // Now db is guaranteed to be started
            return cache.Initialize(db)
        },
    })
    return cache
}
```

#### Debugging Tools

```go
// Enable Fx debug logging
func main() {
    fx.New(
        fx.Provide(/* ... */),
        fx.WithLogger(func() fxevent.Logger {
            return &fxevent.ConsoleLogger{W: os.Stdout}
        }),
    ).Run()
}

// Visualize dependency graph
func main() {
    app := fx.New(
        fx.Provide(/* ... */),
        fx.Invoke(fx.Visualize(os.Stdout)), // Prints DOT graph
    )
    app.Run()
}

// Use fx.DryRun to test initialization
func TestDependencyGraph(t *testing.T) {
    err := fx.ValidateApp(
        fx.Provide(
            NewDatabase,
            NewUserService,
        ),
    )
    assert.NoError(t, err)
}
```

**Tasks**:
1. Create a "broken" app with all common mistakes
2. Fix each issue one by one
3. Document error messages and solutions
4. Create pre-commit hooks to catch issues

---

### Day 28: Final Project - Complete Microservice Platform

**Goal**: Build a production-ready microservice ecosystem

## **Final Capstone Project: Social Media Backend**

### Architecture Overview
```
┌─────────────────────────────────────────────────────┐
│                   API Gateway (Fx)                  │
│         HTTP + gRPC + WebSocket + GraphQL           │
└─────────────────────────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
┌───────▼──────┐  ┌───────▼──────┐  ┌──────▼───────┐
│ Auth Service │  │ User Service │  │ Post Service │
│    (Fx)      │  │    (Fx)      │  │    (Fx)      │
└──────────────┘  └──────────────┘  └──────────────┘
        │                 │                 │
        └─────────────────┼─────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
┌───────▼──────┐  ┌───────▼──────┐  ┌──────▼───────┐
│  PostgreSQL  │  │     Redis    │  │    Kafka     │
└──────────────┘  └──────────────┘  └──────────────┘
```

### Services to Build

#### 1. **Auth Service**
```go
// Features:
// - JWT token generation/validation
// - OAuth2 integration (Google, GitHub)
// - Rate limiting
// - Session management with Redis
// - Password hashing with bcrypt

// Module structure:
var AuthModule = fx.Module("auth",
    fx.Provide(
        NewJWTService,
        NewOAuth2Service,
        NewSessionStore,
        NewPasswordHasher,
        NewAuthRepository,
        NewAuthService,
        NewAuthHandler,
    ),
)
```

#### 2. **User Service**
```go
// Features:
// - User CRUD operations
// - Profile management
// - Follow/unfollow system
// - User search with Elasticsearch
// - Avatar upload to S3

var UserModule = fx.Module("users",
    fx.Provide(
        NewUserRepository,
        NewUserService,
        NewUserHandler,
        NewFollowService,
        NewSearchService,
        NewStorageService,
    ),
)
```

#### 3. **Post Service**
```go
// Features:
// - Create/read/update/delete posts
// - Like/unlike functionality
// - Comment system
// - Feed generation
// - Media processing

var PostModule = fx.Module("posts",
    fx.Provide(
        NewPostRepository,
        NewPostService,
        NewPostHandler,
        NewLikeService,
        NewCommentService,
        NewFeedService,
        NewMediaProcessor,
    ),
)
```

#### 4. **Notification Service**
```go
// Features:
// - Real-time WebSocket notifications
// - Email notifications (async via Kafka)
// - Push notifications
// - Notification preferences

var NotificationModule = fx.Module("notifications",
    fx.Provide(
        NewWebSocketServer,
        NewEmailService,
        NewPushService,
        NewNotificationRepository,
        NewNotificationService,
        NewKafkaConsumer,
    ),
)
```

### Core Infrastructure

```go
// cmd/api-gateway/main.go
package main

import (
    "go.uber.org/fx"
    "yourapp/internal/auth"
    "yourapp/internal/users"
    "yourapp/internal/posts"
    "yourapp/internal/notifications"
    "yourapp/pkg/config"
    "yourapp/pkg/database"
    "yourapp/pkg/cache"
    "yourapp/pkg/messaging"
    "yourapp/pkg/observability"
)

func main() {
    fx.New(
        // Core infrastructure
        config.Module,
        database.Module,
        cache.Module,
        messaging.Module,
        observability.Module,
        
        // Business modules
        auth.Module,
        users.Module,
        posts.Module,
        notifications.Module,
        
        // Server
        fx.Provide(NewAPIGateway),
        
        // Graceful shutdown timeout
        fx.StopTimeout(30*time.Second),
    ).Run()
}
```

### Required Features

**1. Observability Stack**
```go
// pkg/observability/module.go
var Module = fx.Module("observability",
    fx.Provide(
        NewPrometheusMetrics,
        NewJaegerTracer,
        NewStructuredLogger,
        NewHealthChecker,
    ),
)

// Metrics to track:
// - Request latency (per endpoint)
// - Database query time
// - Cache hit/miss ratio
// - Active WebSocket connections
// - Kafka consumer lag
// - Error rates by type
```

**2. Configuration Management**
```go
// pkg/config/config.go
type Config struct {
    Server    ServerConfig
    Database  DatabaseConfig
    Redis     RedisConfig
    Kafka     KafkaConfig
    Auth      AuthConfig
    Storage   StorageConfig
    Features  FeatureFlags
}

// Support multiple sources:
// - Environment variables
// - .env files
// - Consul/etcd
// - AWS Parameter Store
```

**3. Testing Requirements**
```go
// tests/integration_test.go
func TestUserRegistrationFlow(t *testing.T) {
    app := fxtest.New(t,
        // Use test doubles
        fx.Replace(NewMockEmailService()),
        fx.Replace(NewTestDatabase()),
        
        // Real business logic
        auth.Module,
        users.Module,
    )
    
    app.RequireStart()
    defer app.RequireStop()
    
    // Test complete user registration
    // 1. Register user
    // 2. Receive verification email
    // 3. Verify email
    // 4. Login
    // 5. Create profile
}

// Test coverage requirements:
// - Unit tests: 80%+
// - Integration tests: All happy paths
// - E2E tests: Critical user journeys
// - Load tests: 1000 req/s sustained
```

**4. Database Migrations**
```go
// pkg/database/migrations.go
func NewMigrator(db *sql.DB, lc fx.Lifecycle) *Migrator {
    migrator := &Migrator{db: db}
    
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            return migrator.Up(ctx)
        },
    })
    
    return migrator
}
```

**5. API Documentation**
```go
// Use Swagger/OpenAPI
// Auto-generate from code annotations
// Serve at /docs endpoint
```

### Deployment Configuration

```yaml
# docker-compose.yml
version: '3.8'
services:
  api-gateway:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/socialdb
      - REDIS_URL=redis://cache:6379
      - KAFKA_BROKERS=kafka:9092
    depends_on:
      - db
      - cache
      - kafka
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 3

  db:
    image: postgres:15
    environment:
      POSTGRES_DB: socialdb
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
    volumes:
      - postgres_data:/var/lib/postgresql/data

  cache:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data

  kafka:
    image: confluentinc/cp-kafka:latest
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092

  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

volumes:
  postgres_data:
  redis_data:
```

### Success Criteria

**Functionality** (40 points)
- ✅ All services communicate properly
- ✅ Authentication works end-to-end
- ✅ Real-time notifications work
- ✅ Feed generation is efficient
- ✅ File uploads work correctly

**Code Quality** (25 points)
- ✅ Clean Fx module separation
- ✅ No circular dependencies
- ✅ Proper error handling
- ✅ Consistent code style
- ✅ Good documentation

**Testing** (20 points)
- ✅ 80%+ unit test coverage
- ✅ Integration tests pass
- ✅ Load test passes (1000 req/s)
- ✅ All race conditions fixed

**Operations** (15 points)
- ✅ Graceful shutdown works
- ✅ Health checks respond correctly
- ✅ Metrics are collected
- ✅ Logs are structured
- ✅ Deploys via Docker Compose

---

## 📊 **Learning Milestones Checklist**

### Week 1: Foundations ✅
- [ ] Understand DI principles
- [ ] Create first Fx app
- [ ] Use Provide and Invoke
- [ ] Implement lifecycle hooks
- [ ] Complete TODO API project

### Week 2: Intermediate ✅
- [ ] Handle constructor errors
- [ ] Use fx.Annotate for disambiguation
- [ ] Create modular architecture
- [ ] Write testable code with Fx
- [ ] Complete e-commerce API project

### Week 3: Advanced ✅
- [ ] Use value groups
- [ ] Implement feature flags
- [ ] Add observability stack
- [ ] Performance optimization
- [ ] Complete user service project

### Week 4: Mastery ✅
- [ ] Build microservices
- [ ] Advanced testing strategies
- [ ] Production deployment
- [ ] Debug complex issues
- [ ] Complete final capstone project

---

## 📚 **Additional Resources**

### Official Documentation
- Fx Docs: https://uber-go.github.io/fx/
- Dig Docs: https://uber-go.github.io/dig/
- Go Blog: https://go.dev/blog/

### Video Tutorials
- "Building Robust Applications with Fx" (Uber Engineering)
- "Dependency Injection in Go" (JustForFunc)
- "Production Go Services" (Ardan Labs)

### Example Repositories
- Fx examples: https://github.com/uber-go/fx/tree/master/_examples
- Real-world Fx apps: Search GitHub for "fx.New" in Go repos

### Books
- "Let's Go" by Alex Edwards
- "Cloud Native Go" by Matthew Titmus
- "Microservices in Go" by Matt Heath

### Community
- Gophers Slack: #fx channel
- Reddit: r/golang
- Stack Overflow: [go] [fx] tags

---

## 🎯 **Post-Learning Action Items**

After completing this 4-week plan:

1. **Refactor an existing project** to use Fx
2. **Contribute to Fx** - Report bugs, submit PRs
3. **Write a blog post** about your Fx journey
4. **Mentor others** - Answer questions on Stack Overflow
5. **Build a side project** using Fx patterns
6. **Stay updated** - Follow Fx releases and best practices

---

## 💡 **Pro Tips for Learning Fx**

1. **Start small** - Don't refactor your entire app on day 1
2. **Read the source** - Fx's codebase is well-documented
3. **Use fxtest extensively** - It makes testing much easier
4. **Don't over-engineer** - Not every project needs Fx
5. **Join the community** - Ask questions, share learnings
6. **Practice daily** - Consistency beats intensity
7. **Build real projects** - Theory only goes so far

---

## 🚀 **Your Next Steps**

**Day 1 Action Item**: Set up your development environment and complete Exercise 1 (Manual DI). Share your code with the community for feedback.

**End Goal**: By Day 28, you should be able to confidently build, test, and deploy production-ready Go applications using Fx, understanding when to use it and when simpler approaches suffice.

Good luck on your Fx learning journey! 🎉