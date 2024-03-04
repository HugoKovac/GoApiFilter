# Domain Status Service Documentation

## Overview
The Domain Status Service is a RESTful API designed to handle domain status checks. It provides two endpoints:

1. **Submit Domain Endpoint**: Accepts POST requests to submit a domain for status checking and stores the status in a Redis database.
2. **Domain Status Endpoint**: Accepts GET requests to retrieve the status of a domain from the Redis database.

## Install
You should have docker and docker compose installed.

Then just run: 

```sh
docker compose up
```

## Structure

In `docker-compose.yaml` you will see that there is 3 containers:

- `redis`: Running a Redis server to store data. Volume pointing to `./redis-data` for persistence
- `server`: Running the Go HTTP **Server**. The source code is into a volume pointing to `./server`
- `client`: Running the Go HTTP **Client**. The source code is into a volume pointing to `./client`

## Endpoints

### 1. Submit Domain Endpoint
- **Path**: `/v1/submit_domain`
- **Method**: POST

#### Request Body
The endpoint expects a JSON object containing the domain name.

```json
{
    "domain": "example.com"
}
```

#### Response
- **201 Created**: If the domain is successfully submitted and added to the database.
- **400 Bad Request**: If the request body is invalid.

### 2. Domain Status Endpoint
- **Path**: `/v1/domain_status`
- **Method**: GET

#### Query Parameters
- **domain**: The domain name for which status is requested.

#### Response
- The response contains the status of the domain.
  - Possible values: "blocked", "allowed", or "unknown" (if the domain is not found in the database).

## Dependencies
- **Redis**: The service utilizes Redis as a backend for storing domain status data.

## Concurrency
- The service employs concurrency to avoid blocking API calls while checking the status of domains.

## Code Structure
- The code is structured into multiple files:
  - `main.go`: Contains the main server logic and HTTP request handlers.
  - `domain.go`: Defines the `Domain` struct for JSON parsing.
  - External dependencies are imported via `import` statements.

## Redis Configuration
- The Redis client is configured to connect to Redis server running at `redis:6379` with no authentication.

## Running the Service
- The service can be started by running the `main` function in `main.go`.
- By default, the service listens on port `80`.

## Dependencies
- The service relies on the following external libraries:
  - `github.com/redis/go-redis/v9` for Redis client.

## Error Handling
- The service handles errors such as invalid requests, database errors, and unknown domains gracefully, returning appropriate HTTP status codes and error messages.
