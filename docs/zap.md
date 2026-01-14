# Complete Zap Logging Learning Plan - Zero to Production

---

### **Day 2: Structured Fields**

#### **Theory: Adding Context with Fields**


#### **Performance Tip: Avoid zap.Any()**

```go
// ❌ SLOW: Uses reflection
logger.Info("User data", zap.Any("user", user))

// ✅ FAST: Strongly typed
logger.Info("User data",
	zap.Int("user_id", user.ID),
	zap.String("username", user.Username),
)

// ✅ BETTER: Custom marshaler (Day 5)
logger.Info("User data", zap.Object("user", user))
```

**Tasks**:
1. Log all field types
2. Benchmark `zap.Any()` vs typed fields
3. Create a helper function for common field combinations
4. Practice using `zap.Namespace()` for grouping

---

### **Day 3: Log Levels & Sampling**

#### **Theory: Log Levels**

Zap has 6 log levels (in order):
1. **Debug**: Detailed debugging information
2. **Info**: General informational messages
3. **Warn**: Warning messages (something unexpected but not critical)
4. **Error**: Error messages (operation failed but app continues)
5. **DPanic**: Development panic (panics in dev, logs error in prod)
6. **Panic**: Logs and panics
7. **Fatal**: Logs and calls os.Exit(1)

#### **Practice: Level Control**



---

### **Day 4: Logger Configuration**

#### **Theory: Configuration Options**

Zap is highly configurable through `zap.Config`:

```go
type Config struct {
	Level            AtomicLevel      // Minimum log level
	Development      bool             // Dev mode (more verbose)
	Encoding         string           // "json" or "console"
	EncoderConfig    EncoderConfig    // How to format logs
	OutputPaths      []string         // Where to write (stdout, file)
	ErrorOutputPaths []string         // Where to write internal errors
	InitialFields    map[string]interface{} // Fields in every log
	Sampling         *SamplingConfig  // Sampling configuration
}
```

#### **Practice: Custom Configuration**

```go
// Exercise 5: Full custom configuration
package main

import (
	"os"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// 1. Basic custom config
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout", "/log/app.log"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields: map[string]interface{}{
			"app":     "Fortress_API",
			"version": "1.0.0",
		},
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("Application started")

	// 2. Console encoding for development
	consoleConfig := zap.NewDevelopmentConfig()
	consoleConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleLogger, _ := consoleConfig.Build()

	consoleLogger.Info("Colorful console output!")
	consoleLogger.Warn("Warning in color")
	consoleLogger.Error("Error in color")


	// 4. Environment-based config
	var logger2 *zap.Logger
	env := os.Getenv("ENV") // "development" or "production"
	
	if env == "production" {
		logger2, _ = zap.NewProduction()
	} else {
		logger2, _ = zap.NewDevelopment()
	}
	defer logger2.Sync()

	logger2.Info("Environment-aware logging")
}
```

**Tasks**:
1. Create custom encoder configuration
2. Log to multiple destinations (stdout + file)
3. Add global fields to all logs
4. Configure based on environment variables

---



### **Day 6-7: Week 1 Project - REST API with Comprehensive Logging**

**Goal**: Build a complete REST API with best-practice logging

```go
// Project: User Management API with Zap
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger middleware
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Create request-scoped logger
			reqLogger := logger.With(
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("request_id", generateRequestID()),
			)

			// Add logger to context
			ctx := context.WithValue(r.Context(), "logger", reqLogger)
			r = r.WithContext(ctx)

			reqLogger.Info("Request started")

			// Wrap response writer to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			reqLogger.Info("Request completed",
				zap.Int("status", wrapped.statusCode),
				zap.Duration("duration", duration),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Repository with logging
type UserRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewUserRepository(db *sql.DB, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger.Named("repository"),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
	logger := getLogger(ctx).With(zap.String("operation", "create_user"))
	
	logger.Debug("Creating user", zap.String("username", user.Username))
	
	start := time.Now()
	_, err := r.db.ExecContext(ctx, 
		"INSERT INTO users (username, email) VALUES (?, ?)",
		user.Username, user.Email,
	)
	duration := time.Since(start)

	if err != nil {
		logger.Error("Failed to create user",
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return err
	}

	logger.Info("User created successfully",
		zap.String("username", user.Username),
		zap.Duration("duration", duration),
	)
	return nil
}

// Handler
type UserHandler struct {
	repo   *UserRepository
	logger *zap.Logger
}

func NewUserHandler(repo *UserRepository, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		repo:   repo,
		logger: logger.Named("handler"),
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	logger := getLogger(r.Context())

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Warn("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), &user); err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Helper to get logger from context
func getLogger(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value("logger").(*zap.Logger); ok {
		return logger
	}
	logger, _ := zap.NewProduction()
	return logger
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func generateRequestID() string {
	return "req-123"
}

func main() {
	// Setup logger
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()
	defer logger.Sync()

	// Setup database
	db, _ := sql.Open("sqlite3", ":memory:")
	db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT, email TEXT)")

	// Setup handlers
	repo := NewUserRepository(db, logger)
	handler := NewUserHandler(repo, logger)

	// Setup router with middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/users", handler.Create)

	// Wrap with logging middleware
	loggedMux := LoggingMiddleware(logger)(mux)

	logger.Info("Starting server", zap.String("addr", ":8080"))
	http.ListenAndServe(":8080", loggedMux)
}
```

**Requirements**:
- [ ] HTTP request/response logging
- [ ] Database query logging with duration
- [ ] Error logging with stack traces
- [ ] Request-scoped loggers with correlation IDs
- [ ] Different log levels for different operations
- [ ] Structured JSON output
- [ ] Performance metrics (query time, request duration)

---

## 📚 **Week 2: Advanced Patterns & Production**

For performance, implement `zapcore.ObjectMarshaler` instead of using `zap.Any()`:
---

### **Day 8: Custom Encoders & Formatters**

### **Day 9: Integration with Popular Frameworks**

#### **Echo Framework**

```go
// Exercise 9: Zap with Echo
package main

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"time"
)

func EchoLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			req := c.Request()
			res := c.Response()

			logger.Info("Request",
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.Int("status", res.Status),
				zap.Duration("latency", time.Since(start)),
				zap.String("ip", c.RealIP()),
			)

			return err
		}
	}
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	e := echo.New()
	e.Use(EchoLogger(logger))

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello!")
	})

	e.Start(":8080")
}
```


### **Day 10: Error Logging & Stack Traces**

#### **Theory: Proper Error Logging**

```go
// Exercise 11: Advanced error logging
package main

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
)

// Custom error types
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type DatabaseError struct {
	Operation string
	Err       error
}

func (e DatabaseError) Error() string {
	return fmt.Sprintf("database %s failed: %v", e.Operation, e.Err)
}

func (e DatabaseError) Unwrap() error {
	return e.Err
}

// Implement zapcore.ObjectMarshaler for rich error logging
func (e DatabaseError) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("operation", e.Operation)
	enc.AddString("error", e.Err.Error())
	return nil
}

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// 1. Simple error logging
	err := errors.New("something went wrong")
	logger.Error("Operation failed", zap.Error(err))

	// 2. Wrapped errors
	originalErr := errors.New("connection refused")
	wrappedErr := fmt.Errorf("failed to connect to database: %w", originalErr)
	logger.Error("Database error", zap.Error(wrappedErr))

	// 3. Custom error types
	validationErr := ValidationError{
		Field:   "email",
		Message: "invalid email format",
	}
	logger.Warn("Validation failed", zap.Error(validationErr))

	// 4. Structured error logging
	dbErr := DatabaseError{
		Operation: "insert",
		Err:       originalErr,
	}
	logger.Error("Database operation failed",
		zap.Object("db_error", dbErr),
		zap.String("table", "users"),
	)

	// 5. Stack traces
	logger.Error("Critical error",
		zap.Error(err),
		zap.Stack("stacktrace"), // Adds full stack trace
	)

	// 6. Multiple errors
	errs := []error{
		errors.New("error 1"),
		errors.New("error 2"),
		errors.New("error 3"),
	}
	logger.Error("Multiple errors occurred",
		zap.Errors("errors", errs),
	)

	// 7. Error with context
	logger.Error("Failed to process user",
		zap.Error(err),
		zap.String("user_id", "123"),
		zap.String("operation", "update_profile"),
		zap.Int("retry_count", 3),
	)
}
```

### **Day 11: Sampling & High-Performance Logging**

#### **Theory: Sampling for High-Frequency Logs**

When you have logs that fire thousands of times per second, sampling prevents log spam:

```go
// Exercise 13: Sampling strategies
package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func main() {
	// 1. Basic sampling
	config := zap.NewProductionConfig()
	config.Sampling = &zap.SamplingConfig{
		Initial:    100, // Log first 100 entries per second
		Thereafter: 100, // Then log every 100th entry
	}
	logger, _ := config.Build()
	defer logger.Sync()

	// Simulate high-frequency logging
	for i := 0; i < 10000; i++ {
		logger.Info("High frequency log", zap.Int("iteration", i))
	}

	// 2. Custom sampling decision
	core := zapcore.NewSamplerWithOptions(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			zapcore.InfoLevel,
		),
		time.Second, // Sample per second
		100,         // Initial
		10,          // Thereafter
	)

	sampledLogger := zap.New(core)
	
	// 3. Different sampling for different levels
	// Only sample Info and Debug, never sample Error
	dynamicSampler := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewSamplerWithOptions(
			core,
			time.Second,
			100,
			10,
			zapcore.SamplerHook(func(entry zapcore.Entry, decision zapcore.SamplingDecision) {
				// Never sample errors
				if entry.Level >= zapcore.ErrorLevel {
					decision = zapcore.LogDropped // Override to always log
				}
			}),
		)
	})

	configWithDynamicSampling := zap.NewProductionConfig()
	logger2, _ := configWithDynamicSampling.Build(dynamicSampler)

	// All errors are logged
	for i := 0; i < 1000; i++ {
		logger2.Error("Critical error", zap.Int("iteration", i))
	}

	// Info is sampled
	for i := 0; i < 1000; i++ {
		logger2.Info("Info message", zap.Int("iteration", i))
	}
}
```

#### **Zero-Allocation Logging**

```go
// Exercise 14: Performance optimization
package main

import (
	"go.uber.org/zap"
	"testing"
)

// Pre-allocate fields for hot paths
var (
	userIDField = zap.String("user_id", "")
	actionField = zap.String("action", "")
)

func logUserAction(logger *zap.Logger, userID, action string) {
	// Reuse pre-allocated fields
	logger.Info("User action",
		zap.String("user_id", userID),
		zap.String("action", action),
	)
}

// Benchmark different approaches
func BenchmarkLogging(b *testing.B) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	b.Run("String concatenation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Info("User " + "123" + " performed " + "login")
		}
	})

	b.Run("Structured fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Info("User action",
				zap.String("user_id", "123"),
				zap.String("action", "login"),
			)
		}
	})

	b.Run("zap.Any (slow)", func(b *testing.B) {
		user := struct {
			ID   string
			Name string
		}{"123", "Alice"}

		for i := 0; i < b.N; i++ {
			logger.Info("User action", zap.Any("user", user))
		}
	})

	b.Run("Custom marshaler (fast)", func(b *testing.B) {
		user := User{ID: "123", Name: "Alice"}

		for i := 0; i < b.N; i++ {
			logger.Info("User action", zap.Object("user", user))
		}
	})
}

type User struct {
	ID   string
	Name string
}

func (u User) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("id", u.ID)
	enc.AddString("name", u.Name)
	return nil
}

func main() {
	// Run benchmarks
	// go test -bench=. -benchmem
}
```

**Tasks**:
1. Implement sampling for high-frequency logs
2. Benchmark your logging code
3. Identify and optimize hot paths
4. Create sampling rules based on log level

---

#### **Queryable Logs**

```go
// Exercise 16: Structured logs for queries
package main

import (
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Design logs for common queries

	// Query 1: "Show me all failed login attempts"
	logger.Warn("Login failed",
		zap.String("event_type", "authentication"),
		zap.String("event_action", "login_failed"),
		zap.String("user_email", "alice@example.com"),
		zap.String("failure_reason", "invalid_password"),
		zap.String("ip_address", "192.168.1.1"),
		zap.String("user_agent", "Mozilla/5.0"),
	)
	// Query: event_type:"authentication" AND event_action:"login_failed"

	// Query 2: "Show me all slow database queries"
	logger.Warn("Slow query",
		zap.String("event_type", "database"),
		zap.String("query_type", "SELECT"),
		zap.String("table", "users"),
		zap.Duration("duration", 2500*time.Millisecond),
		zap.Int("rows_affected", 1000),
	)
	// Query: event_type:"database" AND duration:>2000

	// Query 3: "Show me all 5xx errors by endpoint"
	logger.Error("HTTP 500",
		zap.String("event_type", "http"),
		zap.String("method", "POST"),
		zap.String("path", "/api/users"),
		zap.Int("status_code", 500),
		zap.String("error_type", "internal_server_error"),
	)
	// Query: event_type:"http" AND status_code:>=500

	// Query 4: "Show me user journey for user X"
	userID := "user-123"
	logger.Info("Page view",
		zap.String("event_type", "user_action"),
		zap.String("action", "page_view"),
		zap.String("user_id", userID),
		zap.String("page", "/dashboard"),
	)
	logger.Info("Button click",
		zap.String("event_type", "user_action"),
		zap.String("action", "button_click"),
		zap.String("user_id", userID),
		zap.String("button_id", "create_post"),
	)
	// Query: user_id:"user-123" | sort by timestamp
}
```

**Tasks**:
1. Configure Zap for your log management system
2. Design a consistent event taxonomy
3. Create dashboards for common queries
4. Set up alerts based on log patterns

---


**Goal**: Build a complete microservice with world-class logging

```go
// Final Project: E-commerce Order Service
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Configuration
type Config struct {
	Environment string
	LogLevel    string
	Service     string
	Version     string
}

func NewConfig() *Config {
	return &Config{
		Environment: getEnv("ENV", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Service:     "order-service",
		Version:     "1.0.0",
	}
}

// Logger factory
func NewLogger(cfg *Config) (*zap.Logger, error) {
	var zapConfig zap.Config

	if cfg.Environment == "production" {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}

	// Parse log level
	level, err := zapcore.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	// Add global fields
	zapConfig.InitialFields = map[string]interface{}{
		"service":     cfg.Service,
		"version":     cfg.Version,
		"environment": cfg.Environment,
	}

	// Custom encoder config
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := zapConfig.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return logger, err
}

// Request logger middleware
func RequestLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = generateRequestID()
			}

			// Create request-scoped logger
			reqLogger := logger.With(
				zap.String("request_id", requestID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)

			// Add to context
			ctx := context.WithValue(r.Context(), "logger", reqLogger)
			ctx = context.WithValue(ctx, "request_id", requestID)
			r = r.WithContext(ctx)

			reqLogger.Info("Request started")

			// Wrap response writer
			wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			// Log completion
			logFunc := reqLogger.Info
			if wrapped.statusCode >= 500 {
				logFunc = reqLogger.Error
			} else if wrapped.statusCode >= 400 {
				logFunc = reqLogger.Warn
			}

			logFunc("Request completed",
				zap.Int("status_code", wrapped.statusCode),
				zap.Duration("duration", duration),
				zap.Int("response_size", wrapped.size),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// Order domain
type Order struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (o Order) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("id", o.ID)
	enc.AddString("customer_id", o.CustomerID)
	enc.AddFloat64("amount", o.Amount)
	enc.AddString("status", o.Status)
	enc.AddTime("created_at", o.CreatedAt)
	return nil
}

// Repository
type OrderRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewOrderRepository(db *sql.DB, logger *zap.Logger) *OrderRepository {
	return &OrderRepository{
		db:     db,
		logger: logger.Named("repository"),
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *Order) error {
	logger := getLogger(ctx).With(zap.Object("order", *order))

	logger.Debug("Creating order")
	start := time.Now()

	_, err := r.db.ExecContext(ctx,
		"INSERT INTO orders (id, customer_id, amount, status, created_at) VALUES (?, ?, ?, ?, ?)",
		order.ID, order.CustomerID, order.Amount, order.Status, order.CreatedAt,
	)

	duration := time.Since(start)

	if err != nil {
		logger.Error("Failed to create order",
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return err
	}

	logger.Info("Order created successfully",
		zap.Duration("duration", duration),
	)

	return nil
}

// Service
type OrderService struct {
	repo   *OrderRepository
	logger *zap.Logger
}

func NewOrderService(repo *OrderRepository, logger *zap.Logger) *OrderService {
	return &OrderService{
		repo:   repo,
		logger: logger.Named("service"),
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, customerID string, amount float64) (*Order, error) {
	logger := getLogger(ctx).With(
		zap.String("customer_id", customerID),
		zap.Float64("amount", amount),
	)

	logger.Info("Processing order creation")

	// Validation
	if amount <= 0 {
		logger.Warn("Invalid order amount",
			zap.String("validation_error", "amount_must_be_positive"),
		)
		return nil, errors.New("amount must be positive")
	}

	order := &Order{
		ID:         generateOrderID(),
		CustomerID: customerID,
		Amount:     amount,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, order); err != nil {
		logger.Error("Order

#### **gRPC**

```go
// Exercise 10: Zap with gRPC
package main

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		logger.Debug("gRPC request started",
			zap.String("method", info.FullMethod),
		)

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		if err != nil {
			logger.Error("gRPC request failed",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
				zap.Error(err),

			)
		} else {
			logger.Info("gRPC request completed",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
			)
		}

		return resp, err
	}
}

func StreamServerInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()

		logger.Debug("gRPC stream started",
			zap.String("method", info.FullMethod),
			zap.Bool("is_client_stream", info.IsClientStream),
			zap.Bool("is_server_stream", info.IsServerStream),
		)

		err := handler(srv, ss)

		duration := time.Since(start)

		if err != nil {
			logger.Error("gRPC stream failed",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
				zap.Error(err),
			)
		} else {
			logger.Info("gRPC stream completed",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
			)
		}

		return err
	}
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	server := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(logger)),
		grpc.StreamInterceptor(StreamServerInterceptor(logger)),
	)

	// Register your services...
	// server.Serve(listener)
}
```

**Tasks**:
1. Integrate Zap with your framework of choice
2. Create middleware for request/response logging
3. Add panic recovery with stack traces
4. Log slow requests (> 1 second)

---
