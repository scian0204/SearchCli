package main

import "encoding/xml"

// SearchResult represents the top-level structure of search results
type SearchResult struct {
	SearchInfo SearchInfo `json:"search_info"`
	Results    []Result   `json:"results"`
}

// SearchInfo contains metadata about the search
type SearchInfo struct {
	TotalResults     string `json:"total_results,omitempty"`
	SearchTime       string `json:"search_time,omitempty"`
	FormattedResults string `json:"formatted_results,omitempty"`
	FormattedTime    string `json:"formatted_time,omitempty"`
	Query            string `json:"query,omitempty"`
}

// Result represents a single search result
type Result struct {
	Title       string          `json:"title"`
	Link        string          `json:"link"`
	Snippet     string          `json:"snippet,omitempty"`
	DisplayLink string          `json:"display_link,omitempty"`
	CrawledContent *CrawledContent `json:"crawled_content,omitempty"`
}

// CrawledContent contains the extracted content from a crawled URL
type CrawledContent struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Keywords    string `json:"keywords,omitempty"`
	Headings    []string `json:"headings,omitempty"`
	Paragraphs  []string `json:"paragraphs,omitempty"`
	Links       []string `json:"links,omitempty"`
	SourceURL   string `json:"source_url"`
}

// GoogleSearchXML represents the XML structure returned by Google Custom Search
type GoogleSearchXML struct {
	XMLName    xml.Name `xml:"rss"`
	Channel    Channel  `xml:"channel"`
}

// Channel represents the RSS channel in Google search results
type Channel struct {
	Title       string       `xml:"title"`
	Link        string       `xml:"link"`
	Description string       `xml:"description"`
	TotalResults string      `xml:"totalResults,attr,omitempty"`
	Items       []Item       `xml:"item"`
}

// Item represents a single search result item in the XML
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	DisplayLink string `xml:"displayLink,omitempty"`
}
