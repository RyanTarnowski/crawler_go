package main

import (
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string
	Heading        string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func extractPageData(html, pageURL string) PageData {
	outgoingLinks := []string{}
	imageURLs := []string{}

	srcURL, err := url.Parse(strings.TrimSpace(pageURL))
	if err == nil {
		outgoingLinks, _ = getURLsFromHTML(html, srcURL)
		imageURLs, _ = getImagesFromHTML(html, srcURL)
	}

	return PageData{
		URL:            pageURL,
		Heading:        getHeadingFromHTML(html),
		FirstParagraph: getFirstParagraphFromHTML(html),
		OutgoingLinks:  outgoingLinks,
		ImageURLs:      imageURLs,
	}
}

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
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		hrefURL, err := url.Parse(strings.TrimSpace(href))
		if err != nil {
			return
		}

		resolvedURL := baseURL.ResolveReference(hrefURL)
		urls = append(urls, resolvedURL.String())
	})

	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return []string{}, err
	}

	imgs := []string{}

	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		// For each '<a href>' it finds, it will run this function.
		src, exists := s.Attr("src")
		if !exists {
			return
		}

		srcURL, err := url.Parse(strings.TrimSpace(src))
		if err != nil {
			return
		}

		resolvedURL := baseURL.ResolveReference(srcURL)
		imgs = append(imgs, resolvedURL.String())
	})

	return imgs, nil
}
