package main

import (
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// FetchSearch performs a search using the specified engine and returns results
func FetchSearch(query, engine string) (*SearchResult, error) {
	var data []byte
	var err error

	switch engine {
	case "ddg", "duckduckgo":
		data, err = fetchDuckDuckGo(query)
	case "bing", "":
		data, err = fetchBing(query)
	default:
		return nil, fmt.Errorf("unsupported search engine: %s", engine)
	}

	if err != nil {
		return nil, err
	}

	// Detect format and parse
	if isXML(data) {
		result, err := ParseSearchXML(data)
		if err != nil {
			return nil, err
		}
		result.SearchInfo.Query = query
		return result, nil
	}

	result, err := ParseSearchHTML(data)
	if err != nil {
		return nil, err
	}
	result.SearchInfo.Query = query
	return result, nil
}

// fetchDuckDuckGo performs a search using DuckDuckGo
func fetchDuckDuckGo(query string) ([]byte, error) {
	searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch search results: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// fetchBing performs a search using Bing
func fetchBing(query string) ([]byte, error) {
	searchURL := fmt.Sprintf("https://www.bing.com/search?q=%s", url.QueryEscape(query))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch search results: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// readResponseBody reads the response body, handling gzip compression
func readResponseBody(resp *http.Response) ([]byte, error) {
	var rc io.ReadCloser = resp.Body

	// Check if response is gzip encoded
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		rc = gz
		defer rc.Close()
	}

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ParseSearchXML parses the XML response from search and converts it to SearchResult
func ParseSearchXML(xmlData []byte) (*SearchResult, error) {
	var rss GoogleSearchXML
	if err := xml.Unmarshal(xmlData, &rss); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	results := make([]Result, 0, len(rss.Channel.Items))
	for _, item := range rss.Channel.Items {
		result := Result{
			Title:       item.Title,
			Link:        item.Link,
			Snippet:     strings.TrimSpace(item.Description),
			DisplayLink: item.DisplayLink,
		}
		results = append(results, result)
	}

	searchInfo := SearchInfo{
		TotalResults: rss.Channel.TotalResults,
	}

	if rss.Channel.Title != "" {
		searchInfo.Query = extractQueryFromTitle(rss.Channel.Title)
	}

	return &SearchResult{
		SearchInfo: searchInfo,
		Results:    results,
	}, nil
}

// ParseSearchHTML parses the HTML response from search engines and converts it to SearchResult
func ParseSearchHTML(htmlData []byte) (*SearchResult, error) {
	doc, err := html.Parse(strings.NewReader(string(htmlData)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	results := extractResultsFromHTML(doc)

	return &SearchResult{
		SearchInfo: SearchInfo{},
		Results:    results,
	}, nil
}

// extractResultsFromHTML extracts search results from HTML response
func extractResultsFromHTML(doc *html.Node) []Result {
	var results []Result

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Key == "class" {
					// DuckDuckGo individual result: "result results_links results_links_deep web-result"
					if strings.Contains(attr.Val, "result results_links") {
						result := extractDDGResult(n)
						if result.Title != "" && result.Link != "" {
							results = append(results, result)
						}
						return
					}
					// Bing uses b_algo class (on <li> elements)
					if strings.Contains(attr.Val, "b_algo") {
						result := extractBingResult(n)
						if result.Title != "" && result.Link != "" {
							results = append(results, result)
						}
						return
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(doc)
	return results
}

// extractDDGResult extracts a result from DuckDuckGo HTML
func extractDDGResult(n *html.Node) Result {
	result := Result{}

	var findTitle, findLink, findSnippet func(*html.Node)
	findTitle = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "h2" {
			for _, attr := range node.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, "result__title") {
					var text strings.Builder
					extractTextFromNode(node, &text)
					result.Title = strings.TrimSpace(text.String())
					return
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			findTitle(c)
		}
	}

	findLink = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, "result__a") {
					for _, attr := range node.Attr {
						if attr.Key == "href" {
							// DuckDuckGo uses redirect URLs, extract the actual URL from uddg parameter
							link := extractDDGRedirectURL(attr.Val)
							if link != "" {
								result.Link = link
								return
							}
						}
					}
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			findLink(c)
		}
	}

	findSnippet = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, "result__snippet") {
					var text strings.Builder
					extractTextFromNode(node, &text)
					result.Snippet = strings.TrimSpace(text.String())
					return
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			findSnippet(c)
		}
	}

	findTitle(n)
	findLink(n)
	findSnippet(n)

	return result
}

// extractDDGRedirectURL extracts the actual URL from DuckDuckGo's redirect URL
func extractDDGRedirectURL(link string) string {
	if !strings.Contains(link, "uddg=") {
		return link
	}

	u, err := url.Parse(link)
	if err != nil {
		return link
	}

	uddg := u.Query().Get("uddg")
	if uddg != "" {
		return uddg
	}
	return link
}

// extractBingResult extracts a result from Bing HTML
func extractBingResult(n *html.Node) Result {
	result := Result{}

	var findTitle, findLink, findSnippet func(*html.Node)
	findTitle = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "h2" {
			var text strings.Builder
			extractTextFromNode(node, &text)
			result.Title = strings.TrimSpace(text.String())
			return
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			findTitle(c)
		}
	}

	findLink = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					link := attr.Val
					// Bing uses redirect URLs with /ck/a path, extract the actual URL
					if strings.Contains(link, "/ck/a") {
						link = extractBingRedirectURL(link)
					}
					if strings.HasPrefix(link, "http") && !strings.Contains(link, "bing.com") {
						result.Link = link
						return
					}
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			findLink(c)
		}
	}

	findSnippet = func(node *html.Node) {
		// Bing snippet is in <div class="b_caption"><p class="b_lineclamp2">
		if node.Type == html.ElementNode && node.Data == "div" {
			for _, attr := range node.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, "b_caption") {
					var text strings.Builder
					extractTextFromNode(node, &text)
					trimmed := strings.TrimSpace(text.String())
					if trimmed != "" {
						result.Snippet = trimmed
						return
					}
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			findSnippet(c)
		}
	}

	findTitle(n)
	findLink(n)
	findSnippet(n)

	return result
}

// extractBingRedirectURL extracts the actual URL from Bing's redirect URL
func extractBingRedirectURL(link string) string {
	// First decode HTML entities (amp; -> &)
	link = strings.ReplaceAll(link, "&amp;", "&")

	u, err := url.Parse(link)
	if err != nil {
		return ""
	}

	// Bing uses 'u' parameter for the encoded URL
	encodedURL := u.Query().Get("u")
	if encodedURL != "" {
		// URL is base64 encoded with a1 prefix
		if strings.HasPrefix(encodedURL, "a1") {
			encodedURL = encodedURL[2:]
		}
		// Decode base64
		decoded, err := base64.StdEncoding.DecodeString(encodedURL)
		if err != nil {
			// Try URL decoding as fallback
			decodedStr, err2 := url.QueryUnescape(encodedURL)
			if err2 != nil {
				return ""
			}
			return decodedStr
		}
		return string(decoded)
	}
	return ""
}

// extractTextFromNode extracts text content from a node
func extractTextFromNode(n *html.Node, text *strings.Builder) {
	if n.Type == html.TextNode {
		text.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractTextFromNode(c, text)
	}
}

// extractQueryFromTitle tries to extract the search query from RSS title format
func extractQueryFromTitle(title string) string {
	parts := strings.Split(title, " - ")
	if len(parts) > 1 {
		return parts[1]
	}
	return title
}

// isXML checks if the data is XML format
func isXML(data []byte) bool {
	return strings.Contains(string(data), "<?xml") || strings.HasPrefix(strings.TrimSpace(string(data)), "<rss")
}

// ToJSON converts a SearchResult to JSON bytes
func ToJSON(result *SearchResult) ([]byte, error) {
	return json.MarshalIndent(result, "", "  ")
}

// Helper function to clean up text (remove extra whitespace)
func cleanText(text string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(strings.TrimSpace(text), " ")
}
