# Health Checker

The Health Checker service periodically checks server health via HTTP GET requests and notifies users of downtime or recovery. It is built with Go and can be run using Docker.

## Features
- Periodic server health checks.
- Logs errors for unreachable servers or non-204 responses.
- Email notifications for failures and recoveries.
- Configurable check intervals and timeouts.

## Prerequisites

- **Docker**: Install [Docker](https://docs.docker.com/get-docker).
- **Golang** (optional for local testing): Requires Go 1.23.

## Installation
```bash
git clone <url>
cd health-checker
```

### Adjust the Environment
```bash
cp -n .env.example .env
vi .env
```

### Install Swag CLI (for Swagger Documentation)
```bash
go install github.com/swaggo/swag/cmd/swag@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

### Run in the Local
```bash
make docker-compose-local-up
```

### Run in the Background
```bash
make docker-compose-up-detached
```

### Run in attached mode to view the logs
```bash
make docker-compose-up
```

### After the project is up and running, execute SQL commands
```bash
make db-migrations-up
```

## Usage

### Accessing the API
Once the container is running, access the API via http://localhost:${APP_PORT}.

### Swagger UI
RestAPI provides Swagger UI for testing API endpoints.

Visit: http://{env.HOST}:{env.APP_PORT}/swagger/index.html

## Final Notes

We hope you find this helpful. Thank you for using our service! 🚀