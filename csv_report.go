package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func writeCSVReport(pages map[string]PageData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("error writing headers: %w", err)
	}

	for _, pageData := range pages {
		outgoingLinks := strings.Join(pageData.OutgoingLinks, ";")
		imageURLs := strings.Join(pageData.ImageURLs, ";")

		row := []string{
			pageData.URL,
			pageData.H1,
			pageData.FirstParagraph,
			outgoingLinks,
			imageURLs,
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing row: %w", err)
		}
	}

	return nil
}
