# Health Checker

The Health Checker service is designed to periodically check the health of multiple servers by sending HTTP GET requests and sending notifications if the servers are down or if they recover. The service is implemented in Go and can be easily run using Docker.

## Features
- Periodically checks the health of servers via HTTP GET requests.
- Logs errors when a server is unreachable or returns a non-200 HTTP status.
- Sends email notifications on server errors and recoveries.
- Configurable intervals and timeout for health checks.
- Exposes a /status endpoint to check the status of all monitored servers.

## Prerequisites

- **Docker** Make sure Docker installed. Refer to [Docker Installation Guide](https://docs.docker.com/get-docker) if needed.
- **Golang** (optional for local testing): Requires Go 1.23

## Installation
```bash
git clone <url>
cd health_checker
```

### Install Swag CLI (for Swagger Documentation)
```bash
go install github.com/swaggo/swag/cmd/swag@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

### Run in attached mode to view the logs
```bash
make docker-compose-up
```

### After the project is up and running, execute SQL commands
```bash
make pg-migrations-up
```
## Usage

### Accessing the API
Once the container is running, access the API via http://localhost:${APP_PORT}.

### Swagger UI
RestAPI provides Swagger UI for testing API endpoints.

Visit: http://${HOST}:${APP_PORT}/swagger/index.html

## Final Notes

We hope you find this helpful. Thank you for using our service! 🚀