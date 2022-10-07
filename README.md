# Bank Server

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