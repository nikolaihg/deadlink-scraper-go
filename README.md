# 🕸️ Dead Link Web Scraper (in Go)

A simple, recursive dead link checker written in Go.

This project scrapes a given URL, recursively follows internal links, and logs any dead links. Designed as a learning project for exploring Go's core features like concurrency, HTTP/HTML handling, and database interaction.

> Built to practice Go through real-world scraping, concurrency, data persistence, and microservice architecture.

## 🚀 Features

- Accepts a starting URL to begin scraping.
- Recursively follows and checks links **within the same base domain**.
- Detects dead links:
  - Links that return 4xx or 5xx HTTP status codes.
  - Links that timeout.
- Logs all dead links to the console.
- Skips already visited URLs to prevent reprocessing.
- Handles redirects properly (3xx responses).
- Avoids infinite recursion or loops.

### Project Overview
The project is planned to evolve in **three stages**:

- Part 1: CLI Link Scanner
  - Basic recursive scanner using `net/http`, `x/net/html`, and concurrency
- Part 2: Storage and Persistent REST API endpoints
  - Adds PostgreSQL storage and REST endpoints
- Part 3: Scalable Microservice Backend
  - Queue-based job processing, gRPC, observability and async architecture

See [`progress.md`](./progress.md) for development breakdown.

## 🧪 Example

Try scraping this test site:  
🔗 [`https://scrape-me.dreamsofcode.io`](https://scrape-me.dreamsofcode.io)

## ⚙️ How It Works

- Uses the standard `net/http` package for HTTP requests.
- Parses HTML using `golang.org/x/net/html`.
- Maintains a set of visited URLs to avoid rechecking.
- Only scrapes pages within the base domain, but **does** validate external links without crawling them.
- Designed to be simple and extensible.

## 🧠 Design Considerations

### ✅ Handled Edge Cases

- **Redirects:** Follows 3xx redirects and treats them as part of the request lifecycle.
- **Infinite Recursion:** Keeps track of visited URLs to prevent loops.
- **Base Domain Limiting:** Recursively scans only within the original domain; external links are checked but not scraped.

### ❌ Not Yet Supported

- **JavaScript-rendered sites:** These require a headless browser (e.g. Playwright for Go). Could be added as an expansion.
- **Robots.txt or rate limiting:** Not yet respected. Use with caution on real websites.

## 🧵 Potential Improvements

- ✅ Add concurrency with goroutines and channels for faster scanning.
- ⏱ Timeout handling per request.
- 📄 Save results to a file (e.g. JSON or CSV).
- 🌐 Proxy support and user-agent randomization.
- 🧪 Unit tests and structured logging.

## 🛠️ Tech Stack

- [Go](https://golang.org/)
- [`net/http`](https://pkg.go.dev/net/http)
- [`x/net/html`](https://pkg.go.dev/golang.org/x/net/html)

## 🔁 Next Steps
- **[Continue to Part 2 – Persistent Dead Link Monitor](./part2.md)**
- **[Part 3 – Scalable Media Service Architecture](./part3.md)**
- **[`progress.md`](./progress.md) – Feature checklist and roadmap**

**multi-part case study** for portfolio or CV:

- **Part 1**: CLI Tool – Go concurrency, parsing, scraping
- **Part 2**: API Backend – Database modeling, HTTP API, Docker
- **Part 3**: Scalable Microservice – gRPC, Kafka, observability


## Screenshots / Example Output

```bash
$ go run cmd/main.go https://example.com

✅ https://example.com/about
❌ https://example.com/dead-link (404)
⏳ https://example.com/stuck (timeout)

Scan complete. 13 OK, 2 Dead Links.
```