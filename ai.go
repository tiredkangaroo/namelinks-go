package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const AI_COMPLETIONS_URL = "https://ai.hackclub.com/chat/completions"

type Message struct {
	Role    string `json:"role"`
	Message string `json:"content"`
}
type Response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		QueueTime        float64 `json:"queue_time"`
		PromptTokens     int     `json:"prompt_tokens"`
		PromptTime       float64 `json:"prompt_time"`
		CompletionTokens int     `json:"completion_tokens"`
		CompletionTime   float64 `json:"completion_time"`
		TotalTokens      int     `json:"total_tokens"`
		TotalTime        float64 `json:"total_time"`
	} `json:"usage"`
	UsageBreakdown    interface{} `json:"usage_breakdown"`
	SystemFingerprint string      `json:"system_fingerprint"`
	XGroq             struct {
		ID string `json:"id"`
	} `json:"x_groq"`
	ServiceTier string `json:"service_tier"`
	Error       string `json:"error"`
}

const SYSTEM_PROMPT = "You are a bot that, when given a company or service name (like \"Google Cloud\", \"GCP\", or \"Fujifilm\"), returns only the official website URL for that service.\n" +
	"Output nothing else—no text, no explanation, no quotes.\n" +
	"If multiple services share the name, give the main official site.\n" +
	"Examples:\n" +
	"Input: \"Google Cloud\" → Output: https://cloud.google.com\n" +
	"Input: \"GCP\" → Output: https://cloud.google.com\n" +
	"Input: \"Fujifilm\" → Output: https://www.fujifilm.com"

func getAILink(name []byte) ([]byte, error) {
	// just throw speed out the window here
	body, err := json.Marshal(map[string]any{
		"messages": []Message{
			{
				Role:    "system",
				Message: SYSTEM_PROMPT,
			},
			{
				Role:    "user",
				Message: string(name),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, AI_COMPLETIONS_URL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respBody Response
	json.NewDecoder(resp.Body).Decode(&respBody)
	if respBody.Error != "" {
		return nil, fmt.Errorf("api error: %s", respBody.Error)
	}
	if len(respBody.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}
	contentParts := strings.Split(respBody.Choices[0].Message.Content, "\n")
	link := contentParts[len(contentParts)-1]
	return []byte(link), nil
}
