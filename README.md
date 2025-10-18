# Bank Server

A robust, production-ready banking server built with Go, providing comprehensive APIs for banking operations including account management, transaction processing, and secure money transfers.

## Features

- **Account Management**: Create and manage bank accounts with owner information, balance tracking, and multi-currency support
- **Transaction Recording**: Automatic recording of all balance changes with complete audit trail
- **Money Transfers**: ACID-compliant money transfers between accounts with proper transaction handling
- **User Authentication**: Secure JWT/PASETO-based authentication with role-based access control
- **Email Verification**: Async email verification system with background job processing
- **gRPC + REST APIs**: Dual API support with automatic documentation generation
- **Database Migrations**: Version-controlled database schema management
- **Comprehensive Testing**: High test coverage with unit and integration tests

## Architecture

The application follows a clean architecture pattern with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐
│   REST API      │    │   gRPC API      │
│   (Gin)         │    │   (Protocol     │
└─────────────────┘    │    Buffers)     │
         │              └─────────────────┘
         │                       │
         └───────────┬───────────┘
                     │
            ┌─────────────────┐
            │  Business Logic │
            │   (Handlers)    │
            └─────────────────┘
                     │
            ┌─────────────────┐
            │   Data Layer    │
            │   (sqlc + TX)   │
            └─────────────────┘
                     │
            ┌─────────────────┐
            │   PostgreSQL    │
            │   Database      │
            └─────────────────┘
```

### Core Components

- **API Layer**: REST (Gin) and gRPC handlers for client communication
- **Authentication**: JWT/PASETO token-based auth with middleware
- **Business Logic**: Transaction processing, validation, and business rules
- **Data Access**: Type-safe SQL operations using sqlc
- **Background Jobs**: Redis-based async task processing with Asynq
- **Email Service**: SMTP-based email notifications

## API Documentation

### REST API Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/users` | Create new user account | No |
| POST | `/users/login` | User authentication | No |
| GET | `/accounts` | List user accounts | Yes |
| POST | `/accounts` | Create new account | Yes |
| GET | `/accounts/:id` | Get specific account | Yes |
| POST | `/transfers` | Transfer money between accounts | Yes |

### gRPC Services

- **UserService**: User management and authentication
- **AccountService**: Account creation and retrieval
- **TransferService**: Money transfer operations

For complete API documentation, visit the Swagger UI at `/swagger` endpoint.

This service will provide APIs for the frontend to do following things:

1. Create and manage bank accounts, which are composed of owner’s name, balance, and currency.
2. Record all balance changes to each of the account. So every time some money is added to or subtracted from the account, an account entry record will be created.
3. Perform a money transfer between 2 accounts. This should happen within a transaction, so that either both accounts’ balance are updated successfully or none of them are.

---
## Demo 
This project already deployed, you can access Swagger at [https://bank.api.umarhadi.dev/swagger](https://bank.api.umarhadi.dev/swagger).

## Setup local development

### Install tools

- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [TablePlus](https://tableplus.com/)
- [Golang](https://golang.org/)
- [Homebrew](https://brew.sh/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

    ```bash
    brew install golang-migrate
    ```

- [DB Docs](https://dbdocs.io/docs)

    ```bash
    npm install -g dbdocs
    dbdocs login
    ```

- [DBML CLI](https://www.dbml.org/cli/#installation)

    ```bash
    npm install -g @dbml/cli
    dbml2sql --version
    ```

- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

    ```bash
    brew install sqlc
    ```

- [Gomock](https://github.com/golang/mock)

    ``` bash
    go install github.com/golang/mock/mockgen@v1.6.0
    ```

### Setup infrastructure

- Create the bank-network

    ``` bash
    make network
    ```

- Start postgres container:

    ```bash
    make postgres
    ```

- Create simple_bank database:

    ```bash
    make createdb
    ```

- Run db migration up all versions:

    ```bash
    make migrateup
    ```

- Run db migration up 1 version:

    ```bash
    make migrateup1
    ```

- Run db migration down all versions:

    ```bash
    make migratedown
    ```

- Run db migration down 1 version:

    ```bash
    make migratedown1
    ```

### Documentation

- Generate DB documentation:

    ```bash
    make db_docs
    ```

- Access the DB documentation at this [link](https://dbdocs.io/umarhadi/bank_server).
### Generate code

- Generate schema SQL file with DBML:

    ```bash
    make db_schema
    ```

- Generate SQL CRUD with sqlc:

    ```bash
    make sqlc
    ```

- Generate DB mock with gomock:

    ```bash
    make mock
    ```

- Create a new db migration:

    ```bash
    migrate create -ext sql -dir db/migration -seq <migration_name>
    ```

### Run

- Run server:

    ```bash
    make server
    ```

- Run test:

    ```bash
    make test
    ```

## Deploy to kubernetes cluster

- [Install nginx ingress controller](https://kubernetes.github.io/ingress-nginx/deploy/#aws):

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.48.1/deploy/static/provider/aws/deploy.yaml
    ```

- [Install cert-manager](https://cert-manager.io/docs/installation/kubernetes/):

    ```bash
    kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.4.0/cert-manager.yaml
    ```

## Deploy to Fly.io

Check out the getting started guide [here](https://fly.io/docs/getting-started/) for installing the Fly CLI. Once you have the CLI installed, you can deploy the app with the following command:

- Clone the repo:

    ```bash
    git clone https://github.com/umarhadi/bank-server.git && cd bank-server
    ```

- Create postgres database:

    ```bash
    flyctl postgres create
    ```
    After that, you will get the database url, save it for later.

- Set environment variable

    - `DB_SOURCE`
        ```bash
        flyctl secrets set DB_SOURCE="postgres://postgres:db_password@db_host:5432/postgres?sslmode=disable"
        ```
        Replace `db_password` and `db_host` with the actual password and host from the database url. Don't delete `sslmode=disable` part, if you delete it, the app will not be able to connect to the database because connection between the app and the database is not encrypted via fly.io internal network.

    - `TOKEN_SYMMETRIC_KEY`
        ```bash
        flyctl secrets set TOKEN_SYMMETRIC_KEY="your_token_symmetric_key"
        ```
        Replace `your_token_symmetric_key` with your own 32 symmetric key. You can generate it with this [tool](https://www.browserling.com/tools/random-hex) or 
        ```bash
        openssl rand -hex 64 | head -c 32
        ```
    
    - Another environment variable
    
        You can change the default value of the environment variable in the `fly.toml` file.
- Deploy the app:

    ```bash
    flyctl deploy
    ```

## Testing Strategy

### Test Coverage Overview

The project maintains high test coverage across multiple layers:

- **API Layer (100%)**: Complete coverage of REST endpoints with success and error scenarios
- **Utilities (100%)**: Full coverage of helper functions and validation logic
- **Validation (100%)**: Complete testing of input validation rules
- **gRPC Layer (44.3%)**: Partial coverage, focusing on core user operations
- **Token Management (79.6%)**: Good coverage of JWT/PASETO token operations

### Running Tests

```bash
# Run all tests with coverage
make test

# Run tests without external dependencies
go test -short ./...

# Generate detailed coverage report
go test -v -cover -coverprofile=coverage.out ./...
```

### Test Categories

1. **Unit Tests**: Individual function and method testing
2. **Integration Tests**: Database and API endpoint testing
3. **Mock Tests**: External dependency simulation using generated mocks

For continuous integration, tests are automatically run on pull requests with coverage reporting.