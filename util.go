package main

import (
	"io"
	"net/http"
)

func queryHTML(url string) (string, error) {
	// Perform the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	// Ensure the response body is closed after the function returns
	defer resp.Body.Close()

	// Check if the request was successful (HTTP status code 200)
	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Print the page content (the HTML body)
	return string(body), nil
}
