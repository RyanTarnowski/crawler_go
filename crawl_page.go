package main

import (
	"fmt"
	"net/url"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}

	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	isMaxed := cfg.checkMaxPages()
	if isMaxed {
		return
	}

	baseURL, err := url.Parse(cfg.baseURL.String())
	if err != nil {
		fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %v\n", cfg.baseURL.String(), err)
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %v\n", rawCurrentURL, err)
		return
	}

	if baseURL.Hostname() != currentURL.Hostname() {
		return
	}

	normalizeCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - normalizedURL: %v", err)
		return
	}

	isFirst := cfg.addPageVisit(normalizeCurrentURL)
	if !isFirst {
		return
	}

	fmt.Printf("crawling %s\n", rawCurrentURL)
	rawHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - getHTML: %v", err)
		return
	}

	pageData := extractPageData(rawHTML, rawCurrentURL)
	cfg.setPageData(normalizeCurrentURL, pageData)

	for _, v := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(v)
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	_, ok := cfg.pages[normalizedURL]
	if ok {
		return false
	}

	cfg.pages[normalizedURL] = PageData{}
	return true
}

func (cfg *config) setPageData(normalizedURL string, data PageData) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.pages[normalizedURL] = data
}

func (cfg *config) checkMaxPages() (isMaxed bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if len(cfg.pages) >= cfg.maxPages {
		return true
	}
	return false
}
