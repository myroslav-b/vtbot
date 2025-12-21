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
		client: &http.Client{},
	}
}

func (c *Client) GetReportByHash(hash string) (*VTResponse, bool, error) {
	url := fmt.Sprintf("https://www.virustotal.com/api/v3/files/%s", hash)
	req, _ := http.NewRequest("GET", url, nil)
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
	part.Write(content)
	writer.Close()

	req, _ := http.NewRequest("POST", url, body)
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

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("x-apikey", c.apiKey)

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var result VTResponse
		bodyBytes, _ := io.ReadAll(resp.Body)
		json.Unmarshal(bodyBytes, &result)

		if result.Data.Attributes.Status == "completed" {
			return &result, nil
		}
	}
	return nil, fmt.Errorf("timeout")
}
