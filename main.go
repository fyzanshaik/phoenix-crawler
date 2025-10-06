package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("usage: crawler URL maxConcurrency maxPages")
		os.Exit(1)
	}

	if len(os.Args) > 4 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	rawBaseURL := os.Args[1]

	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("error parsing maxConcurrency: %v\n", err)
		os.Exit(1)
	}

	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("error parsing maxPages: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %s\n", rawBaseURL)

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("error parsing base URL: %v\n", err)
		os.Exit(1)
	}

	cfg := &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	cfg.wg.Go(func() {
		cfg.crawlPage(rawBaseURL)
	})

	cfg.wg.Wait()

	err = writeCSVReport(cfg.pages, "report.csv")
	if err != nil {
		fmt.Printf("error writing CSV report: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Crawl complete! Found %d pages. Report saved to report.csv\n", len(cfg.pages))
}

func (cfg *config) addPageVisit(pageData PageData) bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if len(cfg.pages) >= cfg.maxPages {
		return false
	}

	normalizedURL := pageData.URL
	if _, exists := cfg.pages[normalizedURL]; exists {
		return false
	}

	cfg.pages[normalizedURL] = pageData
	return true
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
	}()

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}

	if cfg.baseURL.Hostname() != currentURL.Hostname() {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	if _, exists := cfg.pages[normalizedURL]; exists {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	fmt.Printf("crawling: %s\n", rawCurrentURL)
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error getting HTML: %v\n", err)
		return
	}

	pageData := extractPageData(html, rawCurrentURL)
	pageData.URL = normalizedURL

	isFirst := cfg.addPageVisit(pageData)
	if !isFirst {
		return
	}

	for _, nextURL := range pageData.OutgoingLinks {
		cfg.wg.Go(func() {
			cfg.crawlPage(nextURL)
		})
	}
}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "PhoenixCrawler/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("error status code: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return "", fmt.Errorf("invalid content type: %s", contentType)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(bodyBytes), nil
}
