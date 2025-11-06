package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type CheckResponse struct {
	CorrectedText string `json:"corrected_text"`
	HasChanges    bool   `json:"has_changes"`
	Explanation   string `json:"explanation"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.deepseek.com/v1",
	}
}

func (c *Client) CheckSpellingAndPunctuation(ctx context.Context, text string) (*CheckResponse, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	prompt := fmt.Sprintf(`Ты - эксперт по русской орфографии и пунктуации. Проверь текст на ошибки и исправь их, сохранив исходный смысл и стиль. Верни ТОЛЬКО валидный JSON без дополнительных комментариев.

Формат ответа:
{
  "corrected_text": "исправленный текст",
  "has_changes": true/false,
  "explanation": "краткое объяснение сделанных исправлений или пустая строка если изменений нет"
}

Важно:
- Исправь ВСЕ орфографические, пунктуационные и грамматические ошибки
- Сохрани исходный смысл, тон и стиль текста
- Если ошибок нет, верни исходный текст в corrected_text и has_changes: false
- В explanation кратко опиши что было исправлено

Текст: "%s"`, text)

	requestBody := ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("API error: %s", errorResp.Error.Message)
	}

	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices received")
	}

	responseContent := chatResp.Choices[0].Message.Content

	// Extract JSON from markdown code block if present
	jsonContent := extractJSONFromResponse(responseContent)

	// Parse the JSON response from the AI
	var checkResp CheckResponse
	if err := json.Unmarshal([]byte(jsonContent), &checkResp); err != nil {
		return nil, fmt.Errorf("failed to parse AI response as JSON: %w, content: %s", err, jsonContent)
	}

	return &checkResp, nil
}

func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

// extractJSONFromResponse extracts JSON content from markdown code blocks
func extractJSONFromResponse(content string) string {
	// Remove markdown code block markers
	re := regexp.MustCompile(`(?s)` + "```" + `(?:json)?\s*([\s\S]*?)\s*` + "```")
	matches := re.FindStringSubmatch(content)
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1])
	}
	
	// If no code block found, try to find JSON object boundaries
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "{") && strings.HasSuffix(content, "}") {
		return content
	}
	
	// Look for JSON object within the text
	startIdx := strings.Index(content, "{")
	endIdx := strings.LastIndex(content, "}")
	if startIdx >= 0 && endIdx > startIdx {
		return content[startIdx:endIdx+1]
	}
	
	return content
}
