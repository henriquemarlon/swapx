package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var ROLLUP_HTTP_SERVER_URL = os.Getenv("ROLLUP_HTTP_SERVER_URL")

type GioRequest struct {
	Domain uint16 `json:"domain"`
	Id     string `json:"id"`
}

type GioResponse struct {
	ResponseCode uint16 `json:"response_code"`
	Response     string `json:"response"`
}

func GetListDataStorageAt(req *GioRequest) (*GioResponse, error) {

	url := ROLLUP_HTTP_SERVER_URL + "/gio"

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request JSON: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d - %s", resp.StatusCode, string(body))
	}

	var response *GioResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return response, nil
}
