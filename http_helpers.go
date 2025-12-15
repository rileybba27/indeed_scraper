package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

const USER_AGENT string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36"

func NewGetRequest(url string, cookie string) (*http.Request, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	request.Header.Set("Accept-Language", "en-US,en;q=0.9")
	request.Header.Set("Cookie", cookie)
	request.Header.Set("User-Agent", USER_AGENT)
	return request, nil
}

func GetURL(url string, cookie string, httpClient *http.Client) (string, error) {
	if httpClient == nil {
		return "", errors.New("Missing valid HTTP client pointer in GetURL parameters.")
	}

	request, err := NewGetRequest(url, cookie)
	if err != nil {
		return "", err
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Failed to get 200 OK Status from server when requesting page, Status: %s", response.Status)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
