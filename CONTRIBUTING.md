# Contributing to Bank Server

Thank you for your interest in contributing to the Bank Server project! This document provides guidelines and information for contributors.

## Development Setup

### Prerequisites

Ensure you have the following tools installed:
- [Go 1.23+](https://golang.org/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop)
- [PostgreSQL client tools](https://www.postgresql.org/download/)
- [Make](https://www.gnu.org/software/make/)

### Getting Started

1. Fork the repository and clone your fork:
   ```bash
   git clone https://github.com/umarhadi/bank-server.git
   cd bank-server
   ```

2. Set up the development environment:
   ```bash
   # Create Docker network
   make network
   
   # Start PostgreSQL container
   make postgres
   
   # Create database
   make createdb
   
   # Run migrations
   make migrateup
   ```

3. Install development dependencies:
   ```bash
   go mod download
   ```

## Code Style and Standards

### Go Code Guidelines

- Follow standard Go formatting using `gofmt`
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions small and focused on a single responsibility
- Use proper error handling - don't ignore errors

### Testing Guidelines

- Write unit tests for all new functionality
- Maintain or improve test coverage (aim for >80%)
- Use table-driven tests for multiple test cases
- Mock external dependencies using the existing mock framework
- Test both success and error scenarios

### Database Guidelines

- Use sqlc for database operations
- Write proper SQL migrations with both up and down scripts
- Follow PostgreSQL best practices for schema design
- Add proper indexes for performance

## Project Structure

```
.
├── api/          # REST API handlers (Gin)
├── db/           # Database layer
│   ├── migration/# SQL migration files
│   ├── query/    # SQL queries for sqlc
│   └── sqlc/     # Generated Go code from sqlc
├── doc/          # Documentation and generated API docs
├── gapi/         # gRPC API handlers
├── mail/         # Email service
├── pb/           # Generated protobuf code
├── proto/        # Protocol buffer definitions
├── token/        # JWT/PASETO token management
├── util/         # Utility functions
├── val/          # Input validation
└── worker/       # Background job processing
```

## Development Workflow

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./...

# Run only short tests (no database required)
go test -short ./...
```

### Code Generation

```bash
# Generate SQL code with sqlc
make sqlc

# Generate mocks
make mock

# Generate protobuf code
make proto
```

### Database Operations

```bash
# Create new migration
make new_migration name=your_migration_name

# Run migrations
make migrateup

# Rollback migrations
make migratedown1
```

## Submitting Changes

### Pull Request Process

1. Create a feature branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the code guidelines

3. Add or update tests for your changes

4. Ensure all tests pass:
   ```bash
   make test
   ```

5. Commit your changes with a clear commit message:
   ```bash
   git commit -m "feat: add user account suspension feature"
   ```

6. Push to your fork and create a pull request

### Commit Message Format

Use conventional commit format:
- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `test:` for adding tests
- `refactor:` for code refactoring
- `chore:` for maintenance tasks

## API Documentation

### REST API

The REST API is documented using Swagger/OpenAPI. Access the interactive documentation at:
- Local: `http://localhost:8080/swagger`
- Production: `https://bank.api.umarhadi.dev/swagger`

### gRPC API

gRPC API documentation is generated from protobuf files. Use Evans CLI for testing:
```bash
make evans
```

## Getting Help

- Open an issue for bug reports or feature requests
- Check existing issues before creating new ones
- Join discussions in the issue comments
- Contact maintainers for questions about contribution guidelines

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help create a welcoming environment for all contributors
- Follow GitHub's community guidelines

Thank you for contributing to Bank Server!