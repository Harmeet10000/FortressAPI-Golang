# Lint entire project
golangci-lint run

# Lint specific directory
golangci-lint run ./src/...

# Lint specific file
golangci-lint run ./src/cmd/api/main.go

# Show all issues (including warnings)
golangci-lint run -v

# Fix auto-fixable issues (not all linters support this)
golangci-lint run --fix
