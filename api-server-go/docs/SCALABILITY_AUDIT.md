# API Server Go — Scalability Audit

Prioritized list of missing industry best practices and bottlenecks, ordered by severity.

---

## Critical Severity (Will break under load)

- [ ] **1. No Database Connection Pool Configuration**
  `postgres.go` calls `sql.Open()` and `db.Ping()` but never sets `MaxOpenConns`, `MaxIdleConns`, `ConnMaxLifetime`, or `ConnMaxIdleTime`. Under load, Go's default `sql.DB` opens **unlimited connections**, exhausting the PostgreSQL connection limit and crashing the database.
  ```go
  db.SetMaxOpenConns(25)
  db.SetMaxIdleConns(5)
  db.SetConnMaxLifetime(5 * time.Minute)
  db.SetConnMaxIdleTime(5 * time.Minute)
  ```

- [ ] **2. No Rate Limiting**
  The router has zero rate limiting. A single user or attacker can exhaust all server and database resources. Need per-user and per-IP rate limits using `golang.org/x/time/rate` or a Redis-based limiter.

- [ ] **3. No Pagination on List Endpoints**
  Project and deployment handlers return **all rows**. `GetProjects` uses `json_agg` to embed all deployments per project with no `LIMIT`/`OFFSET`. At scale, a single request could load thousands of records into memory, causing OOM.

- [ ] **4. Kafka Consumer Inserts Logs One-by-One to ClickHouse**
  Each Kafka message triggers an individual `INSERT` into ClickHouse. ClickHouse is optimized for **batch inserts** (1000+ rows at a time). Individual inserts cause massive overhead and will eventually lead to rejected writes.

- [ ] **5. No Graceful HTTP Server Shutdown**
  `app.go` handles Kafka shutdown but the HTTP server has no `server.Shutdown(ctx)`. During deploys, in-flight requests get killed mid-response, causing data corruption and user-facing errors.

---

## High Severity (Significant performance degradation)

- [ ] **6. Zero Caching Layer**
  No Redis, no in-memory cache, no HTTP `Cache-Control` headers. Every request hits PostgreSQL directly. Repeated queries for the same data multiply DB load by 10–100x unnecessarily.

- [ ] **7. No Structured Logging**
  Uses `log.Println` / `log.Printf` everywhere — no levels, no request IDs in context, no structured fields. In production with multiple replicas, debugging is nearly impossible. Use `zerolog` or `zap`.

- [ ] **8. Synchronous ECS Task Launch in Deploy Handler**
  The deployment handler calls ECS `RunTask` synchronously, blocking the HTTP response. If ECS is slow or rate-limited, API threads are tied up. Deployments should be queued asynchronously.

- [ ] **9. No Circuit Breakers for External Services**
  Calls to ECS, ClickHouse, Supabase, and Kafka have no circuit breaker. If any external dependency goes down, the API hangs on every request instead of failing fast. Use `github.com/sony/gobreaker`.

- [ ] **10. No Request Timeouts / Context Deadlines**
  Handlers and DB queries don't set context deadlines. A slow query or hung external call blocks a goroutine indefinitely, eventually exhausting resources.

---

## Medium Severity (Operational and security gaps)

- [ ] **11. Secrets in Plain Environment Variables**
  AWS keys, JWT secrets, Kafka passwords, ClickHouse credentials — all in env vars. At scale, use AWS Secrets Manager or HashiCorp Vault with rotation.

- [ ] **12. No Response Compression**
  No `gzip`/`deflate` middleware in the Chi router. For JSON-heavy APIs, this wastes 60–80% of bandwidth. Chi has `middleware.Compress` available.

- [ ] **13. No API Versioning**
  Routes are `/projects`, `/deploy` — no `/api/v1/` prefix. Any breaking change forces all clients to update simultaneously, making rollouts dangerous.

- [ ] **14. No Distributed Tracing or Metrics**
  No OpenTelemetry, no Prometheus `/metrics` endpoint. Can't identify slow endpoints, measure latencies, or track error rates.

- [ ] **15. No Health Check Dependency Verification**
  `/health` returns 200 without checking if Postgres, ClickHouse, and Kafka are actually reachable. Load balancers will route traffic to broken instances.

- [ ] **16. JWKS Cache Not Distributed**
  `jwks_cache.go` uses in-memory caching with 5-min TTL. With multiple replicas, every instance independently fetches JWKS, and simultaneous cache misses cause a thundering herd against Supabase.

- [ ] **17. No Dead Letter Queue for Kafka**
  Failed Kafka messages in the consumer are lost. Intermittent ClickHouse failures mean permanent log loss. Need a DLQ to retry failed messages.

- [ ] **18. ClickHouse Deletes Use `ALTER TABLE DELETE`**
  Deleting logs uses `ALTER TABLE ... DELETE` which is a **mutation** in ClickHouse — extremely expensive and async. At scale, queued mutations pile up and degrade the cluster.

---

## Lower Severity (Best practices gaps)

- [ ] **19. No Request Body Size Limits**
  No `http.MaxBytesReader` — a client can POST a multi-GB body and exhaust memory.

- [ ] **20. No Idempotency on Deploy Endpoint**
  Double-clicking deploy creates duplicate ECS tasks. Add idempotency keys.

- [ ] **21. No Readiness vs Liveness Probe Separation**
  Kubernetes needs separate probes — liveness (is the process alive?) vs readiness (can it serve traffic?). Only one health endpoint exists.

- [ ] **22. Fixed Kafka Worker Pool Size (50)**
  Hardcoded in code, not configurable. Can't tune without redeploying.

- [ ] **23. No Retry with Exponential Backoff**
  External service calls (ECS, ClickHouse) fail without retry. Transient errors become permanent failures.

- [ ] **24. Missing Security Headers**
  No `Strict-Transport-Security`, `X-Content-Type-Options`, `X-Frame-Options`, or CSP headers.
