# Architecture Documentation

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                          Client Layer                            │
│                     (Web, Mobile, CLI)                           │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         │ HTTP/HTTPS
                         │
┌────────────────────────▼────────────────────────────────────────┐
│                      API Gateway / Load Balancer                 │
└────────────────────────┬────────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────────┐
│                    Application Layer (Echo)                      │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │              Middleware Stack                             │   │
│  │  • Recovery  • Request ID  • Logger  • CORS  • Security  │   │
│  └──────────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                  Handler Layer                            │   │
│  │  • Request Validation  • Response Formatting              │   │
│  └──────────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                  Service Layer                            │   │
│  │  • Business Logic  • Orchestration  • Validation         │   │
│  └──────────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                Repository Layer                           │   │
│  │  • Data Access  • Query Execution  • Transactions        │   │
│  └──────────────────────────────────────────────────────────┘   │
└────────────┬──────────────────────────┬────────────────────────┘
             │                          │
             │                          │
┌────────────▼─────────────┐   ┌────────▼──────────┐
│   PostgreSQL Database    │   │   Redis Cache     │
│   • Categories           │   │   • Sessions      │
│   • Todos                │   │   • Rate Limits   │
│   • Comments             │   │   • Asynq Jobs    │
└──────────────────────────┘   └───────────────────┘
```

## Feature-Based Architecture

Each feature follows a consistent layered pattern:

```
Feature (e.g., Category)
├── dto.go          - Request/Response DTOs with validation tags
├── model.go        - Model mappers (DB → DTO conversion)
├── repository.go   - Data access layer (CRUD operations)
├── service.go      - Business logic layer
└── handler.go      - HTTP request handlers with route registration
```

### Data Flow

```
HTTP Request
    │
    ▼
Handler (Validation, Request Parsing)
    │
    ▼
Service (Business Logic, Orchestration)
    │
    ▼
Repository (Database Operations)
    │
    ▼
Database (PostgreSQL)
    │
    ▼
Repository (Return Data)
    │
    ▼
Service (Transform Data)
    │
    ▼
Handler (Format Response)
    │
    ▼
HTTP Response
```

## Dependency Injection Flow

```
main.go
  │
  ├─► Logger Service
  │     │
  │     ├─► Category Repository ──┐
  │     ├─► Todo Repository ──────┤
  │     └─► Comment Repository ───┤
  │                                │
  ├─► Validator                    │
  │                                │
  ├─► Category Service ◄───────────┤
  │     │                          │
  ├─► Todo Service ◄───────────────┤
  │     │                          │
  ├─► Comment Service ◄────────────┘
  │     │
  ├─► Category Handler ◄─── Category Service + Validator
  ├─► Todo Handler ◄─────── Todo Service + Validator
  └─► Comment Handler ◄──── Comment Service + Validator
        │
        ▼
    Echo Server (Routes Registration)
```

## Error Handling Architecture

```
Error Occurs in Repository/Service
    │
    ▼
Custom AppError Created
    │
    ▼
Bubbled Up Through Layers
    │
    ▼
Global Error Handler Middleware
    │
    ├─► Log Error (if 5xx)
    │
    ├─► Map to HTTP Status Code
    │
    ├─► Create Error Response
    │
    └─► Send JSON Response
```

## Database Schema

### Categories Table
```sql
- id (UUID, PK)
- name (VARCHAR, NOT NULL)
- description (TEXT)
- color (VARCHAR)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

### Todos Table
```sql
- id (UUID, PK)
- title (VARCHAR, NOT NULL)
- description (TEXT)
- status (ENUM: pending, in_progress, completed, cancelled)
- priority (ENUM: low, medium, high, urgent)
- category_id (UUID, FK → categories.id)
- due_date (TIMESTAMP)
- completed_at (TIMESTAMP)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

### Comments Table
```sql
- id (UUID, PK)
- todo_id (UUID, FK → todos.id, CASCADE DELETE)
- content (TEXT, NOT NULL)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

## Configuration Management

```
Configuration Sources (Priority Order):
1. Environment Variables (Highest)
2. config.yaml file
3. Default Values (Lowest)

Koanf Library
    │
    ├─► Load config.yaml
    │
    ├─► Override with ENV vars
    │
    ├─► Unmarshal to Config struct
    │
    └─► Validate with go-playground/validator
```

## Middleware Stack (Execution Order)

```
Request
  │
  ▼
1. Recover (Panic Recovery)
  │
  ▼
2. Request ID (Generate/Extract)
  │
  ▼
3. Logger (Log Request)
  │
  ▼
4. CORS (Handle CORS)
  │
  ▼
5. Security Headers
  │
  ▼
6. Body Limit (2MB)
  │
  ▼
7. Gzip Compression
  │
  ▼
Handler Execution
  │
  ▼
Response
```

## Connection Pooling

### Database (pgx)
- Max Open Connections: 25 (configurable)
- Max Idle Connections: 5 (configurable)
- Connection Lifetime: 5 minutes
- Health Check: Every 1 minute

### Redis
- Pool Size: 10 (configurable)
- Min Idle Connections: 5
- Max Retries: 3
- Dial/Read/Write Timeout: 3-5 seconds

## Graceful Shutdown Flow

```
SIGINT/SIGTERM Signal Received
    │
    ▼
Server.Shutdown() Called
    │
    ├─► Stop Accepting New Connections
    │
    ├─► Wait for Active Requests (max: shutdown_timeout)
    │
    ├─► Close HTTP Server
    │
    ├─► Close Database Connections
    │
    ├─► Close Redis Connections
    │
    └─► Exit Application
```

## Security Considerations

1. **Input Validation**: All requests validated before processing
2. **SQL Injection**: SQLC generates type-safe queries
3. **CORS**: Configurable CORS policy
4. **Rate Limiting**: Can be enabled per endpoint
5. **Request Size Limit**: 2MB body limit
6. **Timeout Protection**: Read/Write timeouts configured
7. **Panic Recovery**: Automatic recovery from panics
8. **Structured Logging**: No sensitive data in logs

## Performance Optimizations

1. **Connection Pooling**: Efficient database connection reuse
2. **Query Optimization**: Indexed columns for common queries
3. **Response Compression**: Gzip compression enabled
4. **Minimal Allocations**: Efficient memory usage
5. **Context Propagation**: Request context throughout the stack
6. **Prepared Statements**: Query compilation caching

## Scalability Patterns

### Horizontal Scaling
- Stateless application design
- No session storage in application memory
- Redis for distributed caching
- Database connection pooling

### Vertical Scaling
- Configurable connection pools
- Adjustable worker goroutines
- Memory-efficient operations

## Monitoring & Observability

1. **Health Checks**: `/health` endpoint
2. **Request IDs**: Correlation across services
3. **Structured Logging**: JSON logs for aggregation
4. **Error Tracking**: Detailed error logging
5. **Metrics**: Ready for Prometheus integration
6. **Tracing**: Request ID propagation

## Testing Strategy

```
Unit Tests
  ├─► Repository Layer (Mock DB)
  ├─► Service Layer (Mock Repository)
  └─► Handler Layer (Mock Service)

Integration Tests
  ├─► Database Integration
  ├─► Redis Integration
  └─► API Endpoint Tests

E2E Tests
  └─► Full Request/Response Cycle
```

## Deployment Architecture

```
Production Environment:
  ├─► Load Balancer (nginx/HAProxy)
  │     │
  │     ├─► App Instance 1
  │     ├─► App Instance 2
  │     └─► App Instance N
  │
  ├─► PostgreSQL (Primary + Replica)
  │
  ├─► Redis (Primary + Sentinel)
  │
  └─► Background Workers (Asynq)
```

## Future Enhancements

1. **Authentication**: JWT-based auth
2. **Authorization**: Role-based access control
3. **Pagination**: Cursor-based pagination
4. **WebSockets**: Real-time updates
5. **GraphQL**: Alternative API interface
6. **Caching**: Redis caching layer
7. **Background Jobs**: Async task processing
8. **Email Notifications**: Resend integration
9. **File Uploads**: S3 integration
10. **API Versioning**: Multiple API versions
