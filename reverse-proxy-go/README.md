# Reverse Proxy (Go)

A high-performance reverse proxy server built with Go that routes incoming requests to deployed projects on Cloudflare R2 storage.

## Features

- **Subdomain-based routing**: Routes requests based on subdomain to the corresponding project
- **Custom domain support**: Supports custom domains for projects
- **PostgreSQL integration**: Queries project and deployment information from PostgreSQL
- **Cloudflare R2 integration**: Proxies static assets from Cloudflare R2 storage
- **Graceful shutdown**: Handles SIGINT and SIGTERM signals for clean shutdowns
- **Hot reload**: Development mode with Air for instant reloads

## Architecture

The reverse proxy follows a clean architecture pattern:

```
reverse-proxy-go/
├── main.go                          # Entry point
├── internal/
│   ├── app/                         # Application initialization
│   ├── config/                      # Configuration management
│   ├── db/                          # Database connection
│   ├── domain/                      # Business models
│   │   ├── project/
│   │   └── deployment/
│   ├── handler/                     # HTTP handlers
│   │   └── proxy/
│   ├── middleware/                  # HTTP middleware
│   ├── repository/                  # Data access layer
│   │   └── project/
│   └── router/                      # Route registration
├── go.mod                           # Go modules
├── .env                             # Environment variables
├── .air.toml                        # Air configuration
└── Dockerfile                       # Docker configuration
```

## Prerequisites

- Go 1.23 or higher
- PostgreSQL database
- Cloudflare R2 bucket with public access

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Create a `.env` file:
```env
DATABASE_URL=postgresql://user:password@host:port/database
R2_PUBLIC_URL=https://pub-xxx.r2.dev
PORT=8001
```

3. Run the server:
```bash
# Development mode with hot reload
air

# Production mode
go run main.go
```

## Docker

Build and run with Docker:

```bash
# Build image
docker build -t reverse-proxy-go .

# Run container
docker run -p 8001:8001 --env-file .env reverse-proxy-go
```

## How It Works

1. Incoming request arrives at the proxy (e.g., `myapp.localhost:8001`)
2. Proxy extracts the subdomain (`myapp`)
3. Queries PostgreSQL for project with matching subdomain
4. Retrieves the latest READY deployment
5. Proxies the request to Cloudflare R2: `https://pub-xxx.r2.dev/{project-id}/{path}`
6. Returns the response to the client

## Database Schema

The proxy expects the following PostgreSQL schema:

```sql
-- Projects table
CREATE TABLE projects (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    git_url VARCHAR NOT NULL,
    subdomain VARCHAR UNIQUE NOT NULL,
    custom_domain VARCHAR,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Deployments table
CREATE TABLE deployments (
    id UUID PRIMARY KEY,
    project_id UUID REFERENCES projects(id),
    status VARCHAR CHECK (status IN ('NOT_STARTED', 'QUEUED', 'IN_PROGRESS', 'READY', 'FAIL')),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## Dependencies

- **chi/v5**: Lightweight HTTP router
- **chi/cors**: CORS middleware
- **godotenv**: Environment variable loading
- **lib/pq**: PostgreSQL driver

## Development

```bash
# Install Air for hot reload
go install github.com/air-verse/air@latest

# Run with hot reload
air

# Build binary
go build -o reverse-proxy .

# Run tests
go test ./...
```

## License

MIT
