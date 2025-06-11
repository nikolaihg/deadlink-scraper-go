# gowebscraper
Simple webscraper using the [go programming language](https://go.dev/).

## Idea
Make a webscraper using the go standard library (as much as possible), store the result in a database ([postgreSQL](https://www.postgresql.org/))

I also want the project to enable domain specific scarping to cater to relevance to media or user-facing data by choosing a domain-specific scraping target (e.g. film listings, sports results, or even podcast metadata). (How this is going to be implemented is not yet clear, maybe to predfined url-list).

## Project goals
- Learn [go](https://go.dev/) through a project.
- Apply webscraping + traversing DOM knowlegde from INFO212 & Python into learning golang.
    - Using [net/http](https://pkg.go.dev/net/http) and [x/net/html](https://pkg.go.dev/golang.org/x/net/html)
- Learn database usage in go through storing the results.
    - Want to implement search (maybe fuzzy) from the database.
    - CRUD via the go application.
- Learn concurrency / threading in go (`goroutines`).
- Error handling and logging in go.
- Dockerize the application & server.
- Write tests in go.

## Progress
- [x] Go fundamentals (Syntax, Control flow, functions, Error handling, Concurrency, Standard Library)
- [ ] HTTP requests with `net/http`
- [ ] Implement delays between requests
- [ ] HTML parsing using `/x/net/html`
- [ ] Write some tests in go
- [ ] Concurrent data scarping using goroutines
- [ ] Database connection / implemenation
- [ ] Store scraped data to database
- [ ] CLI usage of the application
- [ ] More testing.
- [ ] Good error handling and logging
- [ ] Multi stage build (Dockerfile)
- [ ] Dockercompose (database container + webscraper container)

### Bonus goals:
- Respect robots.txt
- Implement searching (binary, fuzzy).
- Write fronted to view and search the results.
- CI/CD with github actions and automated testing.

## Tech stack
So far only planned:
- `Go`
    - `net/http`
    - `/x/net/html`
    - `log`
    - `testing`
    - either `"raw SQL"` or `gorm`
- Database: `PostgreSQL`
- Docker: `docker-compose`
