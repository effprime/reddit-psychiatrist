package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/effprime/reddit-psychiatrist/pkg/psyche"
)

type PageData struct {
	Username  string
	Interests []string
	Summary   string
	Error     string
}

var tmpl *template.Template

func main() {
	var (
		addr    = flag.String("addr", ":8080", "HTTP server listen address")
		openai  = flag.String("openai-key", "", "OpenAI API key (optional, will fallback to OPENAI_API_KEY env var)")
		timeout = flag.Int("timeout", 60, "Request timeout in seconds")
	)
	flag.Parse()

	apiKey := *openai
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			log.Fatal("Missing OpenAI API key. Use -openai-key or set OPENAI_API_KEY env var.")
		}
	}

	analyzer := psyche.NewRedditPsychoanalyzer(apiKey)

	var err error
	tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("Template parse error: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl.Execute(w, nil)
			return
		}

		if err := r.ParseForm(); err != nil {
			renderError(w, "Invalid form submission")
			return
		}

		username := r.FormValue("username")
		if username == "" {
			renderError(w, "Username is required")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(*timeout)*time.Second)
		defer cancel()

		result, err := analyzer.Analyze(ctx, username)
		if err != nil {
			renderError(w, err.Error())
			return
		}

		data := PageData{
			Username:  username,
			Interests: result.Interests,
			Summary:   result.Summary,
		}
		tmpl.Execute(w, data)
	})

	log.Printf("Serving on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func renderError(w http.ResponseWriter, msg string) {
	tmpl.Execute(w, PageData{Error: msg})
}
