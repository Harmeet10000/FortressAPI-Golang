# Health Check System

Complete health check utilities for monitoring application, system, and dependency health.

## Overview

The health check system provides comprehensive monitoring capabilities for:

- **System Metrics**: CPU usage, memory, platform information
- **Application Metrics**: Uptime, memory usage, goroutines
- **Dependencies**: PostgreSQL, Redis connectivity and performance
- **Resources**: Memory usage percentage, disk accessibility

## Components

### Utility Functions (`utils/health.go`)

#### System Health

```go
// Get system-level metrics (CPU, memory, platform)
systemHealth := utils.GetSystemHealth()

// Returns:o
// {
//   "cpu_usage": [0.5, 0.6, 0.7],
//   "cpu_usage_percent": "25.50%",
//   "total_memory": "8192.00 MB",
//   "free_memory": "4096.00 MB",
//   "platform": "linux",
//   "arch": "amd64"
// }
```

#### Application Health

```go
// Get application-level metrics
appHealth := utils.GetApplicationHealth()

// Returns:
// {
//   "environment": "production",
//   "uptime": "3600.50 seconds",
//   "memory_usage": {
//     "alloc_mb": "128.50",
//     "total_alloc_mb": "256.75",
//     "sys_mb": "512.25",
//     "num_gc": 42,
//     "heap_alloc_mb": "128.50",
//     "heap_sys_mb": "512.25",
//     "stack_alloc_mb": "5.75"
//   },
//   "pid": 12345,
//   "go_version": "go1.25.5"
// }
```

#### Database Check

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Single connection check
dbHealth := utils.CheckDatabase(ctx, singleConn)

// Connection pool check (recommended)
dbHealth := utils.CheckDatabasePool(ctx, pool)

// Returns:
// {
//   "status": "healthy",
//   "response_time_ms": 12,
//   "details": {
//     "connection": "connected",
//     "driver": "pgx",
//     "acquired_conns": 5,
//     "idle_conns": 10
//   }
// }
```

#### Redis Check

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

redisHealth := utils.CheckRedis(ctx, redisClient)

// Returns:
// {
//   "status": "healthy",
//   "response_time_ms": 8,
//   "details": {
//     "connection": "connected",
//     "message": "PONG",
//     "response_time_ms": 8
//   }
// }
```

#### Memory Check

```go
memHealth := utils.CheckMemory()

// Returns:
// {
//   "status": "healthy",      // or "warning" if > 90%, "critical" if > 95%
//   "total_mb": 512,
//   "used_mb": 256,
//   "usage_percent": 50
// }
```

#### Disk Check

```go
diskHealth := utils.CheckDisk()

// Returns:
// {
//   "status": "healthy",
//   "accessible": true
// }
```

### Controller (`features/controllers/health.go`)

The `HealthController` provides HTTP endpoints for health checks:

#### Comprehensive Health Check

```go
GET /health

// Returns full system, application, and dependency status
{
  "status": "healthy",
  "timestamp": "2025-01-13T10:30:00Z",
  "system": {...},
  "application": {...},
  "database": {...},
  "redis": {...},
  "memory": {...},
  "disk": {...},
  "cpu": {...},
  "checks": {
    "database": "healthy",
    "redis": "healthy",
    "memory": "healthy",
    "disk": "healthy"
  }
}
```

#### Kubernetes Probes

**Liveness Probe** (checks if app is running):
```go
GET /health/live

// Returns 200 if alive
{
  "status": "alive",
  "time": "2025-01-13T10:30:00Z"
}
```

**Readiness Probe** (checks if app can handle requests):
```go
GET /health/ready

// Returns 200 if ready and all dependencies are healthy
// Returns 503 if not ready
{
  "status": "ready",
  "time": "2025-01-13T10:30:00Z"
}
```

#### Detailed Endpoints

**System Health**:
```go
GET /health/system

// Returns CPU, memory, platform information
```

**Application Health**:
```go
GET /health/app

// Returns uptime, memory usage, version
```

**Memory Health**:
```go
GET /health/memory

// Returns memory usage and status
```

## Integration

### Setup in FX Module

```go
package di

import (
	"go.uber.org/fx"
	"github.com/Harmeet10000/Fortress_API/src/internal/features/controllers"
)

func HealthModule() fx.Option {
	return fx.Provide(
		controllers.NewHealthController,
	)
}
```

### Register Routes

```go
package routes

func RegisterHealthRoutes(e *echo.Echo, hc *controllers.HealthController) {
	healthRoutes.HealthRoutes(e, hc)
}
```

### In main.go

```go
app := fx.New(
	// ... other modules
	di.HealthModule(),
	fx.Invoke(routes.RegisterHealthRoutes),
)
```

## Docker Compose Configuration

```yaml
# docker-compose.yml
services:
  api:
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health/ready"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s
```

## Kubernetes Configuration

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fortress-api
spec:
  template:
    spec:
      containers:
      - name: api
        image: fortress-api:latest
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

## Status Codes

| Endpoint | Status 200 | Status 503 |
|----------|-----------|-----------|
| `/health` | Healthy | Unhealthy |
| `/health/live` | Running | Never returns 503 |
| `/health/ready` | Ready | Not ready / dependencies down |
| `/health/system` | Always 200 | N/A |
| `/health/app` | Always 200 | N/A |
| `/health/memory` | Always 200 | N/A |

## Memory Thresholds

Memory health uses the following thresholds:

- **Healthy**: ≤ 90% heap usage
- **Warning**: 90% - 95% heap usage
- **Critical**: > 95% heap usage

## Configuration

Environment variables used:

```env
# Logging level (affects health check logging)
LEVEL=debug

# Environment identification
ENVIRONMENT=development

# Health check timeouts (5 seconds)
# Can be customized in the code
```

## Testing Health Checks

### Using curl

```bash
# Full health check
curl http://localhost:8080/health | jq

# Liveness probe
curl http://localhost:8080/health/live

# Readiness probe
curl http://localhost:8080/health/ready

# System metrics
curl http://localhost:8080/health/system | jq

# Application metrics
curl http://localhost:8080/health/app | jq

# Memory status
curl http://localhost:8080/health/memory | jq
```

### Monitoring with watch

```bash
# Continuous monitoring every 2 seconds
watch -n 2 'curl -s http://localhost:8080/health | jq ".checks"'
```

## Performance Considerations

- **Timeout**: All external checks use 5-10 second timeouts
- **Database Pool**: Uses pool stats instead of individual queries for efficiency
- **Memory**: Reads from runtime.MemStats (no allocations)
- **Caching**: Results are generated on-demand (no caching)

## Best Practices

1. **Use readiness probes** for load balancers and orchestration
2. **Use liveness probes** for automatic container restarts
3. **Monitor memory trends** to detect leaks early
4. **Set appropriate timeouts** in your orchestration platform
5. **Log health check failures** for debugging

## Troubleshooting

### Redis shows unhealthy
- Check Redis connection string in `.env`
- Verify Redis is running: `redis-cli ping`
- Check firewall rules

### Database shows unhealthy
- Verify PostgreSQL is running
- Check connection string in `.env`
- Run migrations: `make migrate-up`

### Memory usage warning
- Check for goroutine leaks: `GET /health/app`
- Monitor `num_gc` for excessive garbage collection
- Consider increasing memory limits

## Future Enhancements

- [ ] Custom health check registration
- [ ] Metric export (Prometheus format)
- [ ] Health check history/trending
- [ ] Configurable thresholds
- [ ] Alert integrations (PagerDuty, Slack)
