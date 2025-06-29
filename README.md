# Dead Link Web Scraper written in Go

Deadlink Scraper Go is a fast, concurrent CLI tool that crawls a website to detect and report broken internal and external links.

This project scrapes a given URL, recursively follows internal links, and logs any dead links. Designed as a learning project for exploring Go's core features like concurrency, HTTP/HTML handling, and database interaction.

The second / third part of this project transforms the simple dead link checker into a **scalable, observable microservice**. Something that is suitable for integration with a media monitoring platform or editorial workflow.

## Motivation

> Built to practice Go through real-world scraping, concurrency, data persistence, and microservice architecture.

### Use case
A real world scenario where this application could be used is for quality control for online newspapers, forums or media houses, where high content quality is essential. A critical component of user experience is ensuring that internal and external links remain functional across the site and app. Dead or broken links impact both SEO and audience trust.

## Features
- Accepts a starting URL via CLI to begin scraping and crawling.
- Recursively follows only internal pages (same base domain).
- Link validation via HEAD / GET.
  - Both internal and external `<a href="(...)">`-tags.
- Skips in-page anchors and non-HTTP resources.
- Keeps track of already visited URLs to avoid loops.
- Detailed stats: total, internal/external, alive/dead, per-status code
- HTTP requests
  - Detects dead links: links that return 4xx or 5xx, or that timeout.
  - Handles redirects properly (3xx responses).
- Concurrency with a goroutine worker-pool (channels + WaitGroup) for faster, safe parallelism.
- Designed to be simple and extensible.


### Crawling Rules
1. Only HTML pages under the base domain are fetched.
2. All `<a>`, `<link>`, `<img>`, `<script>`, `<iframe>` URLs are extracted.
3. Internal links → queued for crawling  
   External links → validated but not enqueued
4. Anchor-only links (e.g. `#section`) → skipped as “page links”

## CLI usage
```bash
$ go run cmd/main.go https://example.com -concurrency 10 -timeout 5s -depth 200
```

Usage: deadlink-scraper [flags] <url>  

Flags:
  -timeout duration     HTTP request timeout (default 10s)
  -concurrency int      Number of concurrent workers (default 10)
  -depth int            Max size of total items queued

## How It Works
- Uses the standard `net/http` package for HTTP requests.
- Parses HTML using `golang.org/x/net/html`.
- Contains custom types for links, sets, stats
  - `linktype/link.go`, `set/set.go`, `stats/stats.go`

### Concurrency
1. **Worker-pool**  
   We spin up N goroutines (workers) that all listen on a single `jobs chan Link`.
2. **Task tracking**  
   A `sync.WaitGroup` (“taskWg”) is incremented every time we enqueue a new link and decremented when a worker finishes processing one.  
   When it drops to zero, we know there’s no more work, so we close the channel.
3. **Safe shared state**  
   We protect the shared `visited` and `checked` sets, plus our `LinkStats` counters, with `sync.Mutex` locks to avoid races.
4. **Lifecycle**  
   - Main seeds the first URL and adds it to the WaitGroup.  
   - Workers pull links off `jobs`, call `crawl()`, then signal Done.  
   - Any newly discovered internal links are `Add`ed to the WaitGroup and sent back into `jobs`.  
   - Once `taskWg.Wait()` unblocks, we `close(jobs)`, workers exit, and we print stats.

### Diagrams
(Insert flow charts)

## Project Overview
The project is planned to evolve in **three stages**:

- [Part 1: CLI Link Scanner](./docs/part1.md)
  - Basic recursive scanner using `net/http`, `x/net/html`, and concurrency
- [Part 2: Storage and Persistent REST API endpoints](./docs/part2.md)
  - Adds PostgreSQL storage and REST endpoints
- [Part 3: Scalable Microservice Backend](./docs/part3.md)
  -  Queue-based job processing, gRPC, observability and async architecture

See [progress.md](./docs/progress.md) for development breakdown.

##  Design Considerations

###  Handled Edge Cases

- **Redirects:** Follows 3xx redirects and treats them as part of the request lifecycle.
- **Infinite Recursion:** Keeps track of visited URLs to prevent loops.
- **Base Domain Limiting:** Recursively scans only within the original domain; external links are checked but not scraped.

###  Not Yet Supported

- **JavaScript-rendered sites:** These require a headless browser (e.g. Playwright for Go). Could be added as an expansion.
- **Robots.txt** Not yet respected. Use with caution on real websites.

###  Potential Improvements
- Unit tests.
- Structured logging.
- Benchmarking.
- Fuzzing

## Example Output

```bash
$ .\deadlink-scraper-go.exe https://example.com
2025/06/23 22:24:53 [ALIVE]  https://example.com (200 OK)
2025/06/23 22:24:53 [Crawling]: https://example.com
2025/06/23 22:24:54 [ALIVE]  https://www.iana.org/domains/example (200 OK)
2025/06/23 22:24:54 Scan complete:
2025/06/23 22:24:54 Total:    2
2025/06/23 22:24:54 Internal: 1
2025/06/23 22:24:54 External: 1
2025/06/23 22:24:54 Alive:    2
2025/06/23 22:24:54 Dead:     0
2025/06/23 22:24:54 Skipped:  0
2025/06/23 22:24:54 Status codes distribution:
2025/06/23 22:24:54   200: 2
2025/06/23 22:24:54 Links visisted: [{https://example.com 0}]
```

Try scraping this test site:  
🔗 [`https://scrape-me.dreamsofcode.io`](https://scrape-me.dreamsofcode.io)


## Performance Benchmarking
I also benchmarked the non concurrent version and the concurrent version to compare the two versions of the program. 
There are binaries located in the `builds/` folder:
  * `deadlink-scraper-go-nonconcurrent` — Non-concurrent version
  * `deadlink-scraper-go-concurrent` — Concurrent version

In `tools/benchmarks/`, there is a program (`benchmarks.go`) that benchmarks the two builds and compares them.

The script will:
* Run each binary 5 times (configurable in the script)
* Print execution time for each run
* Print average execution time for both versions

```
$ go run .\benchmarks.go

Running NonConcurrent version:
Run 1: 26.8868068s
Run 2: 29.7986582s
Run 3: 27.5302179s
Run 4: 27.1635429s
Run 5: 27.3745297s

Average time for NonConcurrent: 27.7507511s

Running Concurrent version:
Run 1: 11.1094938s
Run 2: 11.4766514s
Run 3: 11.7514964s 
Run 4: 11.5983008s
Run 5: 11.5512609s

Average time for Concurrent: 11.49744066s
```

##  Next Steps
- **[See Part 1 - CLI Dead Link Scanner](./docs/part1.md)**
- **[Continue to Part 2 – Persistent Dead Link Monitor](./docs/part2.md)**
- **[Skip to Part 3 – Scalable Media Service Architecture](./docs/part3.md)**
- **[progress.md](./docs/progress.md) – Feature checklist and roadmap**
