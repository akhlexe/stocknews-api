package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	log.Println("🧠 Calling Ollama for summary...")

	payload := OllamaRequest{
		Model:  "llama3",
		Prompt: prompt,
		Stream: false,
	}

	body, _ := json.Marshal(payload)
	log.Printf("📤 Prompt: %s\n", prompt)

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewReader(body))
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
