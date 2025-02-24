package coprocessor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

var ROLLUP_HTTP_SERVER_URL = os.Getenv("ROLLUP_HTTP_SERVER_URL")

func SendPost(endpoint string, jsonData []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, ROLLUP_HTTP_SERVER_URL+"/"+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	return http.DefaultClient.Do(req)
}

func SendFinish(finish *FinishRequest) (*http.Response, error) {
	body, err := json.Marshal(finish)
	if err != nil {
		return &http.Response{}, err
	}

	return SendPost("finish", body)
}

func SendNotice(notice *NoticeRequest) (*http.Response, error) {
	body, err := json.Marshal(notice)
	if err != nil {
		return &http.Response{}, err
	}

	return SendPost("notice", body)
}

func SendException(exception *ExceptionRequest) (*http.Response, error) {
	body, err := json.Marshal(exception)
	if err != nil {
		return &http.Response{}, err
	}

	return SendPost("exception", body)
}
