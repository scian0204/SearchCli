package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Define CLI flags
	query := flag.String("q", "", "Search query")
	engine := flag.String("engine", "ddg", "Search engine: ddg (DuckDuckGo), bing (Bing)")
	crawl := flag.Bool("crawl", false, "Enable crawling of search result links")
	maxLinks := flag.Int("max-links", 5, "Maximum number of links to crawl (default: 5)")
	output := flag.String("output", "", "Output file path (if empty, prints to stdout)")
	help := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if *help || *query == "" {
		printUsage()
		if *help {
			os.Exit(0)
		}
		os.Exit(1)
	}

	// Fetch and parse search results
	fmt.Println("Fetching search results...")
	searchResult, err := FetchSearch(*query, *engine)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d results\n", len(searchResult.Results))

	// Crawl links if requested
	if *crawl {
		fmt.Printf("Crawling up to %d links...\n", *maxLinks)
		CrawlResults(searchResult, *maxLinks)
	}

	// Convert to JSON
	jsonData, err := ToJSON(searchResult)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting to JSON: %v\n", err)
		os.Exit(1)
	}

	// Output results
	if *output != "" {
		if err := os.WriteFile(*output, jsonData, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Results saved to %s\n", *output)
	} else {
		fmt.Println(string(jsonData))
	}
}

func printUsage() {
	fmt.Println(`Search CLI Tool

Usage:
  searchcli [options]

Options:
  -q string           Search query (required)
  -engine string      Search engine: ddg (DuckDuckGo), bing (Bing) (default: ddg)
  -crawl              Enable crawling of search result links
  -max-links int      Maximum number of links to crawl (default: 5)
  -output string      Output file path (if empty, prints to stdout)
  -help               Show this help message

Examples:
  # Search with DuckDuckGo
  ./searchcli -q "python programming"

  # Search with link crawling
  ./searchcli -q "python programming" -crawl -max-links 5

  # Save results to file
  ./searchcli -q "python programming" -crawl -output results.json

  # Use Bing
  ./searchcli -q "python programming" -engine bing`)
}
