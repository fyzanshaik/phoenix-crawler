package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func getH1FromHTML(html string) (string, error) {
	htmlReader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return "", fmt.Errorf("error creating document from input string: %w", err)
	}

	h1Selector := doc.Find("h1")
	// fmt.Println(h1Selector.Html())
	return h1Selector.Text(), nil
}

func getFirstParagraphFromHTML(html string) (string, error) {
	htmlReader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return "", fmt.Errorf("error creating document from input string: %w", err)
	}

	mainSelector := doc.Find("main")
	if mainSelector.Length() > 0 {
		pSelector := mainSelector.Find("p").First()
		return pSelector.Text(), nil
	}

	pSelector := doc.Find("p").First()
	return pSelector.Text(), nil
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	htmlReader := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return nil, fmt.Errorf("error creating document from input string: %w", err)
	}

	var urls []string
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		parsedURL, err := url.Parse(href)
		if err != nil {
			return
		}

		absoluteURL := baseURL.ResolveReference(parsedURL)
		urls = append(urls, absoluteURL.String())
	})

	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	htmlReader := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return nil, fmt.Errorf("error creating document from input string: %w", err)
	}

	var images []string
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}

		parsedURL, err := url.Parse(src)
		if err != nil {
			return
		}

		absoluteURL := baseURL.ResolveReference(parsedURL)
		images = append(images, absoluteURL.String())
	})

	return images, nil
}

func extractPageData(html, pageURL string) PageData {
	baseURL, err := url.Parse(pageURL)
	if err != nil {
		return PageData{URL: pageURL}
	}

	h1, _ := getH1FromHTML(html)
	firstParagraph, _ := getFirstParagraphFromHTML(html)
	outgoingLinks, _ := getURLsFromHTML(html, baseURL)
	imageURLs, _ := getImagesFromHTML(html, baseURL)

	return PageData{
		URL:            pageURL,
		H1:             h1,
		FirstParagraph: firstParagraph,
		OutgoingLinks:  outgoingLinks,
		ImageURLs:      imageURLs,
	}
}
