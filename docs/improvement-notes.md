### Notes for improvements to each of the functions in the application

- `extractLinks` :
    - Add context or logger support, if this will be used in larger-scale crawlers.
    - Unit tests for commom HTML input variations.
    - Metrics / counters (how many links of each type fiund) if needed for diagnostics

- `filterLinks` :
    - Add retry logic or timeout handling in `fetchStatus`.
    - Throttle or rate-limit external link validation if it's used in bulk.
    - Add contextual metadata to logs (e.g. sourcepage or tag name) for better debugging.

- `fetchStatus` :
    - Support robots.txt.
    - Implement caching or deduplication for repeated URLs.
    - Use `HEAD` with fallback logic based on status codes only (not errors) to prevent issues like network timeouts or DNS erros from causing false negatives. 

- `normalize` :
    - Collapse duplicate slashes in paths (e.g. `/foo//bar` $\rightarrow$ `/foo/bar`).

- `crawl` :
    - Add depth limiting or page count limit to avoid infinite crawling.
    - Retry logic for trasient fetch failures (e.g. 5xx)