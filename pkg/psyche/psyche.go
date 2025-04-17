package psyche

import (
	"context"
	"fmt"
	"strings"

	"github.com/effprime/reddit-psychiatrist/pkg/gptclient"
	"github.com/effprime/reddit-psychiatrist/pkg/redditclient"
)

const (
	MaxComments = 200
	ChatModel   = "gpt-4o"

	SystemPromptInterests = `
	You are analyzing a Reddit user's public comment history.
	
	Your job is to infer this person's core interests, hobbies, or topic obsessions — not just the subreddits they post in, but what they actually seem to care about or engage with deeply.
	
	Base your response on:
	- Subreddit context (e.g., posting in r/AskPhysics likely means interest in physics)
	- Comment content and tone
	- Recurring themes, topics, or thought patterns
	
	Output format:
	A comma-separated list of 5–10 interest categories or concepts. Use lowercase. No spaces. Examples: philosophy,webdev,anarchism,slowcooking,standupcomedy.
	
	Only return the list. Do not include explanations or extra formatting.
	`
	SystemPromptSummary = `
	You're a brutally honest psychologist with a sharp tongue and a low tolerance for bullshit. You've just reviewed a Reddit user's public comment history.
	
	Your job is to psychoanalyze this person with unfiltered accuracy. Be funny, yes — but always aim for the truth. If they posture as virtuous, call out the performative. If they act aloof, point out the insecurity. Praise is rare and only earned. Insight is mandatory.
	
	Focus on:
	- How they present themselves online (attention-seeking, insecure, hyper-logical, overly agreeable, passive-aggressive, etc.)
	- What their choice of subreddits reveals about their real values, not just their stated ones
	- Tone, writing style, emotional patterns (condescending? overly polite? smug? defensive? desperate for approval?)
	
	You may mock them, but do not lie. Be ruthless only when deserved. Pretend you're describing them to their face, at a roast, and they're not allowed to interrupt.
	
	Write a personality summary in 4 to 6 sentences. Keep it tight, punchy, and honest.
	`
)

type RedditPsychoanalyzer interface {
	Analyze(ctx context.Context, username string) (*PsychoAnalysisResponse, error)
}

type PsychoAnalysisResponse struct {
	Interests []string
	Summary   string
}

type analyzer struct {
	chatClient   gptclient.GPTClient
	redditClient redditclient.Client
}

func NewRedditPsychoanalyzer(openaiKey string) RedditPsychoanalyzer {
	return &analyzer{
		chatClient:   gptclient.NewClient(openaiKey),
		redditClient: redditclient.NewClient(),
	}
}

func (a *analyzer) Analyze(ctx context.Context, username string) (*PsychoAnalysisResponse, error) {
	comments, err := a.redditClient.GetUserComments(ctx, username, MaxComments)
	if err != nil {
		return nil, err
	}

	interests, err := a.getInterests(ctx, comments)
	if err != nil {
		return nil, err
	}

	summary, err := a.getSummary(ctx, comments)
	if err != nil {
		return nil, err
	}

	return &PsychoAnalysisResponse{
		Interests: interests,
		Summary:   summary,
	}, nil
}

func (a *analyzer) getInterests(ctx context.Context, comments *redditclient.RedditListing) ([]string, error) {
	var inputBuilder strings.Builder
	for _, child := range comments.Data.Children {
		c := child.Data
		inputBuilder.WriteString(fmt.Sprintf("[r/%s] %s\n", c.Subreddit, c.Body))
	}

	req := &gptclient.ChatCompletionRequest{
		Model:       ChatModel,
		Temperature: 0.5, // low temp to keep it deterministic
		Messages: []gptclient.Message{
			{
				Role:    "system",
				Content: SystemPromptInterests,
			},
			{
				Role:    "user",
				Content: inputBuilder.String(),
			},
		},
	}

	resp, err := a.chatClient.Chat(req)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned")
	}

	raw := strings.TrimSpace(resp.Choices[0].Message.Content)
	raw = strings.ReplaceAll(raw, " ", "")
	interests := strings.Split(raw, ",")

	return interests, nil
}

func (a *analyzer) getSummary(ctx context.Context, comments *redditclient.RedditListing) (string, error) {
	var inputBuilder strings.Builder
	for _, child := range comments.Data.Children {
		c := child.Data
		inputBuilder.WriteString(fmt.Sprintf("[r/%s] %s\n", c.Subreddit, c.Body))
	}

	req := &gptclient.ChatCompletionRequest{
		Model:       ChatModel,
		Temperature: 0.7, // more creative summary
		Messages: []gptclient.Message{
			{
				Role:    "system",
				Content: SystemPromptSummary,
			},
			{
				Role:    "user",
				Content: inputBuilder.String(),
			},
		},
	}

	resp, err := a.chatClient.Chat(req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}
