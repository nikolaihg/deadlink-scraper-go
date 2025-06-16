# GoWebScraper
Simple webscraper using the [go](https://go.dev/). Developed to explore concurrency, data handling and analysis through a media-relevant use cases.

This project is part of a learning journey to become proficient in **Go for backend development**, with a tech stack and architectural style inspired by modern develompent practices.

## Idea
Build a domain-specific web scraper that collects data such as **film listings**, **sports results**, or **podcast / music metadata**. And then stores it in a relational database. The applications is built with a strong focuse on:
- Learning Go throug a practical problem / implementation.
- Exploring concurrency, error handling, and database interaction.
- Deploying services using Docker and 'docker-compose`. 
- *Maybe* Creating a reusable and extensible backend pipeline for data scraping. 

Knowledge from previous courses ([INFO215](https://www4.uib.no/en/studies/courses/info215)) and languages (Python) is leveraged—especially DOM traversal, structured data extraction, and web protocols.

## Project goals
### Data Collection
* Scrape data from defined URLs.
* Discover relevant URLs from `<a>`-tags.
* Parse pages for:
  * Title/Header
  * Text nodes (store \~500 words of content)
  * `<a>` links
* Index pages with hashed URL as primary key.

### Storage
* Store structured data in a **PostgreSQL** database.
* Implement search capabilities (exact and *maybe* fuzzy).
* Enable CRUD operations via the Go application.

### Architecture and Concurrency
* Use **goroutines** for concurrent scraping.
* Implement delays and rate limiting.
* Respect `robots.txt`.
* Implement robust **error handling and logging**.

### CLI & API
* Create a simple CLI interface for scraper configuration or triggering.
* *Bonus:* Serve a minimal REST API (e.g., `GET /latest`).
* *Optional:* Build a frontend to view and search results.

## Progress
Progress is tracked in [`progress.md`](./progress.md).

## Tech Stack
| Area             | Tools                                          |
| ---------------- | ---------------------------------------------- |
| Language         | Go                                             |
| Libraries        | `net/http`, `x/net/html`, `log`, `testing`     |
| Database         | PostgreSQL                                     |
| Bonus (optional) | REST API (`net/http`), React frontend, CLI     |

## WebPage Object Model
The core data structure for scraped content is the **WebPage** object, which represents a snapshot of a single parsed webpage. It includes:

```go
type WebPage struct {
    URL         string   // The original URL
    URLHash     string   // A hashed version of the URL (used as primary key)
    Title       string   // Page title or main header
    Text        string   // Up to ~500 words of main content
    Links       []string // All discovered <a href=""> URLs
    Timestamp   time.Time // Time of scrape
}

This structure will be used consistently in parsing logic and in the database schema.
