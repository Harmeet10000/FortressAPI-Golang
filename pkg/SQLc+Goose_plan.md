# Complete Guide: sqlc + Goose - The Perfect Combo

## 📚 **Table of Contents**

1. [Introduction & Setup](#introduction)
2. [Project Structure](#project-structure)
3. [Installation](#installation)
4. [Database Schema Design](#schema-design)
5. [Migrations with Goose](#migrations)
6. [Queries with sqlc](#queries)
7. [Complete CRUD Operations](#crud-operations)
8. [Advanced Patterns](#advanced-patterns)
9. [Testing](#testing)
10. [Production Deployment](#deployment)
11. [Best Practices](#best-practices)
12. [Troubleshooting](#troubleshooting)

---

<a name="introduction"></a>
## 📖 **1. Introduction**

### **Why sqlc + Goose?**

This combination is the **gold standard** for Go database work:

- **Goose**: Handles database migrations (schema changes)
- **sqlc**: Generates type-safe Go code from SQL queries

```
┌─────────────────────────────────────────┐
│           Your Application              │
├─────────────────────────────────────────┤
│                                         │
│  ┌──────────┐         ┌──────────┐     │
│  │  Goose   │────────▶│ Database │     │
│  │(Migrate) │         │ (Schema) │     │
│  └──────────┘         └────┬─────┘     │
│                            │            │
│  ┌──────────┐              │            │
│  │  sqlc    │──────────────┘            │
│  │(Queries) │                           │
│  └──────────┘                           │
│                                         │
└─────────────────────────────────────────┘
```

**Perfect for:**
- Production applications
- Team projects
- Long-term maintainability
- Type safety obsessives
- SQL lovers

---

<a name="project-structure"></a>
## 🏗️ **2. Project Structure**

### **Recommended Structure**

```
my-app/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── db/
│   ├── migrations/              # Goose migrations
│   │   ├── 00001_create_users_table.sql
│   │   ├── 00002_create_posts_table.sql
│   │   ├── 00003_create_comments_table.sql
│   │   └── 00004_add_indexes.sql
│   ├── queries/                 # sqlc queries
│   │   ├── users.sql
│   │   ├── posts.sql
│   │   └── comments.sql
│   └── schema.sql               # Complete schema (for reference)
├── internal/
│   ├── database/                # Generated sqlc code
│   │   ├── db.go
│   │   ├── models.go
│   │   ├── querier.go
│   │   ├── users.sql.go
│   │   ├── posts.sql.go
│   │   └── comments.sql.go
│   ├── repository/              # Business logic layer
│   │   ├── user_repository.go
│   │   └── post_repository.go
│   ├── service/                 # Service layer
│   │   ├── user_service.go
│   │   └── post_service.go
│   └── handler/                 # HTTP handlers
│       ├── user_handler.go
│       └── post_handler.go
├── pkg/
│   └── config/
│       └── config.go            # Configuration
├── scripts/
│   ├── migrate.sh               # Migration helper scripts
│   └── seed.sh                  # Seed data
├── .env.example                 # Environment variables template
├── .gitignore
├── docker-compose.yml           # Local development
├── Dockerfile                   # Production container
├── go.mod
├── go.sum
├── Makefile                     # Development commands
├── README.md
└── sqlc.yaml                    # sqlc configuration
```

---

<a name="installation"></a>
## 🚀 **3. Installation**

### **Step 1: Install Tools**

```bash
# Install Goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Verify installation
goose version
sqlc version
```

### **Step 2: Initialize Project**

```bash
# Create project
mkdir my-blog-api
cd my-blog-api

# Initialize Go module
go mod init github.com/yourusername/my-blog-api

# Install database driver
go get github.com/lib/pq  # PostgreSQL

# Create directory structure
mkdir -p cmd/server
mkdir -p db/{migrations,queries}
mkdir -p internal/{database,repository,service,handler}
mkdir -p pkg/config
mkdir -p scripts
```

### **Step 3: Create Configuration Files**

**sqlc.yaml**
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "db/queries"
    schema: "db/migrations"
    gen:
      go:
        package: "database"
        out: "internal/database"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_exact_table_names: false
        emit_empty_slices: true
        emit_exported_queries: true
        emit_result_struct_pointers: true
        emit_params_struct_pointers: false
        emit_methods_with_db_argument: false
        emit_pointers_for_null_types: true
        emit_enum_valid_method: true
        emit_all_enum_values: true
        json_tags_case_style: "snake"
```

**.env.example**
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=myblog
DB_SSLMODE=disable

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# Environment
ENV=development
```

**.gitignore**
```
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
bin/
dist/

# Test binary
*.test

# Output of the go coverage tool
*.out

# Dependencies
vendor/

# Go workspace file
go.work

# Environment variables
.env

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Database
*.db
*.sqlite
*.sqlite3
```

**Makefile**
```makefile
# Variables
DB_URL := postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATIONS_DIR := db/migrations

# Load environment variables
include .env
export

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: install-tools
install-tools: ## Install required tools
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: setup
setup: ## Initial project setup
	cp .env.example .env
	@echo "Please edit .env file with your settings"

# Database commands
.PHONY: db-up
db-up: ## Start database (Docker)
	docker-compose up -d postgres

.PHONY: db-down
db-down: ## Stop database (Docker)
	docker-compose down

.PHONY: db-create
db-create: ## Create database
	createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)

.PHONY: db-drop
db-drop: ## Drop database
	dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)

# Migration commands
.PHONY: migrate-up
migrate-up: ## Run all pending migrations
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

.PHONY: migrate-down
migrate-down: ## Rollback last migration
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

.PHONY: migrate-status
migrate-status: ## Show migration status
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

.PHONY: migrate-reset
migrate-reset: ## Rollback all migrations
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" reset

.PHONY: migrate-redo
migrate-redo: ## Rollback and reapply last migration
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" redo

.PHONY: migrate-create
migrate-create: ## Create new migration (usage: make migrate-create name=create_users)
	@read -p "Enter migration name: " name; \
	goose -dir $(MIGRATIONS_DIR) create $$name sql

# sqlc commands
.PHONY: sqlc-generate
sqlc-generate: ## Generate Go code from SQL
	sqlc generate

.PHONY: sqlc-vet
sqlc-vet: ## Check SQL queries
	sqlc vet

# Development commands
.PHONY: run
run: ## Run the application
	go run cmd/server/main.go

.PHONY: build
build: ## Build the application
	go build -o bin/server cmd/server/main.go

.PHONY: test
test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

.PHONY: test-coverage
test-coverage: test ## Run tests with coverage report
	go tool cover -html=coverage.out

.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out

# Docker commands
.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t my-blog-api:latest .

.PHONY: docker-up
docker-up: ## Start all services with Docker Compose
	docker-compose up -d

.PHONY: docker-down
docker-down: ## Stop all services
	docker-compose down

.PHONY: docker-logs
docker-logs: ## Show Docker logs
	docker-compose logs -f

# Development workflow
.PHONY: dev-setup
dev-setup: install-tools setup db-up migrate-up sqlc-generate ## Complete development setup

.PHONY: dev-reset
dev-reset: migrate-reset migrate-up sqlc-generate ## Reset database and regenerate code

.PHONY: new-migration
new-migration: migrate-create migrate-up sqlc-generate ## Create and apply new migration
```

---

<a name="schema-design"></a>
## 🗄️ **4. Database Schema Design**

### **Planning Your Schema**

Before writing migrations, plan your database schema:

**Example: Blog Application**

```
┌─────────────┐
│    users    │
├─────────────┤
│ id          │──┐
│ username    │  │
│ email       │  │
│ password    │  │
│ created_at  │  │
│ updated_at  │  │
└─────────────┘  │
                 │
                 │  ┌─────────────┐
                 └─▶│    posts    │
                    ├─────────────┤
                    │ id          │──┐
                    │ user_id     │  │
                    │ title       │  │
                    │ content     │  │
                    │ published   │  │
                    │ created_at  │  │
                    │ updated_at  │  │
                    └─────────────┘  │
                                     │
                                     │  ┌─────────────┐
                                     └─▶│  comments   │
                                        ├─────────────┤
                                        │ id          │
                                        │ post_id     │
                                        │ user_id     │
                                        │ content     │
                                        │ created_at  │
                                        └─────────────┘
```

**db/schema.sql** (Reference only, not used by tools)
```sql
-- This file is for documentation purposes
-- Actual schema is created via migrations

-- Users table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    bio TEXT,
    avatar_url VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Posts table
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    excerpt TEXT,
    published BOOLEAN NOT NULL DEFAULT false,
    published_at TIMESTAMP,
    view_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Comments table
CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id BIGINT REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Tags table
CREATE TABLE tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    slug VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Post tags junction table
CREATE TABLE post_tags (
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    tag_id BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

-- Indexes for performance
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_published ON posts(published);
CREATE INDEX idx_posts_slug ON posts(slug);
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_post_tags_post_id ON post_tags(post_id);
CREATE INDEX idx_post_tags_tag_id ON post_tags(tag_id);
```

---

<a name="migrations"></a>
## 🔄 **5. Migrations with Goose**

### **Creating Migrations**

Migrations should be **small, focused, and reversible**.

#### **Migration 1: Create Users Table**

```bash
make migrate-create name=create_users_table
# Creates: db/migrations/00001_create_users_table.sql
```

**db/migrations/00001_create_users_table.sql**
```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    bio TEXT,
    avatar_url VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- Comments for documentation
COMMENT ON TABLE users IS 'Stores user account information';
COMMENT ON COLUMN users.password_hash IS 'Bcrypt hashed password';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users CASCADE;
-- +goose StatementEnd
```

#### **Migration 2: Create Posts Table**

```bash
make migrate-create name=create_posts_table
```

**db/migrations/00002_create_posts_table.sql**
```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    excerpt TEXT,
    published BOOLEAN NOT NULL DEFAULT false,
    published_at TIMESTAMP,
    view_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Foreign key constraint
    CONSTRAINT fk_posts_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_published ON posts(published);
CREATE INDEX idx_posts_slug ON posts(slug);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);

-- Full-text search index
CREATE INDEX idx_posts_title_search ON posts USING gin(to_tsvector('english', title));
CREATE INDEX idx_posts_content_search ON posts USING gin(to_tsvector('english', content));

-- Comments
COMMENT ON TABLE posts IS 'Blog posts';
COMMENT ON COLUMN posts.slug IS 'URL-friendly version of title';
COMMENT ON COLUMN posts.excerpt IS 'Short preview of post content';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts CASCADE;
-- +goose StatementEnd
```

#### **Migration 3: Create Comments Table**

```bash
make migrate-create name=create_comments_table
```

**db/migrations/00003_create_comments_table.sql**
```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    parent_id BIGINT, -- For nested comments
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Foreign key constraints
    CONSTRAINT fk_comments_post
        FOREIGN KEY (post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_comments_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_comments_parent
        FOREIGN KEY (parent_id)
        REFERENCES comments(id)
        ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_comments_created_at ON comments(created_at DESC);

-- Prevent circular references (comment can't be its own parent)
ALTER TABLE comments ADD CONSTRAINT check_not_self_parent 
    CHECK (id != parent_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS comments CASCADE;
-- +goose StatementEnd
```

#### **Migration 4: Create Tags System**

```bash
make migrate-create name=create_tags_system
```

**db/migrations/00004_create_tags_system.sql**
```sql
-- +goose Up
-- +goose StatementBegin
-- Tags table
CREATE TABLE tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    slug VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Post-Tags junction table (many-to-many)
CREATE TABLE post_tags (
    post_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (post_id, tag_id),
    
    CONSTRAINT fk_post_tags_post
        FOREIGN KEY (post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_post_tags_tag
        FOREIGN KEY (tag_id)
        REFERENCES tags(id)
        ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_tags_slug ON tags(slug);
CREATE INDEX idx_post_tags_post_id ON post_tags(post_id);
CREATE INDEX idx_post_tags_tag_id ON post_tags(tag_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS post_tags CASCADE;
DROP TABLE IF EXISTS tags CASCADE;
-- +goose StatementEnd
```

#### **Migration 5: Add Database Functions**

```bash
make migrate-create name=add_update_timestamp_function
```

**db/migrations/00005_add_update_timestamp_function.sql**
```sql
-- +goose Up
-- +goose StatementBegin
-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply trigger to users table
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Apply trigger to posts table
CREATE TRIGGER update_posts_updated_at
    BEFORE UPDATE ON posts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Apply trigger to comments table
CREATE TRIGGER update_comments_updated_at
    BEFORE UPDATE ON comments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_posts_updated_at ON posts;
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;
DROP FUNCTION IF EXISTS update_updated_at_column();
-- +goose StatementEnd
```

### **Running Migrations**

```bash
# Apply all pending migrations
make migrate-up

# Check status
make migrate-status

# Output:
# Applied At                  Migration
# =======================================
# Wed Nov 15 10:30:45 2024 -- 00001_create_users_table.sql
# Wed Nov 15 10:30:45 2024 -- 00002_create_posts_table.sql
# Wed Nov 15 10:30:45 2024 -- 00003_create_comments_table.sql
# Wed Nov 15 10:30:45 2024 -- 00004_create_tags_system.sql
# Wed Nov 15 10:30:45 2024 -- 00005_add_update_timestamp_function.sql

# Rollback last migration
make migrate-down

# Reset all migrations (careful!)
make migrate-reset

# Redo last migration (down + up)
make migrate-redo
```

---

<a name="queries"></a>
## 🔍 **6. Queries with sqlc**

### **Writing SQL Queries**

Now that we have our schema, let's write queries.

#### **User Queries**

**db/queries/users.sql**
```sql
-- name: CreateUser :one
INSERT INTO users (
    username, email, password_hash, bio, avatar_url
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET 
    bio = COALESCE(sqlc.narg(bio), bio),
    avatar_url = COALESCE(sqlc.narg(avatar_url), avatar_url),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET 
    password_hash = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: SearchUsers :many
SELECT * FROM users
WHERE 
    username ILIKE '%' || $1 || '%' OR
    email ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUserWithStats :one
SELECT 
    u.*,
    COUNT(DISTINCT p.id) as post_count,
    COUNT(DISTINCT c.id) as comment_count
FROM users u
LEFT JOIN posts p ON p.user_id = u.id
LEFT JOIN comments c ON c.user_id = u.id
WHERE u.id = $1
GROUP BY u.id;
```

#### **Post Queries**

**db/queries/posts.sql**
```sql
-- name: CreatePost :one
INSERT INTO posts (
    user_id, title, slug, content, excerpt, published, published_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetPostByID :one
SELECT * FROM posts
WHERE id = $1 LIMIT 1;

-- name: GetPostBySlug :one
SELECT * FROM posts
WHERE slug = $1 LIMIT 1;

-- name: GetPostWithAuthor :one
SELECT 
    p.*,
    u.id as author_id,
    u.username as author_username,
    u.avatar_url as author_avatar
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE p.id = $1 LIMIT 1;

-- name: ListPosts :many
SELECT * FROM posts
WHERE published = true
ORDER BY published_at DESC NULLS LAST, created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListPostsByUser :many
SELECT * FROM posts
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListDraftPosts :many
SELECT * FROM posts
WHERE user_id = $1 AND published = false
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePost :one
UPDATE posts
SET 
    title = COALESCE(sqlc.narg(title), title),
    slug = COALESCE(sqlc.narg(slug), slug),
    content = COALESCE(sqlc.narg(content), content),
    excerpt = COALESCE(sqlc.narg(excerpt), excerpt),
    published = COALESCE(sqlc.narg(published), published),
    published_at = COALESCE(sqlc.narg(published_at), published_at),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: PublishPost :one
UPDATE posts
SET 
    published = true,
    published_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UnpublishPost :one
UPDATE posts
SET 
    published = false,
    published_at = NULL,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;

-- name: IncrementPostViewCount :exec
UPDATE posts
SET view_count = view_count + 1
WHERE id = $1;

-- name: CountPosts :one
SELECT COUNT(*) FROM posts
WHERE published = true;

-- name: CountPostsByUser :one
SELECT COUNT(*) FROM posts
WHERE user_id = $1;

-- name: SearchPosts :many
SELECT * FROM posts
WHERE 
    published = true AND (
        to_tsvector('english', title) @@ plainto_tsquery('english', $1) OR
        to_tsvector('english', content) @@ plainto_tsquery('english', $1)
    )
ORDER BY published_at DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: GetPopularPosts :many
SELECT * FROM posts
WHERE published = true
ORDER BY view_count DESC, published_at DESC
LIMIT $1;

-- name: GetRecentPosts :many
SELECT 
    p.*,
    u.username as author_username
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE p.published = true
AND p.published_at > NOW() - INTERVAL '30 days'
ORDER BY p.published_at DESC
LIMIT $1;
```

#### **Comment Queries**

**db/queries/comments.sql**
```sql
-- name: CreateComment :one
INSERT INTO comments (
    post_id, user_id, parent_id, content
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetCommentByID :one
SELECT * FROM comments
WHERE id = $1 LIMIT 1;

-- name: GetCommentWithAuthor :one
SELECT 
    c.*,
    u.username as author_username,
    u.avatar_url as author_avatar
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.id = $1 LIMIT 1;

-- name: ListCommentsByPost :many
SELECT 
    c.*,
    u.username as author_username,
    u.avatar_url as author_avatar
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.post_id = $1
ORDER BY c.created_at ASC;

-- name: ListTopLevelComments :many
SELECT 
    c.*,
    u.username as author_username,
    u.avatar_url as author_avatar
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.post_id = $1 AND c.parent_id IS NULL
ORDER BY c.created_at ASC;

-- name: ListReplies :many
SELECT 
    c.*,
    u.username as author_username,
    u.avatar_url as author_avatar
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.parent_id = $1
ORDER BY c.created_at ASC;

-- name: ListCommentsByUser :many
SELECT 
    c.*,
    p.title as post_title,
    p.slug as post_slug
FROM comments c
JOIN posts p ON c.post_id = p.id
WHERE c.user_id = $1
ORDER BY c.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateComment :one
UPDATE comments
SET 
    content = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1;

-- name: CountCommentsByPost :one
SELECT COUNT(*) FROM comments
WHERE post_id = $1;

-- name: CountCommentsByUser :one
SELECT COUNT(*) FROM comments
WHERE user_id = $1;
```

#### **Tag Queries**

**db/queries/tags.sql**
```sql
-- name: CreateTag :one
INSERT INTO tags (