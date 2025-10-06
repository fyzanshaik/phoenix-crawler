# Phoenix - A Web Crawler

A CLI application that generates 'internal links' reports for any website on the internet by crawling each page of the site!

## Usage

```bash
go build -o crawler
./crawler <URL> <maxConcurrency> <maxPages>
```

Example:
```bash
./crawler "https://example.com" 5 100
```

Arguments:
- URL: The website to crawl (starting point)
- maxConcurrency: Maximum number of concurrent HTTP requests (e.g., 5)
- maxPages: Maximum number of pages to crawl (e.g., 100)

## Output

The crawler generates a `report.csv` file with the following columns:
- page_url: The normalized URL of the page
- h1: The main heading (H1 tag) of the page
- first_paragraph: The first paragraph of content (prioritizes main tag)
- outgoing_link_urls: All links found on the page (semicolon-separated)
- image_urls: All image URLs found on the page (semicolon-separated)

This CSV can be opened in Excel, Google Sheets, or any spreadsheet program for analysis.

# Motivation
A web crawler tool or commonly referred as a bot is used by search engine companies to index new websites and update their search results! AI model companies like OpenAI, Meta, Gemini or ANthropic are using their own web search for creating the best query results as well as scraping data for training AI models.

This is my motivation for creating Phoenix which not only can scrape data but also generate link reports (a page that references other pages on the website), and add a markdown file(s) for each page crawled for llms to utilize.

## Build Log

- URL normalization (remove scheme and trailing slashes)
- Extract H1 tags from HTML
- Extract first paragraph from HTML (with main tag priority)
- Extract all links from HTML pages (convert relative to absolute URLs)
- Extract all image URLs from HTML (convert relative to absolute URLs)
- Structured page data extraction (PageData struct)
- CLI argument validation and handling
- HTTP fetching with User-Agent header and content-type validation
- Recursive web crawling with same-domain restriction
- Link reference counting across entire website
- Concurrent crawling with goroutines and WaitGroups
- Thread-safe page tracking with mutexes
- Configurable concurrency control via buffered channels
- Configurable max pages limit to prevent excessive crawling
- CSV report export with full page data (URL, H1, paragraph, links, images)

