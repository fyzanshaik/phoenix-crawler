package main

import (
	"fmt"
	"net/url"
)

/*
Input:
https://blog.boot.dev/path/
https://blog.boot.dev/path
http://blog.boot.dev/path/
http://blog.boot.dev/path

Output: blog.boot.dev/path
*/
func normalizeURL(link string) (string, error) {
	parsedUrl, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("error parsing url: %w", err)
	}
	if len(parsedUrl.Path) > 0 && parsedUrl.Path[len(parsedUrl.Path)-1] == '/' {
		parsedUrl.Path = parsedUrl.Path[:len(parsedUrl.Path)-1]
	}
	normalizedUrl := parsedUrl.Hostname() + parsedUrl.Path
	return normalizedUrl, nil
}
