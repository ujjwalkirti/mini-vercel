# Mini Vercel

A self-hosted deployment platform that demonstrates a Vercel-like workflow. Deploy your static sites with Git integration, real-time build logs, and automatic subdomain routing.

## Overview

Mini Vercel consists of four main services working together:

| Service | Port | Description |
|---------|------|-------------|
| **api-server** | 9000 | Control plane API that handles projects, deployments, and build orchestration |
| **frontend** | 5173 | React dashboard for managing projects and viewing deployment logs |
| **reverse-proxy** | 8001 | Routes subdomain requests to deployed static assets |
| **build-project** | Container | Clones repos, runs builds, uploads artifacts to storage |

## Architecture

```
                                    +------------------+
                                    |    Frontend      |
                                    |   (React UI)     |
                                    +--------+---------+
                                             |
                                             v
+---------------+    Deploy Request    +------------------+
|    User       | ------------------> |   API Server     |
+---------------+                      |      (Go)        |
                                      +--------+---------+
                                               |
                          +--------------------+--------------------+
                          |                    |                    |
                          v                    v                    v
                   +--------------+     +--------------+      +--------------+
                   |  PostgreSQL  |     |    Kafka     |      |  AWS ECS     |
                   |  (Supabase)  |     |              |      |  (Fargate)   |
                   +--------------+     +------+-------+      +------+-------+
                                               |                     |
                                               |               +-----v------+
                                               |               |   Build    |
                                               |               | Container  |
                                               |               +-----+------+
                                               |                     |
                                               v                     v
                                        +--------------+     +--------------+
                                        | ClickHouse   |     |Cloudflare R2 |
                                        |   (Logs)     |     |  (Storage)   |
                                        +--------------+     +------+-------+
                                                                    |
+---------------+    subdomain.localhost:8001    +------------------v--------+
|  End User     | <---------------------------- |    Reverse Proxy          |
+---------------+                                +---------------------------+
```

## Features

- **Git-based Deployments**: Deploy directly from GitHub repositories
- **Real-time Build Logs**: Stream build output via Kafka to ClickHouse
- **Automatic Subdomains**: Each project gets a unique subdomain
- **Deployment Status Tracking**: Track builds through QUEUED → IN_PROGRESS → READY/FAIL
- **User Authentication**: Supabase powered auth with JWT verification
- **Concurrent Message Processing**: Worker pool handles Kafka messages with 50 concurrent workers
- **Clean Architecture**: Repository pattern with dependency injection for maintainability

## Tech Stack

### Backend
- **Runtime**: Go 1.21+
- **Framework**: Chi router with custom middleware
- **Database**: PostgreSQL (Supabase) with pgx driver
- **Message Queue**: Apache Kafka with IBM Sarama
- **Analytics**: ClickHouse for log storage
- **Container Orchestration**: AWS ECS Fargate
- **Storage**: Cloudflare R2 (S3 compatible)
- **Authentication**: Supabase Auth with JWT verification and JWKS caching

### Frontend
- **Framework**: React 19 with TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **UI Components**: Radix UI
- **Routing**: React Router DOM

## Repository Structure

```
mini-vercel/
├── api-server-go/             # Control plane API (Go)
│   ├── main.go                # Application entry point
│   ├── internal/
│   │   ├── app/               # App initialization and wiring
│   │   ├── auth/              # JWT verification, JWKS caching
│   │   ├── client/            # External service clients (ClickHouse, Supabase)
│   │   ├── config/            # Configuration loaders (AWS, Kafka, ClickHouse)
│   │   ├── db/                # Database connection
│   │   ├── domain/            # Domain models (Project, Deployment, BuildLog)
│   │   ├── handler/           # HTTP handlers (health, project, deployment)
│   │   ├── kafka/             # Kafka consumer with worker pool
│   │   ├── middleware/        # Auth middleware, context management
│   │   ├── repository/        # Data access layer
│   │   ├── router/            # Route registration
│   │   ├── service/           # Business logic (deployment, logs, ECS)
│   │   └── utils/             # Response helpers, UUID generation
│   ├── Dockerfile
│   └── docker-compose.yml     # Local Kafka & ClickHouse
│
├── frontend/                   # React dashboard
│   ├── src/
│   │   ├── pages/             # Page components
│   │   ├── components/        # Reusable UI components
│   │   ├── contexts/          # Auth context
│   │   └── lib/               # API client, Supabase
│   └── vite.config.ts
│
├── reverse-proxy-go/           # Static asset routing (Go)
│   ├── main.go                # Proxy entry point
│   ├── internal/
│   │   ├── app/               # App initialization
│   │   ├── config/            # Configuration
│   │   ├── db/                # Database connection
│   │   ├── domain/            # Domain models
│   │   ├── handler/           # Proxy handler
│   │   ├── middleware/        # Request logging
│   │   ├── repository/        # Data access layer
│   │   └── router/            # Route setup
│   └── Dockerfile
│
├── build-project/              # Build runner container
│   ├── script.js              # Build pipeline
│   ├── kafkaProducer.js       # Log streaming
│   ├── r2Blob.js              # R2 upload
│   └── Dockerfile
│
└── Readme.md
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Node.js 20+ (for frontend and build container)
- Docker and Docker Compose
- AWS Account (for ECS)
- Supabase Project
- Cloudflare R2 Bucket
- Kafka Cluster (Aiven or self hosted)
- ClickHouse Instance

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/mini-vercel.git
cd mini-vercel

# Install Go dependencies
cd api-server-go
go mod download
cd ../reverse-proxy-go
go mod download

# Install frontend dependencies
cd ../frontend
npm install

# Install build container dependencies
cd ../build-project
npm install
```

### Environment Setup

#### API Server (`api-server-go/.env`)

```env
# AWS ECS
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=us-east-1

# ECS Task Configuration
ECS_CLUSTER_NAME=mini-vercel-cluster
ECS_TASK_DEFINITION=build-task
ECS_SUBNETS=subnet-xxx,subnet-yyy
ECS_SECURITY_GROUPS=sg-xxx
ECS_ASSIGN_PUBLIC_IP=ENABLED
ECS_IMAGE_NAME=your-registry/build-project:latest
ECS_LAUNCH_TYPE=FARGATE
ECS_COUNT=1

# Database (Supabase PostgreSQL)
DATABASE_URL=postgresql://user:password@host:5432/database

# Kafka
KAFKA_BROKERS=broker1:9092,broker2:9092
KAFKA_CLIENT_ID=api-server
KAFKA_USERNAME=your_username
KAFKA_PASSWORD=your_password

# Cloudflare R2
R2_ACCOUNT_ID=your_account_id
R2_ACCESS_KEY_ID=your_r2_access_key
R2_SECRET_ACCESS_KEY=your_r2_secret
R2_BUCKET_NAME=mini-vercel-builds
R2_PUBLIC_URL=https://pub-xxx.r2.dev

# ClickHouse
CLICKHOUSE_HOST=localhost
CLICKHOUSE_PORT=9000
CLICKHOUSE_DATABASE=logs
CLICKHOUSE_USERNAME=default
CLICKHOUSE_PASSWORD=your_password

# Server
PORT=9000
FRONTEND_URL=http://localhost:5173

# Supabase Auth
SUPABASE_URL=https://xxx.supabase.co
SUPABASE_SERVICE_ROLE_KEY=your_service_role_key
```

#### Frontend (`frontend/.env`)

```env
VITE_API_URL=http://localhost:9000
VITE_SUPABASE_URL=
VITE_SUPABASE_ANON_KEY=
```

#### Reverse Proxy (`reverse-proxy-go/.env`)

```env
R2_PUBLIC_URL=https://pub-xxx.r2.dev
DATABASE_URL=postgresql://user:password@host:5432/database
PORT=8001
```

### Database Setup

The Go services use the pgx driver directly instead of an ORM. Create the required tables in your PostgreSQL database:

```sql
-- Projects table
CREATE TABLE IF NOT EXISTS projects (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    git_url TEXT NOT NULL,
    sub_domain VARCHAR(255) NOT NULL UNIQUE,
    custom_domain VARCHAR(255),
    user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Deployments table
CREATE TABLE IF NOT EXISTS deployments (
    id VARCHAR(36) PRIMARY KEY,
    project_id VARCHAR(36) NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'NOT_STARTED',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ClickHouse logs table
CREATE DATABASE IF NOT EXISTS logs;

USE logs;

CREATE TABLE IF NOT EXISTS log_events (
    event_id String,
    deployment_id String,
    log String,
    timestamp DateTime64(3)
) ENGINE = MergeTree()
ORDER BY (deployment_id, timestamp)
PARTITION BY toYYYYMM(timestamp)
TTL timestamp + INTERVAL 30 DAY;
```

### Running Locally

```bash
# Terminal 1: Start infrastructure (Kafka, ClickHouse)
cd api-server-go
docker-compose up

# Terminal 2: Start API server
cd api-server-go
go run main.go

# Terminal 3: Start Reverse Proxy
cd reverse-proxy-go
go run main.go

# Terminal 4: Start Frontend
cd frontend
npm run dev
```

## Docker Deployment

### Building Docker Images

```bash
# Build API Server
cd api-server-go
docker build -t mini-vercel-api:latest .

# Build Reverse Proxy
cd ../reverse-proxy-go
docker build -t mini-vercel-proxy:latest .

# Build Project Builder
cd ../build-project
docker build -t mini-vercel-build:latest .
```

### Running with Docker

```bash
# Run API Server
docker run -d \
  --name mini-vercel-api \
  -p 9000:9000 \
  --env-file .env \
  mini-vercel-api:latest

# Run Reverse Proxy
docker run -d \
  --name mini-vercel-proxy \
  -p 8001:8001 \
  --env-file .env \
  mini-vercel-proxy:latest
```

### Pushing to Registry

```bash
# Tag and push to your registry
docker tag mini-vercel-api:latest your-registry/mini-vercel-api:latest
docker push your-registry/mini-vercel-api:latest

docker tag mini-vercel-proxy:latest your-registry/mini-vercel-proxy:latest
docker push your-registry/mini-vercel-proxy:latest

docker tag mini-vercel-build:latest your-registry/mini-vercel-build:latest
docker push your-registry/mini-vercel-build:latest
```

## API Endpoints

### Health
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |

### Projects
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/projects` | List all projects |
| GET | `/projects/:id` | Get project details |
| POST | `/projects` | Create new project |
| DELETE | `/projects/:id` | Delete project |

### Deployments
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/projects/:projectId/deployments` | List deployments |
| GET | `/deployments/:id` | Get deployment details |
| GET | `/deployments/:id/logs` | Get deployment logs |
| POST | `/deploy` | Trigger new deployment |

## Deployment Flow

1. **User triggers deployment** via frontend or API
2. **API server creates deployment** record with status `QUEUED`
3. **AWS ECS task spawns** build container with project config
4. **Build container**:
   - Clones Git repository
   - Runs `npm install` and `npm run build`
   - Streams logs to Kafka topic `build-events`
   - Uploads `dist/` folder to Cloudflare R2
5. **API server consumes Kafka logs** using consumer group with worker pool:
   - Worker pool processes messages concurrently (50 workers)
   - Updates status to `IN_PROGRESS` on build start
   - Updates status to `READY` on success
   - Updates status to `FAIL` on error
   - Stores all logs in ClickHouse for querying
6. **Reverse proxy routes traffic** to the deployed assets on R2

## Database Schema

### PostgreSQL Tables

```sql
-- Projects
CREATE TABLE projects (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    git_url TEXT NOT NULL,
    sub_domain VARCHAR(255) NOT NULL UNIQUE,
    custom_domain VARCHAR(255),
    user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Deployments
CREATE TABLE deployments (
    id VARCHAR(36) PRIMARY KEY,
    project_id VARCHAR(36) NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'NOT_STARTED',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Deployment Status Values
-- NOT_STARTED | QUEUED | IN_PROGRESS | READY | FAIL
```

### ClickHouse Logs Table

```sql
CREATE TABLE log_events (
    event_id String,
    deployment_id String,
    log String,
    timestamp DateTime64(3)
) ENGINE = MergeTree()
ORDER BY (deployment_id, timestamp)
PARTITION BY toYYYYMM(timestamp)
TTL timestamp + INTERVAL 30 DAY;
```

## Go Implementation Features

### Authentication
- JWT token verification with RS256 algorithm
- JWKS (JSON Web Key Set) fetching and caching from Supabase
- Automatic key refresh every 24 hours
- Protected routes with user context injection

### Kafka Consumer
- Consumer group with IBM Sarama library
- Worker pool with 50 concurrent workers for parallel message processing
- Graceful shutdown with context cancellation
- TLS support for secure connections

### Repository Pattern
- Clean separation between business logic and data access
- Generic repository interfaces for flexibility
- ClickHouse adapter for log storage with query flexibility
- Mock repositories available for testing

### Error Handling
- Structured error responses with consistent format
- Context aware error logging
- HTTP status code mapping

## Ports Reference

| Service | Port | Protocol |
|---------|------|----------|
| API Server | 9000 | HTTP |
| Frontend (dev) | 5173 | HTTP |
| Reverse Proxy | 8001 | HTTP |
| Zookeeper | 2181 | TCP |
| Kafka | 9092-9094 | TCP |
| ClickHouse HTTP | 8123 | HTTP |
| ClickHouse Native | 9001 | TCP |

## Production Considerations

- Enable HTTPS and TLS on all services
- Use proper secrets management like AWS Secrets Manager or HashiCorp Vault
- Set up monitoring and alerting with Prometheus and Grafana
- Configure CDN for static assets to improve global performance
- Implement deployment rollbacks for quick recovery
- Add build timeout handling to prevent hanging tasks
- Configure database connection pooling for better resource utilization
- Set up proper CORS policies for frontend access
- Enable rate limiting at the API gateway level
- Use environment specific configurations
- Implement log rotation policies in ClickHouse

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## License

MIT License - See LICENSE file for details.
