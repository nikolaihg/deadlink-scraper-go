# Progress – Dead Link Web Scraper

Tracks the development progress of the Dead Link Web Scraper project. Divided into 3 main parts to reflect learning and architectural complexity.

## PART 1 – CLI Dead Link Scanner (MVP)

Minimal tool to detect dead (4xx/5xx) or timeout links via recursive crawling.

### Core Scraper Logic
- [X] ~Accept single URL input~
- [X] ~Fetch HTML content using `net/http`~
- [X] ~Parse links with `x/net/html`~
- [x] ~Filter internal vs external links~
- [x] ~Filter page navigation links (if href == "#")~
- [ ] Identify and log on current page
- [ ] Track visited links to avoid loops
- [ ] Recursively crawl internal links
- [ ] Identify and log:
  - [ ] Dead links (4xx/5xx)
  - [ ] Timeout / unreachable links
  - [ ] Redirect chains

### CLI Tooling & UX
- [ ] Simple CLI interface (e.g., `go run main.go https://example.com`)
- [ ] Option to set crawl depth
- [ ] Flag for verbosity/debug output

### Concurrency & Control
- [ ] Use goroutines for parallel link checking
- [ ] Use channels for communication
- [ ] Shared memory-safe visited link store (`sync.Map`)
- [ ] Limit parallelism (semaphore or buffered channel)
- [ ] Basic rate limiting between requests

---

## PART 2 – Persistent Backend (PostgreSQL + REST API)

Adds long-term storage, REST interface, and improved crawl control.

### API Endpoints
- [ ] `POST /scan` – Trigger a new crawl
- [ ] `GET /results` – List all scans
- [ ] `GET /results/{id}` – View results of a specific scan
- [ ] `GET /url/{hash}` – Get result by hashed URL

### PostgreSQL Integration
- [ ] Create `scans`, `pages`, and `dead_links` tables
- [ ] Connect Go app to Postgres (`pgx` or `gorm`)
- [ ] Store results with scan metadata (duration, status, timestamps)
- [ ] Write tests for DB layer (insert, fetch, clean)

### DevOps & Configuration
- [ ] Dockerfile for scraper service
- [ ] `docker-compose.yml` with Postgres
- [ ] Load config from `.env` (URL, DB connection, max workers)
- [ ] Seed/test scan in development env

---

## PART 3 – Scalable Architecture (Async Jobs, Monitoring)

Advanced service architecture: distributed jobs, observability, clean architecture.

### Job Queue Architecture
- [ ] Define job format (e.g., JSON or protobuf)
- [ ] Introduce job queue (Redis or Kafka)
- [ ] Producer: API queues scan job
- [ ] Worker: pulls from queue, runs scan, updates DB
- [ ] Track job lifecycle (`queued`, `running`, `done`, `failed`)

### nter-service Communication
- [ ] gRPC between API and scan worker
- [ ] Define and generate shared protobufs

### Observability & Monitoring
- [ ] Structured logs with trace IDs
- [ ] Prometheus metrics:
  - [ ] Number of dead links per scan
  - [ ] Request latencies
  - [ ] Concurrent scan count
- [ ] Optional: Grafana dashboard for scans

### Hardening & Extensions
- [ ] Token-based API access
- [ ] Optional CORS config for frontend
- [ ] Retry queue for failed links (exponential backoff)
- [ ] Robots.txt parser integration

## Optional Features / Future Work
- [ ] Frontend to visualize dead links
- [ ] CSV/JSON export of scan results
- [ ] Full-text search on pages
- [ ] Kubernetes deployment
- [ ] Email alert on dead link scan
- [ ] CLI + API auth keys