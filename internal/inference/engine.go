package inference

import (
"context"
"encoding/json"
"fmt"
"io"
"net/http"
"sync"
"time"
)

type Engine struct {
mu      sync.RWMutex
models  map[string]*Model
config  *Config
client  *http.Client
}

type Model struct {
Name        string    `json:"name"`
Path        string    `json:"path"`
Loaded      bool      `json:"loaded"`
Size        int64     `json:"size"`
Modified    time.Time `json:"modified"`
Digest      string    `json:"digest"`
Details     ModelDetails `json:"details"`
}

type ModelDetails struct {
Format            string   `json:"format"`
Family            string   `json:"family"`
Families          []string `json:"families"`
ParameterSize     string   `json:"parameter_size"`
QuantizationLevel string   `json:"quantization_level"`
}

type Config struct {
OllamaURL   string  `yaml:"ollama_url"`
MaxTokens   int     `yaml:"max_tokens"`
Temperature float64 `yaml:"temperature"`
TopP        float64 `yaml:"top_p"`
Timeout     time.Duration `yaml:"timeout"`
}

type Request struct {
Model   string `json:"model"`
Prompt  string `json:"prompt"`
Stream  bool   `json:"stream"`
System  string `json:"system,omitempty"`
Options map[string]interface{} `json:"options,omitempty"`
}

type Response struct {
Model              string `json:"model"`
CreatedAt          time.Time `json:"created_at"`
Response           string `json:"response"`
Done               bool   `json:"done"`
Context            []int  `json:"context,omitempty"`
TotalDuration      time.Duration `json:"total_duration"`
LoadDuration       time.Duration `json:"load_duration"`
PromptEvalCount    int    `json:"prompt_eval_count"`
PromptEvalDuration time.Duration `json:"prompt_eval_duration"`
EvalCount          int    `json:"eval_count"`
EvalDuration       time.Duration `json:"eval_duration"`
}

func NewEngine() *Engine {
return &Engine{
models: make(map[string]*Model),
config: &Config{
OllamaURL:   "http://localhost:11434",
MaxTokens:   512,
Temperature: 0.7,
TopP:        0.9,
Timeout:     30 * time.Second,
},
client: &http.Client{
Timeout: 30 * time.Second,
},
}
}

func (e *Engine) Process(ctx context.Context, req *Request) (*Response, error) {
ollamaReq := map[string]interface{}{
"model":  req.Model,
"prompt": req.Prompt,
"stream": false,
"options": map[string]interface{}{
"temperature": e.config.Temperature,
"top_p":       e.config.TopP,
"max_tokens":  e.config.MaxTokens,
},
}

jsonData, err := json.Marshal(ollamaReq)
if err != nil {
return nil, fmt.Errorf("failed to marshal request: %w", err)
}

httpReq, err := http.NewRequestWithContext(ctx, "POST", e.config.OllamaURL+"/api/generate", 
io.NopCloser(bytes.NewReader(jsonData)))
if err != nil {
return nil, fmt.Errorf("failed to create request: %w", err)
}
httpReq.Header.Set("Content-Type", "application/json")

resp, err := e.client.Do(httpReq)
if err != nil {
return nil, fmt.Errorf("failed to send request: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
return nil, fmt.Errorf("ollama API error: %s", resp.Status)
}

var response Response
if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
return nil, fmt.Errorf("failed to decode response: %w", err)
}

return &response, nil
}

func (e *Engine) ListModels(ctx context.Context) ([]Model, error) {
req, err := http.NewRequestWithContext(ctx, "GET", e.config.OllamaURL+"/api/tags", nil)
if err != nil {
return nil, fmt.Errorf("failed to create request: %w", err)
}

resp, err := e.client.Do(req)
if err != nil {
return nil, fmt.Errorf("failed to send request: %w", err)
}
defer resp.Body.Close()

var result struct {
Models []Model `json:"models"`
}
if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
return nil, fmt.Errorf("failed to decode response: %w", err)
}

return result.Models, nil
}

func (e *Engine) LoadModel(ctx context.Context, modelName string) error {
loadReq := map[string]string{
"model": modelName,
}

jsonData, err := json.Marshal(loadReq)
if err != nil {
return fmt.Errorf("failed to marshal load request: %w", err)
}

req, err := http.NewRequestWithContext(ctx, "POST", e.config.OllamaURL+"/api/pull",
io.NopCloser(bytes.NewReader(jsonData)))
if err != nil {
return fmt.Errorf("failed to create load request: %w", err)
}
req.Header.Set("Content-Type", "application/json")

resp, err := e.client.Do(req)
if err != nil {
return fmt.Errorf("failed to load model: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
return fmt.Errorf("failed to load model: %s", resp.Status)
}

return nil
}
