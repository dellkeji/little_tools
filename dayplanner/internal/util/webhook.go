package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// WeChatTextMessage represents a WeChat Work text message
type WeChatTextMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

// SendWeChatTextMessage sends a text message to WeChat Work webhook
// Returns error if the request fails or times out
func SendWeChatTextMessage(webhookURL, content string) error {
	// Create message payload
	message := WeChatTextMessage{
		MsgType: "text",
	}
	message.Text.Content = content

	// Marshal to JSON
	payload, err := json.Marshal(message)
	if err != nil {
		log.Printf("ERROR: Failed to marshal webhook message: %v", err)
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create POST request
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("ERROR: Failed to create webhook request: %v", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: Failed to send webhook request: %v", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("WARN: Failed to read webhook response body: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR: Webhook request failed with status %d: %s", resp.StatusCode, string(body))
		return fmt.Errorf("webhook request failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("INFO: Webhook message sent successfully")
	return nil
}
