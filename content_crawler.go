package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// CrawlURL fetches and extracts content from a single URL
func CrawlURL(url string) (*CrawledContent, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return nil, fmt.Errorf("not an HTML page: %s", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	content := ExtractContent(doc)
	content.SourceURL = url
	return content, nil
}

// ExtractContent extracts meaningful content from an HTML document
func ExtractContent(doc *html.Node) *CrawledContent {
	content := &CrawledContent{}

	// Extract title
	if title := findTitle(doc); title != "" {
		content.Title = title
	}

	// Extract meta description
	if desc := findMetaDescription(doc); desc != "" {
		content.Description = desc
	}

	// Extract meta keywords
	if keywords := findMetaKeywords(doc); keywords != "" {
		content.Keywords = keywords
	}

	// Extract headings (h1, h2, h3, etc.)
	content.Headings = extractHeadings(doc)

	// Extract paragraphs
	content.Paragraphs = extractParagraphs(doc)

	// Extract links
	content.Links = extractLinks(doc)

	return content
}

// findTitle extracts the page title
func findTitle(doc *html.Node) string {
	return findMetaTag(doc, "title")
}

// findMetaDescription extracts the meta description
func findMetaDescription(doc *html.Node) string {
	return findMetaTag(doc, "description")
}

// findMetaKeywords extracts the meta keywords
func findMetaKeywords(doc *html.Node) string {
	return findMetaTag(doc, "keywords")
}

// findMetaTag finds a meta tag by its name attribute
func findMetaTag(doc *html.Node, name string) string {
	var findMeta func(*html.Node) string
	findMeta = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "meta" {
			for _, attr := range n.Attr {
				if strings.EqualFold(attr.Key, "name") && strings.EqualFold(attr.Val, name) {
					// Find the content attribute
					for _, contentAttr := range n.Attr {
						if strings.EqualFold(contentAttr.Key, "content") {
							return contentAttr.Val
						}
					}
				}
				if strings.EqualFold(attr.Key, "property") && strings.EqualFold(attr.Val, "og:title") {
					for _, contentAttr := range n.Attr {
						if strings.EqualFold(contentAttr.Key, "content") {
							return contentAttr.Val
						}
					}
				}
			}
		}
		if n.Type == html.ElementNode && n.Data == "title" {
			if name == "title" {
				var title strings.Builder
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						title.WriteString(c.Data)
					}
				}
				return title.String()
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if result := findMeta(c); result != "" {
				return result
			}
		}
		return ""
	}
	return findMeta(doc)
}

// extractHeadings extracts all heading text (h1-h6)
func extractHeadings(doc *html.Node) []string {
	var headings []string
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && isHeading(n.Data) {
			var text strings.Builder
			extractText(n, &text)
			if text.String() != "" {
				headings = append(headings, text.String())
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}
	extract(doc)
	return headings
}

// isHeading checks if a tag is a heading (h1-h6)
func isHeading(tag string) bool {
	return tag == "h1" || tag == "h2" || tag == "h3" ||
		tag == "h4" || tag == "h5" || tag == "h6"
}

// extractParagraphs extracts paragraph text
func extractParagraphs(doc *html.Node) []string {
	var paragraphs []string
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			var text strings.Builder
			extractText(n, &text)
			trimmedText := strings.TrimSpace(text.String())
			if trimmedText != "" && len(trimmedText) > 10 { // Filter out very short paragraphs
				paragraphs = append(paragraphs, trimmedText)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}
	extract(doc)
	return paragraphs
}

// extractLinks extracts all href links
func extractLinks(doc *html.Node) []string {
	var links []string
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}
	extract(doc)
	return links
}

// extractText extracts text content from a node and its children
func extractText(n *html.Node, text *strings.Builder) {
	if n.Type == html.TextNode {
		text.WriteString(strings.TrimSpace(n.Data))
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, text)
	}
}

// CrawlResults crawls multiple URLs and updates the search results with crawled content
func CrawlResults(searchResult *SearchResult, maxLinks int) {
	if maxLinks <= 0 {
		return
	}

	count := 0
	for i := range searchResult.Results {
		if count >= maxLinks {
			break
		}

		if searchResult.Results[i].Link != "" {
			content, err := CrawlURL(searchResult.Results[i].Link)
			if err != nil {
				fmt.Printf("Warning: failed to crawl %s: %v\n", searchResult.Results[i].Link, err)
				continue
			}
			searchResult.Results[i].CrawledContent = content
			count++
		}
	}
}
