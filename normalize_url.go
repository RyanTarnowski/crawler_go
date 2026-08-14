package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(inputURL string) (string, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("NormalizeURL error: %w", err)
	}

	return strings.ToLower(u.Host + strings.TrimSuffix(u.Path, "/")), nil
}
