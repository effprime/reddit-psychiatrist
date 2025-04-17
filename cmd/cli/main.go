package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/effprime/reddit-psychiatrist/pkg/psyche"
)

func main() {
	username := flag.String("u", "", "Reddit username to analyze (required)")
	apiKey := flag.String("k", "", "OpenAI API key (can also set OPENAI_API_KEY env var)")
	timeout := flag.Int("t", 60, "Timeout in seconds for the analysis")
	flag.Parse()

	if *username == "" {
		fmt.Println("Error: -u <username> is required")
		os.Exit(1)
	}

	key := *apiKey
	if key == "" {
		key = os.Getenv("OPENAI_API_KEY")
		if key == "" {
			fmt.Println("Error: OpenAI API key not provided. Use -k or set OPENAI_API_KEY env var.")
			os.Exit(1)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	analyzer := psyche.NewRedditPsychoanalyzer(key)

	fmt.Printf("Analyzing Reddit user: u/%s\n", *username)
	resp, err := analyzer.Analyze(ctx, *username)
	if err != nil {
		log.Fatalf("Analysis failed: %v", err)
	}

	fmt.Println("\n🧠 Personality Summary:")
	fmt.Println(resp.Summary)

	fmt.Println("\n🔥 Inferred Interests:")
	for _, interest := range resp.Interests {
		fmt.Printf("- %s\n", interest)
	}
}
