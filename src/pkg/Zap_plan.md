# Complete Zap Logging Learning Plan - Zero to Production

A comprehensive 2-week plan to master Uber's Zap logging library.

---

## 📚 **Week 1: Fundamentals & Core Concepts**

---

### **Day 1: Understanding Structured Logging**

#### **Theory: Why Zap?**

**Problems with traditional logging**:
```go
// Traditional logging (fmt/log package)
log.Printf("User %s logged in from %s at %v", username, ip, time.Now())
// Output: User alice logged in from 192.168.1.1 at 2024-11-15 10:30:45

// Problems:
// ❌ Hard to parse programmatically
// ❌ No log levels
// ❌ Poor performance
// ❌ Can't filter by fields
// ❌ No structured querying
```

**Zap's advantages**:
- **Fast**: 10x faster than other Go loggers
- **Structured**: JSON output for easy parsing
- **Type-safe**: Compile-time checking
- **Zero allocations**: In hot paths
- **Levels**: Debug, Info, Warn, Error, Fatal, Panic
- **Fields**: Key-value pairs for context

#### **Practice: First Zap Logger**

```go
// Exercise 1: Basic Zap usage
package main

import (
	"go.uber.org/zap"
)

func main() {
	// Method 1: Quick development logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // Flush any buffered log entries

	logger.Info("Hello from Zap!")
	logger.Debug("This is a debug message")
	logger.Warn("This is a warning")
	logger.Error("This is an error")

	// Method 2: Production logger (JSON output)
	prodLogger, _ := zap.NewProduction()
	defer prodLogger.Sync()

	prodLogger.Info("Production log entry")
}
```

**Output comparison**:
```bash
# Development logger (human-readable)
2024-11-15T10:30:45.123+0530	INFO	main.go:11	Hello from Zap!

# Production logger (JSON)
{"level":"info","ts":1700029845.123,"caller":"main/main.go:11","msg":"Hello from Zap!"}
```

**Tasks**:
1. Create both development and production loggers
2. Log messages at all levels (Debug, Info, Warn, Error)
3. Notice the difference in output format
4. Try `logger.Fatal()` and see what happens (app exits!)

---

### **Day 2: Structured Fields**

#### **Theory: Adding Context with Fields**

Zap uses strongly-typed fields for better performance and safety.

```go
// ❌ BAD: String concatenation
logger.Info("User " + username + " logged in from " + ip)

// ✅ GOOD: Structured fields
logger.Info("User logged in",
	zap.String("username", username),
	zap.String("ip", ip),
)
```

#### **Practice: Field Types**

```go
// Exercise 2: All field types
package main

import (
	"time"
	"go.uber.org/zap"
)

type User struct {
	ID       int
	Username string
	Email    string
}

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// String fields
	logger.Info("String field",
		zap.String("name", "Alice"),
	)

	// Numeric fields
	logger.Info("Numeric fields",
		zap.Int("age", 30),
		zap.Int64("user_id", 123456789),
		zap.Float64("price", 99.99),
	)

	// Boolean fields
	logger.Info("Boolean field",
		zap.Bool("is_admin", true),
	)

	// Time fields
	logger.Info("Time field",
		zap.Time("login_at", time.Now()),
		zap.Duration("request_duration", 250*time.Millisecond),
	)

	// Complex types
	user := User{ID: 1, Username: "alice", Email: "alice@example.com"}
	logger.Info("Complex object",
		zap.Any("user", user), // Uses reflection
		zap.Reflect("user_reflect", user), // Explicit reflection
	)

	// Arrays and slices
	logger.Info("Collections",
		zap.Strings("tags", []string{"golang", "logging", "zap"}),
		zap.Ints("scores", []int{85, 90, 95}),
	)

	// Error fields (special)
	err := performOperation()
	if err != nil {
		logger.Error("Operation failed",
			zap.Error(err), // Special error field
		)
	}

	// Binary data
	logger.Info("Binary data",
		zap.Binary("data", []byte{0x01, 0x02, 0x03}),
	)

	// Namespace (group related fields)
	logger.Info("Grouped fields",
		zap.Namespace("request"),
		zap.String("method", "POST"),
		zap.String("path", "/api/users"),
		zap.Int("status", 200),
	)
}

func performOperation() error {
	return nil
}
```

**Output (Development)**:
```
2024-11-15T10:30:45.123+0530	INFO	main.go:15	String field	{"name": "Alice"}
2024-11-15T10:30:45.124+0530	INFO	main.go:19	Numeric fields	{"age": 30, "user_id": 123456789, "price": 99.99}
```

**Output (Production JSON)**:
```json
{
  "level": "info",
  "ts": 1700029845.123,
  "caller": "main/main.go:15",
  "msg": "String field",
  "name": "Alice"
}
```

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

```go
// Exercise 3: Log levels and filtering
package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// 1. Default loggers
	devLogger, _ := zap.NewDevelopment()  // Debug level and above
	prodLogger, _ := zap.NewProduction()  // Info level and above

	// 2. Custom level
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel) // Only Warn and above
	logger, _ := config.Build()
	defer logger.Sync()

	logger.Debug("This won't appear")
	logger.Info("This won't appear either")
	logger.Warn("This will appear")
	logger.Error("This will appear too")

	// 3. Dynamic level changes
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.InfoLevel)

	dynamicConfig := zap.NewProductionConfig()
	dynamicConfig.Level = atomicLevel
	dynamicLogger, _ := dynamicConfig.Build()

	dynamicLogger.Debug("Hidden")
	dynamicLogger.Info("Visible")

	// Change level at runtime
	atomicLevel.SetLevel(zapcore.DebugLevel)
	dynamicLogger.Debug("Now visible!")

	// 4. Sampling (reduce log volume)
	samplingConfig := zap.NewProductionConfig()
	samplingConfig.Sampling = &zap.SamplingConfig{
		Initial:    100, // Log first 100 messages
		Thereafter: 10,  // Then log every 10th message
	}
	sampledLogger, _ := samplingConfig.Build()

	// This will log first 100, then every 10th
	for i := 0; i < 1000; i++ {
		sampledLogger.Info("Frequent message", zap.Int("iteration", i))
	}
}
```

#### **HTTP Endpoint for Dynamic Level Changes**

```go
// Exercise 4: Change log level via HTTP
package main

import (
	"net/http"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	atomicLevel := zap.NewAtomicLevel()
	
	config := zap.NewProductionConfig()
	config.Level = atomicLevel
	logger, _ := config.Build()
	defer logger.Sync()

	// HTTP endpoint to change level
	// GET /log/level - Returns current level
	// PUT /log/level?level=debug - Sets level
	http.HandleFunc("/log/level", atomicLevel.ServeHTTP)

	logger.Info("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
```

**Test it**:
```bash
# Get current level
curl http://localhost:8080/log/level

# Set to debug
curl -X PUT http://localhost:8080/log/level?level=debug

# Set to error
curl -X PUT http://localhost:8080/log/level?level=error
```

**Tasks**:
1. Create loggers with different default levels
2. Implement dynamic level changes
3. Add sampling for high-frequency logs
4. Create admin endpoint to change levels in production

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
		OutputPaths:      []string{"stdout", "/var/log/app.log"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields: map[string]interface{}{
			"app":     "my-service",
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

	// 3. Multiple outputs (stdout + file)
	file, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(file),
			zapcore.InfoLevel,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		),
	)

	multiLogger := zap.New(core)
	multiLogger.Info("Logged to both console and file")

	// 4. Environment-based config
	var logger2 *zap.Logger
	env := os.Getenv("APP_ENV") // "development" or "production"
	
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

### **Day 5: Named Loggers & Child Loggers**

#### **Theory: Logger Hierarchy**

Create specialized loggers for different parts of your application:

```go
// Root logger
logger := zap.NewProduction()

// Named loggers
dbLogger := logger.Named("database")
apiLogger := logger.Named("api")
cacheLogger := logger.Named("cache")

// Child loggers with persistent fields
userLogger := logger.With(
	zap.String("user_id", "123"),
	zap.String("session_id", "abc"),
)
```

#### **Practice: Logger Organization**

```go
// Exercise 6: Named and child loggers
package main

import (
	"go.uber.org/zap"
)

type Application struct {
	logger      *zap.Logger
	dbLogger    *zap.Logger
	apiLogger   *zap.Logger
	cacheLogger *zap.Logger
}

func NewApplication() *Application {
	logger, _ := zap.NewProduction()

	return &Application{
		logger:      logger,
		dbLogger:    logger.Named("db"),
		apiLogger:   logger.Named("api"),
		cacheLogger: logger.Named("cache"),
	}
}

func (app *Application) HandleRequest(userID string) {
	// Create request-scoped logger
	reqLogger := app.apiLogger.With(
		zap.String("user_id", userID),
		zap.String("request_id", generateRequestID()),
	)

	reqLogger.Info("Request started")

	// All logs from this point include user_id and request_id
	app.queryDatabase(reqLogger, userID)
	app.checkCache(reqLogger, userID)

	reqLogger.Info("Request completed")
}

func (app *Application) queryDatabase(logger *zap.Logger, userID string) {
	// Database operation
	logger.Info("Querying database")
	// user_id and request_id are automatically included
}

func (app *Application) checkCache(logger *zap.Logger, userID string) {
	// Cache operation
	logger.Info("Checking cache")
	// user_id and request_id are automatically included
}

func generateRequestID() string {
	return "req-12345"
}

func main() {
	app := NewApplication()
	defer app.logger.Sync()

	// Different parts of the app use their named loggers
	app.dbLogger.Info("Database connected")
	app.cacheLogger.Info("Cache initialized")

	// Handle request with scoped context
	app.HandleRequest("user-789")
}
```

**Output (JSON)**:
```json
{"level":"info","ts":1700029845.123,"logger":"db","msg":"Database connected"}
{"level":"info","ts":1700029845.124,"logger":"cache","msg":"Cache initialized"}
{"level":"info","ts":1700029845.125,"logger":"api","msg":"Request started","user_id":"user-789","request_id":"req-12345"}
{"level":"info","ts":1700029845.126,"logger":"api","msg":"Querying database","user_id":"user-789","request_id":"req-12345"}
{"level":"info","ts":1700029845.127,"logger":"api","msg":"Checking cache","user_id":"user-789","request_id":"req-12345"}
{"level":"info","ts":1700029845.128,"logger":"api","msg":"Request completed","user_id":"user-789","request_id":"req-12345"}
```

**Tasks**:
1. Create named loggers for different modules
2. Use `.With()` for request-scoped logging
3. Build a logger middleware for HTTP handlers
4. Implement correlation IDs across services

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

---

### **Day 8: Custom Encoders & Formatters**

#### **Theory: Custom Object Marshaling**

For performance, implement `zapcore.ObjectMarshaler` instead of using `zap.Any()`:

```go
// Exercise 7: Custom object encoding
package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type User struct {
	ID       int
	Username string
	Email    string
	Password string // Sensitive!
}

// Implement zapcore.ObjectMarshaler
func (u User) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("id", u.ID)
	enc.AddString("username", u.Username)
	enc.AddString("email", u.Email)
	// Deliberately omit password
	return nil
}

type HTTPRequest struct {
	Method  string
	Path    string
	Headers map[string]string
}

func (r HTTPRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("method", r.Method)
	enc.AddString("path", r.Path)
	
	// Nest headers
	enc.AddObject("headers", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for k, v := range r.Headers {
			// Redact sensitive headers
			if k == "Authorization" || k == "Cookie" {
				enc.AddString(k, "[REDACTED]")
			} else {
				enc.AddString(k, v)
			}
		}
		return nil
	}))
	
	return nil
}

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	user := User{
		ID:       1,
		Username: "alice",
		Email:    "alice@example.com",
		Password: "secret123",
	}

	// Fast, type-safe logging
	logger.Info("User logged in", zap.Object("user", user))

	req := HTTPRequest{
		Method: "POST",
		Path:   "/api/login",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer token123",
			"User-Agent":    "Mozilla/5.0",
		},
	}

	logger.Info("HTTP request", zap.Object("request", req))
}
```

**Output**:
```json
{
  "level": "info",
  "msg": "User logged in",
  "user": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com"
  }
}

{
  "level": "info",
  "msg": "HTTP request",
  "request": {
    "method": "POST",
    "path": "/api/login",
    "headers": {
      "Content-Type": "application/json",
      "Authorization": "[REDACTED]",
      "User-Agent": "Mozilla/5.0"
    }
  }
}
```

**Tasks**:
1. Implement `MarshalLogObject` for custom types
2. Redact sensitive fields (passwords, tokens)
3. Benchmark vs `zap.Any()`
4. Create reusable marshalers for common types

---

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

#### **Panic Recovery with Logging**

```go
// Exercise 12: Panic recovery
package main

import (
	"go.uber.org/zap"
	"runtime/debug"
)

func RecoverWithLog(logger *zap.Logger) {
	if r := recover(); r != nil {
		logger.Error("Panic recovered",
			zap.Any("panic", r),
			zap.ByteString("stack", debug.Stack()),
		)
	}
}

func riskyOperation(logger *zap.Logger) {
	defer RecoverWithLog(logger)

	// Simulate panic
	panic("something terrible happened")
}

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	riskyOperation(logger)
	logger.Info("Application continues running")
}
```

**Tasks**:
1. Create custom error types with rich logging
2. Implement error wrapping with context
3. Add stack traces for critical errors
4. Create error aggregation for batch operations

---

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

### **Day 12: Log Aggregation & Analysis**

#### **Theory: Integration with Log Management Systems**

```go
// Exercise 15: Structured logs for ELK/Splunk/DataDog
package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func main() {
	// 1. ELK Stack (Elasticsearch, Logstash, Kibana)
	elkConfig := zap.NewProductionConfig()
	elkConfig.EncoderConfig.TimeKey = "@timestamp" // ELK standard
	elkConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	elkConfig.InitialFields = map[string]interface{}{
		"service":     "user-service",
		"environment": "production",
		"version":     "1.2.3",
		"host":        os.Hostname(),
	}
	elkLogger, _ := elkConfig.Build()

	elkLogger.Info("User event",
		zap.String("event_type", "user_login"),
		zap.String("user_id", "123"),
		zap.String("ip_address", "192.168.1.1"),
	)

	// 2. Splunk HEC (HTTP Event Collector)
	splunkConfig := zap.NewProductionConfig()
	splunkConfig.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder, // Unix timestamp
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	splunkConfig.InitialFields = map[string]interface{}{
		"source":     "user-service",
		"sourcetype": "application",
		"index":      "main",
	}
	splunkLogger, _ := splunkConfig.Build()

	splunkLogger.Info("Application metric",
		zap.Int("active_users", 1523),
		zap.Float64("cpu_usage", 45.2),
		zap.Int64("memory_bytes", 1024*1024*500),
	)

	// 3. DataDog
	datadogConfig := zap.NewProductionConfig()
	datadogConfig.EncoderConfig.LevelKey = "status" // DataDog uses "status"
	datadogConfig.EncoderConfig.MessageKey = "message"
	datadogConfig.InitialFields = map[string]interface{}{
		"dd.service": "user-service",
		"dd.env":     "production",
		"dd.version": "1.2.3",
		"dd.trace_id": "1234567890", // Correlation with APM
	}
	datadogLogger, _ := datadogConfig.Build()

	datadogLogger.Info("API call",
		zap.String("http.method", "POST"),
		zap.String("http.url", "/api/users"),
		zap.Int("http.status_code", 201),
		zap.Duration("duration", 45*time.Millisecond),
	)

	// 4. CloudWatch Logs (AWS)
	cloudwatchConfig := zap.NewProductionConfig()
	cloudwatchConfig.EncoderConfig.TimeKey = "timestamp"
	cloudwatchConfig.InitialFields = map[string]interface{}{
		"application": "user-service",
		"environment": "production",
		"region":      "us-east-1",
	}
	cloudwatchLogger, _ := cloudwatchConfig.Build()

	cloudwatchLogger.Info("Lambda invocation",
		zap.String("request_id", "abc-123"),
		zap.String("function_name", "ProcessUser"),
		zap.Duration("duration", 150*time.Millisecond),
		zap.Int("memory_used_mb", 128),
	)
}
```

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

### **Day 13: Testing & Observability**

#### **Theory: Testing Logs**

```go
// Exercise 17: Testing with zap/zaptest
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
)

// Service to test
type UserService struct {
	logger *zap.Logger
}

func NewUserService(logger *zap.Logger) *UserService {
	return &UserService{logger: logger}
}

func (s *UserService) CreateUser(username string) error {
	if username == "" {
		s.logger.Error("Invalid username", zap.String("username", username))
		return errors.New("username required")
	}

	s.logger.Info("User created", zap.String("username", username))
	return nil
}

// Test 1: Using zaptest (logs to testing.T)
func TestUserService_WithZaptest(t *testing.T) {
	logger := zaptest.NewLogger(t)
	service := NewUserService(logger)

	err := service.CreateUser("alice")
	assert.NoError(t, err)
	// Logs appear in test output
}

// Test 2: Using observer (inspect logs)
func TestUserService_WithObserver(t *testing.T) {
	// Create observable logger
	core, recorded := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)
	service := NewUserService(logger)

	// Test successful creation
	err := service.CreateUser("alice")
	assert.NoError(t, err)

	// Verify log was written
	logs := recorded.All()
	assert.Len(t, logs, 1)
	assert.Equal(t, "User created", logs[0].Message)
	assert.Equal(t, "alice", logs[0].ContextMap()["username"])

	// Test failure
	recorded.TakeAll() // Clear previous logs
	err = service.CreateUser("")
	assert.Error(t, err)

	// Verify error was logged
	logs = recorded.All()
	assert.Len(t, logs, 1)
	assert.Equal(t, zapcore.ErrorLevel, logs[0].Level)
	assert.Equal(t, "Invalid username", logs[0].Message)
}

// Test 3: Filtering logs
func TestUserService_FilteredLogs(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)
	service := NewUserService(logger)

	service.CreateUser("alice")
	service.CreateUser("bob")
	service.CreateUser("") // Error

	// Get only error logs
	errorLogs := recorded.FilterLevelExact(zapcore.ErrorLevel).All()
	assert.Len(t, errorLogs, 1)

	// Get only info logs
	infoLogs := recorded.FilterLevelExact(zapcore.InfoLevel).All()
	assert.Len(t, infoLogs, 2)
}

// Test 4: Testing log fields
func TestLogFields(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)

	logger.Info("User action",
		zap.String("user_id", "123"),
		zap.String("action", "login"),
		zap.Int("attempt", 1),
	)

	logs := recorded.All()
	require.Len(t, logs, 1)

	fields := logs[0].ContextMap()
	assert.Equal(t, "123", fields["user_id"])
	assert.Equal(t, "login", fields["action"])
	assert.Equal(t, int64(1), fields["attempt"])
}
```

#### **Integration with Metrics**

```go
// Exercise 18: Logs + Metrics
package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_logs_total",
			Help: "Total number of log entries by level",
		},
		[]string{"level"},
	)

	errorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_errors_total",
			Help: "Total number of errors by type",
		},
		[]string{"error_type"},
	)
)

// Custom core that increments metrics
type metricsCore struct {
	zapcore.Core
}

func (c *metricsCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// Increment log counter
	logCounter.WithLabelValues(entry.Level.String()).Inc()

	// Increment error counter if error level
	if entry.Level >= zapcore.ErrorLevel {
		errorType := "unknown"
		for _, field := range fields {
			if field.Key == "error_type" {
				errorType = field.String
				break
			}
		}
		errorCounter.WithLabelValues(errorType).Inc()
	}

	return c.Core.Write(entry, fields)
}

func NewLoggerWithMetrics() *zap.Logger {
	core := &metricsCore{
		Core: zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			zapcore.InfoLevel,
		),
	}

	return zap.New(core)
}

func main() {
	logger := NewLoggerWithMetrics()
	defer logger.Sync()

	// These logs also increment Prometheus metrics
	logger.Info("Application started")
	logger.Warn("Cache miss")
	logger.Error("Database error", zap.String("error_type", "connection_timeout"))
	logger.Error("API error", zap.String("error_type", "rate_limit"))

	// Prometheus metrics now show:
	// app_logs_total{level="info"} 1
	// app_logs_total{level="warn"} 1
	// app_logs_total{level="error"} 2
	// app_errors_total{error_type="connection_timeout"} 1
	// app_errors_total{error_type="rate_limit"} 1
}
```

**Tasks**:
1. Write unit tests for logged behavior
2. Use observer to inspect logs in tests
3. Create custom cores for metrics integration
4. Set up log-based alerts

---

### **Day 14: Week 2 Project - Production-Grade Microservice**

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