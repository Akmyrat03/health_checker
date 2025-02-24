w# Health Checker 

The Health Checker service is designed to periodically check the health of multiple servers by sending HTTP GET requests and sending notifications if the servers are down or if they recover. The service is implemented in Go and can be easily run using Docker.

## Features
- Periodically checks the health of servers via HTTP GET requests.
- Logs errors when a server is unreachable or returns a non-200 HTTP status.
- Sends email notifications on server errors and recoveries.
- Configurable intervals and timeout for health checks.
- Exposes a /status endpoint to check the status of all monitored servers.

## Prerequisites

- **Docker** Make sure Docker installed. Refer to [Docker Installation Guide](https://docs.docker.com/get-docker) if needed.
- **Golang** (optional for local testing): Requires Go 1.20

## Installation
```bash
git clone <url>
cd health_checker
```