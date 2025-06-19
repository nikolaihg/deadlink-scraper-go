# Part 2 – Persistent Dead Link Monitor

In this second part, the project moves from a CLI-only tool to a **minimal persistent backend service**.

It introduces:
- A REST API to submit scan jobs and view results.
- PostgreSQL database to persist:
  - URLs scanned
  - Dead link results
  - Timestamps for audit/history

## New Features

### PostgreSQL Storage

Stores each scan job and its results.

```sql
CREATE TABLE scans (
    id UUID PRIMARY KEY,
    base_url TEXT,
    created_at TIMESTAMP
);

CREATE TABLE dead_links (
    id SERIAL PRIMARY KEY,
    scan_id UUID REFERENCES scans(id),
    href TEXT,
    status_code INT,
    reason TEXT
);
```

Use the `pgx` driver or `database/sql` + `lib/pq`.

### REST API (using `net/http`)

- `POST /scan` – triggers scan for a given URL
- `GET /results` – returns dead link history
- Optional: pagination, filtering by domain

###  Deployment-Ready

* Add Dockerfile + `docker-compose.yml` with:
  * Go app
  * PostgreSQL

## Sample Flow

```bash
curl -X POST localhost:8080/scan \
  -H "Content-Type: application/json" \
  -d '{"url": "https://tv2.no"}'

# Later:
curl localhost:8080/results
```

## Optional Enhancements

* Basic auth or token header for access.
* UI with Go templates or static React frontend.
* Retry logic for flaky links.

## Technologies Introduced

* `Go` for HTTP + database interaction
* `PostgreSQL` for persistent storage
* `Docker Compose` for orchestration
* `UUID` + timestamps for good audit practice

➡️ Ready to go deeper? Check out [Part 3 – Scalable Media Microservice](./part3.md)