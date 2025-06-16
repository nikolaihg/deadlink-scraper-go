# Progress – GowebScraper

This document tracks the development status of the GowebScraper project.
### **Go Fundamentals**
* [X] Core syntax, control flow, functions
* [ ] Error handling, modules, packages
* [ ] Concurrency (goroutines)
* [ ] Standard libraries: `net/http`, `log`, `time`.

### Core Scraper Logic
* [ ] Send HTTP requests with `net/http`
* [ ] Parse HTML using `x/net/html`
* [ ] Extract:
  * Page `<title>` or primary heading
  * Main body text (limit \~500 words)
  * All `<a href="">` links
* [ ] Implement delay/rate limiting between requests
* [ ] Recursively discover URLs via internal links
* [ ] Use hashed URLs as unique IDs
* [ ] Respect `robots.txt` directives

### Concurrency & Control
* [ ] Launch concurrent scrapes using goroutines
* [ ] Use channels or sync primitives to manage concurrency
* [ ] Limit parallelism and implement backpressure

### Data Storage
* [ ] Design database schema (`WebPage` object model)
* [ ] Connect to PostgreSQL
* [ ] Insert parsed data into DB
* [ ] Implement indexing/search (exact and optional fuzzy)

### Dev Tools & Quality
* [ ] CLI interface for triggering scrapes or querying DB
* [ ] Structured error handling and logging
* [ ] Unit and integration tests
* [ ] Logging of request status, errors, timestamps

### Deployment
* [ ] Add `.env` support for config (ports, DB URL)

### Bonus Goals
* [ ] REST API (e.g., `GET /latest`, `GET /page/{hash}`)
* [ ] Frontend (React or minimal HTML) to view/search data
* [ ] Optional full-text or fuzzy search
* [ ] Metrics/log dashboard or JSON export of results
