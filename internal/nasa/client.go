package nasa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type APOD struct {
	Date        string `json:"date"`
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	MediaType   string `json:"media_type"`
	URL         string `json:"url"`
	HDURL       string `json:"hdurl"`
	Copyright   string `json:"copyright"`
}

const apiURL = "https://api.nasa.gov/planetary/apod"

func Fetch(apiKey string) (*APOD, error) {

	params := url.Values{}
	params.Set("api_key", apiKey)
	params.Set("date", "2026-11-01") // ВРЕМЕННО, потом уберём
	reqURL := apiURL + "?" + params.Encode()

	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("request apod: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("apod: unexpected status %d", resp.StatusCode)
	}

	var apod APOD
	if err := json.NewDecoder(resp.Body).Decode(&apod); err != nil {
		return nil, fmt.Errorf("apod: failed to decode response: %w", err)
	}

	return &apod, nil
}
