package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/effprime/reddit-psychiatrist/pkg/psyche"
)

type requestPayload struct {
	Username string `json:"username"`
}

type responsePayload struct {
	Username  string   `json:"username"`
	Interests []string `json:"interests"`
	Summary   string   `json:"summary"`
	Error     string   `json:"error,omitempty"`
}

func main() {
	var (
		addr    = flag.String("addr", ":8080", "HTTP listen address")
		openai  = flag.String("openai-key", "", "OpenAI API Key (optional, will fallback to env OPENAI_API_KEY)")
		timeout = flag.Int("timeout", 60, "Request timeout in seconds")
	)
	flag.Parse()

	key := *openai
	if key == "" {
		key = os.Getenv("OPENAI_API_KEY")
		if key == "" {
			log.Fatal("Missing OpenAI API key. Use -openai-key or set OPENAI_API_KEY env var.")
		}
	}

	analyzer := psyche.NewRedditPsychoanalyzer(key)

	http.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST supported", http.StatusMethodNotAllowed)
			return
		}

		var req requestPayload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" {
			http.Error(w, "Missing or invalid 'username' in request body", http.StatusBadRequest)
			return
		}

		log.Printf("Analyzing u/%s", req.Username)
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(*timeout)*time.Second)
		defer cancel()

		result, err := analyzer.Analyze(ctx, req.Username)
		resp := responsePayload{Username: req.Username}

		if err != nil {
			log.Printf("Error: %v", err)
			resp.Error = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			resp.Interests = result.Interests
			resp.Summary = result.Summary
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Printf("Listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
