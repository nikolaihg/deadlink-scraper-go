# Part 1 – CLI Dead Link Scanner
The first part of this project lays the foundation: a fast, recursive CLI tool to detect broken links across an entire website.

> It is focused on learning core Go concepts such as HTTP clients, HTML parsing, concurrency (via goroutines and channels), and building reliable, maintainable tooling from scratch.

## Goals
- Build a basic but useful **command-line tool** that can:
  - Crawl internal links on a given website
  - Identify dead or unreachable links
  - Distinguish internal vs. external references
  - Log link status and errors clearly to the console
- Emphasize performance and correctness over completeness
- Serve as a strong starting point for future expansion into a persistent backend or service

## Key Concepts Introduced
### 1. **Recursive Web Crawling**
- Begin at a single base URL.
- Fetch HTML, extract all relevant URLs.
- Recursively follow internal links (same hostname).
- Avoid infinite loops via visited set.
- Respect only certain tag types (`<a>`, `<link>`, `<img>`, `<script>`, `<iframe>`).
### 2. **Link Validation**
- Each link is validated using HTTP HEAD (or fallback to GET).
- A link is considered "dead" if it:
  - Returns a 4xx or 5xx status code
  - Times out
  - Fails to resolve
- Supports following redirects (3xx), but not beyond base domain.
### 3. **Concurrency and Workers**
- Crawling and link-checking benefit from high parallelism.
- Uses goroutines to check many links in parallel.
- Channels and/or semaphore pattern manage safe concurrent access.
- Future-proof design for scaling to job-based async model.
### 4. **CLI Design**
- The tool is invoked from the command line and takes:
  - A target URL
  - Optional flags: timeout, concurrency, crawl depth
- Designed to be fast, portable, and stateless.

## Example Usage
```bash
$ go run cmd/main.go https://example.com -concurrency 10 -timeout 5s -depth 200
```

Usage: deadlink-scraper [flags] <url>  

Flags:
  -timeout duration     HTTP request timeout (default 10s)
  -concurrency int      Number of concurrent workers (default 10)
  -depth int            Max size of total items queued
  
---
- [Part 2 – Persistent Dead Link Monitor](./part2.md)
- [Part 3 – Scalable Microservice Backend](./part3.md)
- Back to the [README.md](../README.md)