# 🏰 Fortress API

> A modern, scalable, and production-ready backend API built with Go. Designed for enterprise-grade applications with a focus on performance, reliability, and developer experience.

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.25.5-blue?logo=go&logoColor=white)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![Code Quality](https://img.shields.io/badge/Code%20Quality-A-brightgreen)](./src)

**[Features](#-features)** • **[Quick Start](#-quick-start)** • **[Architecture](#-architecture)** • **[Documentation](#-documentation)**

</div>

---

## 🌟 Features

<table>
<tr>
<td width="50%">

### Core
- ⚡ **High-Performance** HTTP server with Echo framework
- 🔒 **Enterprise-Grade Security** with JWT authentication
- 📦 **Type-Safe Database Layer** using SQLc + PostgreSQL
- 🔄 **Database Migrations** with Goose
- ✅ **Comprehensive Validation** with go-playground/validator

</td>
<td width="50%">

### Advanced
- 📊 **Structured Logging** with Uber Zap
- 🗃️ **Redis Caching** for high-performance data retrieval
- 📨 **Message Queue** with RabbitMQ for async processing
- 📡 **gRPC Support** with Protocol Buffers
- 🧩 **Dependency Injection** with Uber FX

</td>
</tr>
<tr>
<td colspan="2">

### Observability & Operations
- 📈 **New Relic APM** for performance monitoring
- 🏥 **Health Checks** with configurable intervals
- 🐳 **Docker Support** for dev and production environments
- 🚀 **Infrastructure as Code** templates for AWS, Azure, and GCP
- 🔧 **Environment Configuration** with Viper

</td>
</tr>
</table>

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.25.5** or higher
- **Docker & Docker Compose** (optional, for containerized development)
- **PostgreSQL 17+** for database
- **Redis 7+** for caching
- **RabbitMQ 3.12+** for message queue (optional)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Harmeet10000/Fortress_API.git
   cd Fortress_API
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env.dev
   # Edit .env.dev with your configuration
   ```

3. **Start services with Docker Compose**
   ```bash
   docker-compose up -d
   ```

4. **Run database migrations**
   ```bash
   make migrate-up
   ```

5. **Start the API server**
   ```bash
   make dev
   ```

The API will be available at `http://localhost:8080`

---

## 📁 Project Structure

```
Fortress_API/
├── src/
│   ├── cmd/                    # Entry points
│   │   ├── api/               # Main API server
│   │   ├── migrate/           # Migration utilities
│   │   └── workers/           # Background workers
│   └── internal/              # Private application code
│       ├── config/            # Configuration management
│       ├── connections/       # Database & service connections
│       ├── db/                # Database layer
│       │   ├── migrations/    # Goose migrations
│       │   ├── schema/        # SQLc generated code
│       │   └── seeders/       # Database seeders
│       ├── features/          # Feature modules
│       │   ├── controllers/   # HTTP handlers
│       │   ├── services/      # Business logic
│       │   ├── repository/    # Data access layer
│       │   ├── models/        # Domain models
│       │   ├── routes/        # Route definitions
│       │   ├── validations/   # Input validation
│       │   └── controllers/   # HTTP handlers
│       ├── middlewares/       # HTTP middleware (auth, logging, etc.)
│       ├── helpers/           # Utility functions
│       └── utils/             # Common utilities
├── pkg/                        # Public packages
├── tests/                      # Test suites
│   ├── unit/                  # Unit tests
│   ├── integration/           # Integration tests
│   ├── e2e/                   # End-to-end tests
│   └── performance/           # Performance benchmarks
├── infra/                     # Infrastructure as Code
│   ├── aws/                   # AWS CloudFormation/Terraform
│   ├── azure/                 # Azure Resource Manager templates
│   └── gcp/                   # Google Cloud Deployment Manager
├── docker/                    # Docker configuration
├── scripts/                   # Utility scripts
├── certs/                     # SSL/TLS certificates
└── docs/                      # API documentation
```

---

## 🏗️ Architecture

### Layered Architecture

Fortress API follows a clean, layered architecture pattern for optimal separation of concerns:

```
┌─────────────────────────────────────────┐
│         HTTP Layer (Echo)               │
│    Controllers & Request/Response       │
├─────────────────────────────────────────┤
│         Business Logic Layer            │
│   Services & Domain Logic               │
├─────────────────────────────────────────┤
│         Data Access Layer               │
│   Repositories & SQLc Queries           │
├─────────────────────────────────────────┤
│         Database Layer                  │
│   PostgreSQL + pgx driver               │
└─────────────────────────────────────────┘
```

### Key Design Principles

- **Separation of Concerns**: Each layer has a single, well-defined responsibility
- **Dependency Injection**: All dependencies wired through Uber FX
- **Type Safety**: Leverages Go's type system and SQLc for compile-time guarantees
- **Error Handling**: Comprehensive error wrapping with context
- **Testability**: Interfaces and dependency injection enable easy testing

---

## 🛠️ Technology Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Framework** | [Echo](https://echo.labstack.com/) | Lightweight, high-performance HTTP server |
| **Database** | [PostgreSQL](https://www.postgresql.org/) + [pgx](https://github.com/jackc/pgx) | Reliable relational database |
| **SQL Generation** | [SQLc](https://sqlc.dev/) | Type-safe SQL code generation |
| **Migrations** | [Goose](https://github.com/pressly/goose) | Database schema versioning |
| **Validation** | [go-playground/validator](https://github.com/go-playground/validator) | Struct validation |
| **Config** | [Viper](https://github.com/spf13/viper) + [Koanf](https://github.com/knadh/koanf) | Flexible configuration management |
| **Logging** | [Zap](https://github.com/uber-go/zap) | Structured logging |
| **Cache** | [Redis](https://redis.io/) | High-performance caching layer |
| **Message Queue** | [RabbitMQ](https://www.rabbitmq.com/) | Asynchronous task processing |
| **gRPC** | [Protocol Buffers](https://developers.google.com/protocol-buffers) | High-performance RPC |
| **DI Container** | [Uber FX](https://pkg.go.dev/go.uber.org/fx) | Dependency injection framework |
| **APM** | [New Relic](https://newrelic.com/) | Application performance monitoring |

---

## 🚀 Development

### Available Commands

```bash
# Development
make dev              # Start development server with hot reload
make build            # Build the application
make run              # Run the built application

# Database
make migrate-up       # Apply all pending migrations
make migrate-down     # Rollback the last migration
make migrate-status   # Show migration status
make seed-db          # Run database seeders

# Testing
make test             # Run all tests
make test-unit        # Run unit tests only
make test-integration # Run integration tests
make test-coverage    # Generate coverage report
make test-watch       # Run tests in watch mode

# Code Quality
make lint             # Run linters
make fmt              # Format code
make vet              # Run go vet

# Docker
make docker-build     # Build Docker images
make docker-up        # Start Docker containers
make docker-down      # Stop Docker containers

# Utilities
make clean            # Clean build artifacts
make help             # Show all available commands
```

### Setting Up Development Environment

```bash
# Install Go dependencies
go mod download
go mod tidy

# Generate code from SQLc
sqlc generate

# Set up pre-commit hooks (optional)
pre-commit install
```

### Environment Configuration

The application loads configuration from multiple sources in this order (highest priority first):

1. **Environment Variables** - Override everything
2. **.env** - Local development overrides
3. **.env.dev** - Default development configuration

Example `.env` file:

```env
ENV=development
PORT=8080
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=fortress_dev
REDIS_ADDRESS=redis://localhost:6379
SECRET_KEY=your-secret-key
LEVEL=debug
```

---

## 🧪 Testing

Fortress API follows comprehensive testing practices:

- **Unit Tests**: Test individual functions and services
- **Integration Tests**: Test interactions between layers
- **E2E Tests**: Test complete workflows
- **Performance Tests**: Benchmark critical operations

```bash
# Run all tests with coverage
make test-coverage

# Run specific test suite
go test ./... -run TestUserService -v

# Run benchmarks
go test -bench=. ./...
```

Target coverage: **80%+** on critical paths

---

## 📖 API Documentation

### Authentication

All protected endpoints require a Bearer token in the Authorization header:

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/protected
```

### Example Endpoints

*Documentation structure - expand with your actual endpoints*

#### Health Check
```bash
GET /health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2025-01-13T10:30:00Z"
}
```

#### Create Resource
```bash
POST /api/v1/resources
Content-Type: application/json

{
  "name": "Resource Name",
  "description": "Resource Description"
}
```

---

## 🔒 Security

- **Authentication**: JWT-based authentication with configurable secret
- **Input Validation**: All inputs validated before processing
- **SQL Injection Prevention**: Type-safe SQLc queries prevent SQL injection
- **CORS**: Configurable CORS policy for frontend integration
- **SSL/TLS**: Support for HTTPS with certificate management
- **Secrets Management**: Environment-based secret handling

### Security Best Practices

1. Never commit `.env` files with production secrets
2. Use strong `SECRET_KEY` values in production
3. Enable HTTPS in production
4. Regularly update dependencies: `go get -u ./...`
5. Review security advisories: `go vulnerabilities ./...`

---

## 🐳 Docker

### Development with Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop all services
docker-compose down
```

### Building Production Images

```bash
# Build using prod Dockerfile
docker build -f docker/prod.Dockerfile -t fortress-api:latest .

# Run container
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e DATABASE_HOST=db \
  fortress-api:latest
```

---

## 📊 Monitoring & Observability

### Logging

Structured logging with Zap provides detailed insights:

```go
logger.Info("user created",
    zap.String("user_id", user.ID),
    zap.String("email", user.Email),
)
```

### Health Checks

Configure health checks in your environment:

```env
ENABLED=true
INTERVAL=30s
TIMEOUT=5s
CHECKS=database,redis
```

### New Relic Integration

Enable APM monitoring:

```env
LICENSE_KEY=your-new-relic-license-key
APP_LOG_FORWARDING_ENABLED=true
DISTRIBUTED_TRACING_ENABLED=true
```

---

## 🚀 Deployment

### Infrastructure Templates

Ready-to-use IaC templates for major cloud providers:

- **AWS**: CloudFormation & Terraform templates in `infra/aws/`
- **Azure**: ARM templates in `infra/azure/`
- **GCP**: Deployment Manager templates in `infra/gcp/`

### Environment-Specific Configuration

Configure different environments:

```bash
# Development
ENV=development LEVEL=debug make dev

# Staging
ENV=staging LEVEL=info make build && docker run ...

# Production
ENV=production LEVEL=warn make build && docker run ...
```

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Workflow

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

### Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go) conventions
- Run `make fmt` before committing
- Ensure `make lint` passes
- Write tests for new features

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🛡️ Security Policy

Please report security vulnerabilities to [harmeetsinghfbd@gmail.com](mailto:security@example.com) rather than using the issue tracker. See [SECURITY.md](SECURITY.md) for more details.

---

## 🙋 Support

- **Issues**: [GitHub Issues](https://github.com/Harmeet10000/Fortress_API/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Harmeet10000/Fortress_API/discussions)
- **Email**: [your-email@example.com](mailto:your-email@example.com)

---

<div align="center">

**[⬆ back to top](#-fortress-api)**

Made with ❤️ by [Harmeet](https://github.com/Harmeet10000)

</div>
