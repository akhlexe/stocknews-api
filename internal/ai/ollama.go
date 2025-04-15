package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func GenerateSummary(prompt string) (string, error) {

	ollamaUrl := GetOllamaUrl()

	log.Println("🧠 Calling Ollama for summary...")

	payload := OllamaRequest{
		Model:  "llama3",
		Prompt: prompt,
		Stream: false,
	}

	body, _ := json.Marshal(payload)
	log.Printf("📤 Prompt: %s\n", prompt)

	resp, err := http.Post(ollamaUrl+"/api/generate", "application/json", bytes.NewReader(body))
	if err != nil {
		log.Println("❌ Ollama request failed:", err)
		return "", err
	}

	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	log.Println("✅ Ollama responded")

	var res OllamaResponse
	err = json.Unmarshal(raw, &res)

	if err != nil {
		log.Printf("⚠️ Failed to parse Ollama response: %v\nRaw: %s\n", err, string(raw))
		return "", fmt.Errorf("ollama error: %v - raw: %s", err, raw)
	}

	log.Println("📝 AI Summary done")
	return res.Response, nil
}

func SummarizeArticles(ctx context.Context, combined string) (string, error) {
	ollamaUrl := GetOllamaUrl()

	payload := OllamaRequest{
		Model:  "llama3",
		Prompt: fmt.Sprintf("Summarize the following stock market news: \n\n%s", combined),
		Stream: false,
	}

	data, _ := json.Marshal(payload)

	resp, err := http.Post(ollamaUrl+"/api/generate", "application/json", bytes.NewBuffer(data))

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var result OllamaResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Response, nil
}

func GetOllamaUrl() string {
	return os.Getenv("OLLAMA_URL")
}
