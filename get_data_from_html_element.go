package main

import (
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) string {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
		return ""
	}

	//Get the h1 text or h2 if h1 isn't present
	headerText := doc.Find("h1, h2").First().Text()

	return strings.TrimSpace(headerText)
}

func getFirstParagraphFromHTML(html string) string {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
		return ""
	}

	//Get the first p text, in main if main exists
	paraText := doc.Find("main p").First().Text()
	if paraText == "" {
		paraText = doc.Find("p").First().Text()
	}

	return strings.TrimSpace(paraText)
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return []string{}, err
	}

	urls := []string{}

	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		// For each '<a href>' it finds, it will run this function.
		aurl, exists := s.Attr("href")
		if exists {
			if !strings.Contains(aurl, baseURL.String()) {
				fullURL := baseURL.JoinPath(aurl)
				aurl = fullURL.String()
			}

			urls = append(urls, aurl)
		}
	})

	return urls, nil
}
