package redditclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client interface {
	GetUserComments(ctx context.Context, username string, limit int) (*RedditListing, error)
}

type RedditListing struct {
	Data struct {
		Children []struct {
			Data struct {
				Subreddit string `json:"subreddit"`
				Body      string `json:"body"`
				Permalink string `json:"permalink"`
				Score     int    `json:"score"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type httpRedditClient struct {
	client *http.Client
}

func NewClient() Client {
	return &httpRedditClient{client: &http.Client{Transport: http.DefaultTransport}}
}

func (h *httpRedditClient) GetUserComments(ctx context.Context, username string, limit int) (*RedditListing, error) {
	url := fmt.Sprintf("https://www.reddit.com/user/%s/comments.json?limit=%d", username, limit)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Reddit requires a user-agent header
	req.Header.Set("User-Agent", "golang:reddit.comment.fetcher:v1.0 (by /u/yourusername)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var comments RedditListing
	if err := json.Unmarshal(body, &comments); err != nil {
		return nil, err
	}

	return &comments, nil
}
