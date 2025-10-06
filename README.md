# Phoenix - A Web Crawler

A CLI application that generates 'internal links' reports for any website on the internet by crawling each page of the site!

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

