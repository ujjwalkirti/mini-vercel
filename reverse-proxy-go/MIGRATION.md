# Migration from TypeScript to Go

This document outlines the migration of the reverse proxy from TypeScript to Go.

## Architecture Comparison

### TypeScript Implementation
```
reverse-proxy/
├── src/
│   ├── reverse-proxy.ts          # Main entry point + all logic
│   └── lib/
│       └── prisma.ts              # Prisma client
├── prisma/
│   └── schema.prisma              # Database schema
├── package.json
└── Dockerfile
```

**Tech Stack:**
- Express.js for HTTP server
- http-proxy for reverse proxying
- Prisma ORM for database queries
- TypeScript with tsx for development

### Go Implementation
```
reverse-proxy-go/
├── main.go                        # Entry point
├── internal/
│   ├── app/
│   │   └── app.go                # Application initialization
│   ├── config/
│   │   └── config.go             # Configuration management
│   ├── db/
│   │   └── postgres.go           # Database connection
│   ├── domain/
│   │   ├── project/model.go      # Project model
│   │   └── deployment/model.go   # Deployment model
│   ├── handler/
│   │   └── proxy/
│   │       ├── handler.go        # Proxy logic
│   │       └── routes.go         # Route registration
│   ├── middleware/
│   │   └── logger.go             # Logging middleware
│   ├── repository/
│   │   └── project/postgres.go   # Database queries
│   └── router/
│       └── router.go             # Router setup
├── go.mod
├── .air.toml
└── Dockerfile
```

**Tech Stack:**
- chi/v5 for HTTP routing
- net/http/httputil for reverse proxying (standard library)
- database/sql with lib/pq for PostgreSQL
- Air for hot reload development

## Feature Parity

| Feature | TypeScript | Go | Status |
|---------|-----------|-----|--------|
| Subdomain routing | ✅ | ✅ | ✅ |
| Custom domain support | ✅ | ✅ | ✅ |
| PostgreSQL queries | ✅ (Prisma) | ✅ (database/sql) | ✅ |
| R2 proxy | ✅ | ✅ | ✅ |
| Index.html default | ✅ | ✅ | ✅ |
| Graceful shutdown | ✅ | ✅ | ✅ |
| Health check | ✅ | ✅ | ✅ |
| CORS support | ✅ | ✅ | ✅ |
| Hot reload | ✅ (tsx) | ✅ (air) | ✅ |
| Docker support | ✅ | ✅ | ✅ |

## Key Differences

### 1. Database Layer
**TypeScript:**
```typescript
const project = await prismaClient.project.findUnique({
  where: { subDomain: subdomain },
  include: {
    deployments: {
      where: { status: 'READY' },
      orderBy: { createdAt: 'desc' },
      take: 1
    }
  }
});
```

**Go:**
```go
query := `
  SELECT p.*, d.*
  FROM projects p
  LEFT JOIN deployments d ON p.id = d.project_id
  WHERE p.subdomain = $1
  AND d.status = $2
  ORDER BY d.created_at DESC
  LIMIT 1
`
err := db.QueryRowContext(ctx, query, subdomain, "READY").Scan(...)
```

### 2. Reverse Proxy
**TypeScript:**
```typescript
const proxy = httpProxy.createProxyServer();
proxy.web(req, res, {
  target: `${R2_URL}/${projectId}`,
  changeOrigin: true
});
```

**Go:**
```go
target, _ := url.Parse(targetURL)
proxy := httputil.NewSingleHostReverseProxy(target)
proxy.ServeHTTP(w, r)
```

### 3. Middleware
**TypeScript:**
```typescript
app.use(express.json());
app.use(cors());
```

**Go:**
```go
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
r.Use(cors.Handler(cors.Options{...}))
```

## Performance Benefits

- **Memory**: Go uses ~10-20MB vs Node.js ~50-100MB
- **Startup**: Go cold start ~50ms vs Node.js ~500ms
- **Concurrency**: Go goroutines are more efficient than Node.js event loop for I/O-bound operations
- **Binary size**: Single ~10MB binary vs ~100MB+ node_modules

## Development Workflow

### TypeScript
```bash
npm run dev          # Development with tsx watch
npm run build        # Build to JavaScript
npm start            # Run production build
```

### Go
```bash
air                  # Development with hot reload
go build            # Build binary
./main              # Run binary
```

## Environment Variables

Both implementations use the same environment variables:

```env
DATABASE_URL=postgresql://...
R2_PUBLIC_URL=https://pub-xxx.r2.dev
PORT=8001
```

## Testing Both Implementations

You can run both implementations side by side on different ports:

```bash
# TypeScript on port 8001
cd reverse-proxy
npm run dev

# Go on port 8002
cd reverse-proxy-go
PORT=8002 air
```

## Migration Checklist

- [x] Set up Go project structure
- [x] Implement config loading
- [x] Set up PostgreSQL connection
- [x] Create domain models
- [x] Implement repository layer
- [x] Create proxy handler
- [x] Set up routing and middleware
- [x] Add graceful shutdown
- [x] Create Dockerfile
- [x] Add hot reload configuration
- [x] Document the implementation

## Next Steps

1. Test the Go implementation with actual projects
2. Update docker-compose.yml to use Go version
3. Monitor performance and memory usage
4. Gradually deprecate TypeScript version
5. Remove TypeScript version after validation

## Rollback Plan

If issues arise, the TypeScript implementation remains available:
- Original code is in `reverse-proxy/` directory
- Docker image: Use old Dockerfile
- No database schema changes required
