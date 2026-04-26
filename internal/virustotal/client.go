package virustotal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type Client struct {
	apiKey string
	client *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) GetReportByHash(hash string) (*VTResponse, bool, error) {
	url := fmt.Sprintf("https://www.virustotal.com/api/v3/files/%s", hash)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, false, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("x-apikey", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, false, nil
	}

	if resp.StatusCode != 200 {
		return nil, false, fmt.Errorf("status %d", resp.StatusCode)
	}

	var result VTResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, false, err
	}

	return &result, true, nil
}

func (c *Client) UploadFile(filename string, content []byte) (string, error) {
	url := "https://www.virustotal.com/api/v3/files"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(content); err != nil {
		return "", fmt.Errorf("writing file content: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("closing multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("x-apikey", c.apiKey)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("upload status %d", resp.StatusCode)
	}

	var result VTResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Data.ID, nil
}

func (c *Client) PollAnalysis(analysisID string) (*VTResponse, error) {
	url := fmt.Sprintf("https://www.virustotal.com/api/v3/analyses/%s", analysisID)

	// Чекаємо до 5 хвилин
	for i := 0; i < 20; i++ {
		time.Sleep(15 * time.Second) // Повага до rate limit навіть всередині поллінгу

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}
		req.Header.Add("x-apikey", c.apiKey)

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, err
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return nil, fmt.Errorf("reading response body: %w", err)
		}

		var result VTResponse
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			return nil, fmt.Errorf("decoding response: %w", err)
		}

		if result.Data.Attributes.Status == "completed" {
			return &result, nil
		}
	}
	return nil, fmt.Errorf("timeout")
}
