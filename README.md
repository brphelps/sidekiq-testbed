# sidekiq-testbed

A Go application that uses Sidekiq-compatible workers to schedule and execute jobs with Redis.

## Features

- Schedules jobs to run twice per minute (every 30 seconds)
- Each job sleeps for 10 seconds before completing
- Uses Redis as the backend for job queuing
- Compatible with Sidekiq job format

## Prerequisites

- Go 1.18 or higher (tested with Go 1.24+)
- Redis server running and accessible

## Configuration

The application uses environment variables for Redis connection:

- `REDIS_SERVER`: Redis server address (default: `localhost:6379`)
- `REDIS_PASSWORD`: Redis password (optional, default: empty)
- `REDIS_DB`: Redis database number (default: `0`)

Copy `.env.example` to `.env` and configure as needed:

```bash
cp .env.example .env
```

## Installation

```bash
# Install dependencies
go mod download

# Build the application
go build -o sidekiq-testbed
```

## Running

### With environment variables

```bash
export REDIS_SERVER=localhost:6379
export REDIS_PASSWORD=your_password
export REDIS_DB=0
./sidekiq-testbed
```

### Or build and run directly

```bash
go run main.go
```

## How It Works

1. The application connects to Redis using the provided connection details
2. A scheduler runs in the background, enqueueing jobs every 30 seconds
3. Worker processes pick up jobs from the `sleep_queue`
4. Each job logs its start time, sleeps for 10 seconds, and logs its completion time
5. The process continues until interrupted (Ctrl+C)

## Development

Run the application in development mode:

```bash
go run main.go
```

## Architecture

- **main.go**: Main application entry point
  - Configures Redis connection
  - Registers the `SleepWorker` worker
  - Starts the job scheduler
  - Starts the worker process
  - Handles graceful shutdown

- **SleepWorker**: A simple worker that sleeps for 10 seconds and completes

## Example Output

```
2024/12/16 22:41:00 Starting Sidekiq worker connected to Redis at localhost:6379
2024/12/16 22:41:00 Initial job enqueued
2024/12/16 22:41:00 Starting job at 2024-12-16T22:41:00Z
2024/12/16 22:41:10 Completed job at 2024-12-16T22:41:10Z
2024/12/16 22:41:30 Job enqueued
2024/12/16 22:41:30 Starting job at 2024-12-16T22:41:30Z
2024/12/16 22:41:40 Completed job at 2024-12-16T22:41:40Z
```