package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
)

func executeRequest[T any](logger *slog.Logger, method, endpoint string, payload interface{}) T {
	var bodyBuffer *bytes.Buffer
	if payload != nil && method != "GET" {
		data, err := json.Marshal(payload)
		if err != nil {
			logger.Error("Error marshalling payload", "error", err)
			panic(err)
		}
		bodyBuffer = bytes.NewBuffer(data)
	} else {
		bodyBuffer = bytes.NewBuffer(nil)
	}

	if token == "" && endpoint != "auth/login" {
		logger.Error("Authorization token is not set. Please ensure login() is called before making requests.")
		panic("Authorization token is not set")
	}

	if method == "" {
		logger.Error("HTTP method is empty. Please provide a valid method.")
		panic("HTTP method is empty")
	}
	if endpoint == "" {
		logger.Error("Endpoint is empty. Please provide a valid endpoint.")
		panic("Endpoint is empty")
	}

	logger.Debug("Preparing request", "method", method, "endpoint", endpoint)

	fullURL := fmt.Sprintf("%s/%s", baseURL, endpoint)

	req, err := http.NewRequest(method, fullURL, bodyBuffer)
	if err != nil {
		logger.Error("Error creating request", "method", method, "error", err)
		panic(err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	logger.Info("Executing HTTP request")

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error executing request", "method", method, "error", err)
		panic(err)
	}
	defer resp.Body.Close()

	logger.Info("Response received", "status_code", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		logger.Error("Non-OK HTTP status", "status", resp.Status)
		panic(fmt.Sprintf("Non-OK HTTP status: %s", resp.Status))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error reading response body", "error", err)
		panic(err)
	}

	logger.Debug("Raw response body", "body", string(body))

	var result T
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Error("Error unmarshalling response body", "error", err)
		panic(err)
	}

	logger.Debug("Parsed response", "result", result)
	return result
}

func makeGetRequest[T any](logger *slog.Logger, endpoint string) T {
	return executeRequest[T](logger, "GET", endpoint, nil)
}

func makePostRequest[T any](logger *slog.Logger, endpoint string, payload interface{}) T {
	return executeRequest[T](logger, "POST", endpoint, payload)
}
