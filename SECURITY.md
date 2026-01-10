# Security Policy

This document describes the security practices, expectations, and reporting process for Fortress API.

## Scope

Applies to the codebase and deployment artifacts in this repository. Stack: Go 1.25, Chi, PostgreSQL, Redis, RabbitMQ, gRPC, Docker.

## Reporting a Vulnerability

- Preferred: Open a private GitHub Security Advisory or email harmeetsinghfbd@gmail.com with "Security Vulnerability" in the subject.
- Include:
  1. Affected component/module
  2. Severity (Low, Medium, High, Critical)
  3. Proof-of-concept (code, screenshots, videos)
  4. Steps to reproduce
  5. Affected versions (if known)
  6. Potential mitigations (if any)
- Response timeline:
  - Acknowledgement: within 48 hours
  - Initial triage: within 5 business days
  - Fix/mitigation plan: communicated after triage

Do NOT post exploit details publicly until a fix or mitigation is available.

## Responsible Disclosure

- Avoid mass disclosure. Coordinate with maintainers to allow time for remediation.
- If you require PGP or alternate secure communication, request it in the initial contact.

## Supported Versions

- Only actively maintained branches and releases are covered. Pull requests addressing security will be prioritized.

## Secure Development Practices (Go-specific)

### Input Validation & Sanitization
- Use `go-playground/validator` for all request validation (tags on struct fields)
- Validate path parameters, query strings, and request bodies in controllers
- Never trust user input; validate before passing to services
- Use parameterized queries (SQLc handles this natively)
- Implement custom validation errors with context-aware messages

```go
type CreateUserRequest struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required,min=2,max=100"`
}

if err := validator.Struct(req); err != nil {
    // Return user-friendly error without exposing internals
    logger.Error("validation failed", zap.Error(err))
}
```

### Error Handling
- Use Zap for structured logging; never log sensitive data
- Return generic error messages to clients
- Log full error context internally with structured fields
- Use error wrapping with `fmt.Errorf("%w", err)` for context preservation
- Implement custom error types for specific failure scenarios

```go
if err != nil {
    logger.Error("database operation failed", 
        zap.Error(err), 
        zap.String("operation", "create_user"),
    )
    // Return generic error to client
    http.Error(w, "internal server error", http.StatusInternalServerError)
    return
}
```

### HTTP Security Headers
- Implement middleware to set security headers:
  - `X-Content-Type-Options: nosniff`
  - `X-Frame-Options: DENY`
  - `X-XSS-Protection: 1; mode=block`
  - `Strict-Transport-Security: max-age=31536000; includeSubDomains`
  - `Content-Security-Policy` as needed

```go
func SecurityHeadersMiddleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("X-Content-Type-Options", "nosniff")
            w.Header().Set("X-Frame-Options", "DENY")
            w.Header().Set("X-XSS-Protection", "1; mode=block")
            w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
            next.ServeHTTP(w, r)
        })
    }
}
```

### CORS Configuration
- Implement strict CORS policies
- Only allow necessary origins in production
- Validate `Origin` header on the server

```go
router.Use(cors.Handler(cors.Options{
    AllowedOrigins:   []string{"https://example.com"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
    AllowedHeaders:   []string{"Authorization", "Content-Type"},
    ExposedHeaders:   []string{"Content-Length"},
    MaxAge:           300,
    AllowCredentials: true,
}))
```

## Authentication & Secrets

- JWT secrets, database credentials, API keys, and other secrets MUST be stored in environment variables (`.env.*`); never commit secrets
- Use `godotenv` to load environment variables in development
- Never log sensitive data (passwords, tokens, API keys)
- Implement secret rotation policies
- Use short-lived tokens (JWT exp: 15-30 minutes) with refresh tokens
- Store tokens securely in HTTP-only cookies when possible

```go
// Load from environment
dbPassword := os.Getenv("DB_PASSWORD")
jwtSecret := os.Getenv("JWT_SECRET")

if jwtSecret == "" {
    log.Fatal("JWT_SECRET not set")
}

// Sign tokens with expiration
claims := jwt.MapClaims{
    "sub": userID,
    "exp": time.Now().Add(15 * time.Minute).Unix(),
}
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
```

## Passwords & Tokens

- Hash passwords using `golang.org/x/crypto/bcrypt`
- Use minimum 12 rounds for bcrypt (default is 10, consider 12+)
- Implement rate limiting on authentication endpoints (use middleware)
- Use refresh token rotation with short-lived access tokens
- Invalidate sessions on logout (use Redis for blacklisting)

```go
import "golang.org/x/crypto/bcrypt"

// Hash password
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
    return fmt.Errorf("hash password: %w", err)
}

// Verify password
err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(inputPassword))
if err != nil {
    // Invalid password
}
```

## Database Security (PostgreSQL)

- Use SQLc for type-safe queries (prevents SQL injection)
- Never construct SQL strings dynamically
- Use prepared statements (SQLc enforces this)
- Implement least-privilege database users:
  - App user: SELECT, INSERT, UPDATE, DELETE on specific tables
  - Migration user: Full permissions (used only in migrations)
- Enable row-level security (RLS) for sensitive data
- Use connection pooling with pgx
- Encrypt database connections with TLS

```go
// Connection string with SSL
connStr := "postgres://user:password@localhost/dbname?sslmode=require"
config, err := pgx.ParseConfig(connStr)
if err != nil {
    return fmt.Errorf("parse config: %w", err)
}

// Connection pooling
pool, err := pgxpool.NewWithConfig(ctx, config)
```

## Data Protection

- Encrypt sensitive data at rest (database-level encryption for PostgreSQL)
- Use TLS 1.3+ for all in-transit communication; enforce HTTPS in production
- Implement encryption for sensitive fields (e.g., PII) at application level using `crypto/aes`
- Use proper key derivation functions (PBKDF2, Argon2)
- Regularly back up data and test recovery procedures

```go
import "crypto/aes"

// Encrypt sensitive data
encryptedData, err := encryptAES(plaintext, encryptionKey)
if err != nil {
    return fmt.Errorf("encrypt: %w", err)
}
```

## PostgreSQL Hardening

- Require authentication; use strong passwords or certificate-based auth
- Bind to private network interfaces (not 0.0.0.0)
- Enable TLS for client connections
- Use least-privilege user roles
- Enable logging of failed authentication attempts
- Run with latest security patches
- Implement connection pooling limits

## Redis Security

- Require authentication (use `requirepass` in production)
- Bind to private network interfaces only
- Enable TLS for production deployments
- Use Redis ACLs for role-based access control
- Disable dangerous commands (FLUSHDB, FLUSHALL, KEYS, etc.) or restrict to trusted sources
- Enable persistence and encryption at rest if sensitive data is cached

## RabbitMQ Security

- Require authentication for all users
- Use strong passwords or certificate-based authentication
- Bind to private network interfaces
- Enable TLS/SSL for all connections
- Implement vhost isolation
- Use least-privilege user permissions per queue
- Monitor connection attempts and message patterns

## Container & Host Hardening

- Use distroless or scratch base images to minimize attack surface
- Run containers as non-root users (distroless enforces this)
- Scan images for vulnerabilities using Trivy or similar tools
- Keep base images updated; monitor and apply security patches
- Set resource limits (memory, CPU) on containers
- Use read-only filesystems where possible
- Implement security scanning in CI/CD pipelines

```dockerfile
# Use distroless for production
FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=builder /app/fortress-api /app/fortress-api

# Runs as non-root automatically
ENTRYPOINT ["/app/fortress-api"]
```

## Dependency Management

- Use `go mod` for dependency management
- Run `go mod verify` to check integrity
- Regularly update dependencies: `go get -u ./...`
- Audit for vulnerabilities using:
  - `go list -json -m all | nancy sleuth` (Nancy)
  - `govulncheck ./...` (Go vulnerability scanner)
  - Snyk integration for continuous monitoring
- Remove unused dependencies: `go mod tidy`
- Pin versions in go.mod for reproducible builds
- Implement automated security scanning in CI/CD

```bash
# Check for vulnerabilities
govulncheck ./...

# Tidy and verify
go mod tidy && go mod verify

# Update dependencies safely
go get -u -t ./...
go mod tidy
```

## CI / CD Security

- Store secrets in GitHub Secrets or CI secret managers; never commit secrets
- Enforce branch protection rules and require PR reviews
- Run security checks before merge:
  - `go vet ./...` - Static analysis
  - `govulncheck ./...` - Vulnerability scanning
  - `go test -race ./...` - Race condition detection
  - SAST tools (e.g., Snyk, SonarQube)
- Use signed commits (`git commit -S`)
- Implement artifact signing and verification
- Run container image vulnerability scans in CI

```yaml
# Example GitHub Actions workflow
- name: Security Checks
  run: |
    go vet ./...
    govulncheck ./...
    go test -race ./...
    trivy image --severity HIGH,CRITICAL $IMAGE_NAME
```

## Logging & Monitoring

- Use structured logging (Zap) without exposing secrets
- Never log passwords, tokens, API keys, or PII
- Log security events: failed auth, privilege escalation, unusual patterns
- Implement centralized logging for audit trails
- Monitor authentication failures and rate-limit violations
- Set up alerts for:
  - Multiple failed login attempts
  - Unauthorized API access
  - Database connection errors
  - Service crashes or panics
- Use metrics (Prometheus) for performance and security monitoring

```go
// Structured logging with Zap
logger.Info("user authentication successful",
    zap.String("user_id", user.ID),
    zap.String("method", "password"),
)

logger.Warn("authentication failed",
    zap.String("email", email),
    zap.Error(err),
    zap.Time("timestamp", time.Now()),
)

// Never log sensitive data
// logger.Info("user login", zap.String("password", password)) // ❌ WRONG
```

## API Security

- Implement rate limiting per IP/user
- Use API versioning for backward compatibility
- Validate Content-Type headers
- Implement request signing for critical operations
- Use API keys with scopes and rotation policies
- Implement request logging without sensitive data
- Use consistent authentication (Bearer tokens with JWT)

```go
// Rate limiting middleware
func RateLimitMiddleware(limiter *rate.Limiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

## gRPC Security

- Use TLS for all gRPC connections
- Implement authentication (mutual TLS or token-based)
- Validate input in gRPC handlers
- Use message size limits to prevent DoS
- Implement rate limiting on gRPC services
- Log gRPC errors without exposing internals

```go
// Load certificates
creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
if err != nil {
    log.Fatal(err)
}

// Create secure gRPC server
server := grpc.NewServer(grpc.Creds(creds))
```

## Incident Response

- Triage and contain: revoke affected keys, rotate credentials, apply hotfix.
- Notify affected users and stakeholders per legal/regulatory requirements.
- Post-incident: root cause analysis and preventative measures documented.

## Automated Tools & Tests

- Use static analysis, dependency scanning, and secret detection in pre-commit / CI.
- Add unit/integration tests for security-critical flows (auth, token rotation, password reset).

## Acknowledgements

Following guidelines and patterns from the project's architecture: repository pattern, singleton connections, and centralized utilities for errors, logging, and responses.

For questions or to report issues: open a private GitHub Security Advisory or contact the project maintainers.
