package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func getHTML(rawURL string) (string, error) {
	//new client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	//create request
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("Error creating request: %v", err)
	}

	//set User-Agent
	req.Header.Set("User-Agent", "BootCrawler/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Bad status: %s", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("Content-Type is not html, Content-Type: %s", contentType)
	}

	//read the resp.Body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading content: %v", err)
	}

	return string(bodyBytes), nil
}
