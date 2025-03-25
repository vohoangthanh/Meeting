package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CloudflareService struct {
	baseURL   string
	appID     string
	appSecret string
}

type CloudflareResponse struct {
	ErrorCode        *string                `json:"errorCode"`
	ErrorDescription string                 `json:"errorDescription"`
	SessionID        string                 `json:"sessionId"`
	Result           map[string]interface{} `json:"result"`
}

// NewCloudflareService creates a new Cloudflare service instance
func NewCloudflareService(baseURL, appID, appSecret string) *CloudflareService {
	return &CloudflareService{
		baseURL:   baseURL,
		appID:     appID,
		appSecret: appSecret,
	}
}

// makeCloudflareRequest is a helper function to make requests to Cloudflare API
func (cs *CloudflareService) makeCloudflareRequest(method string, path string, body interface{}) (map[string]interface{}, error) {
	var reqBodyBytes []byte
	var err error

	if body != nil {
		reqBodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %v", err)
		}
	}

	url := fmt.Sprintf("%s/%s%s", cs.baseURL, cs.appID, path)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+cs.appSecret)
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var cfResp CloudflareResponse
	if err := json.Unmarshal(respBody, &cfResp); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	// Check for Cloudflare error
	if cfResp.ErrorCode != nil {
		return nil, fmt.Errorf("Cloudflare API error: %s", cfResp.ErrorDescription)
	}

	// For session creation, return a map with sessionId
	if cfResp.SessionID != "" {
		return map[string]interface{}{
			"sessionId": cfResp.SessionID,
		}, nil
	}

	// Return the result or whole response as map
	if len(cfResp.Result) > 0 {
		return cfResp.Result, nil
	}

	// Convert the whole response to map as fallback
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error converting response to map: %v", err)
	}

	return result, nil
}

// CreateSession creates a new session
func (cs *CloudflareService) CreateSession() (string, error) {
	response, err := cs.makeCloudflareRequest("POST", "/sessions/new", nil)
	if err != nil {
		return "", err
	}

	sessionID, ok := response["sessionId"].(string)
	if !ok {
		return "", fmt.Errorf("sessionId not found in response")
	}

	return sessionID, nil
}

// PublishTracks publishes tracks to a session
func (cs *CloudflareService) PublishTracks(sessionID string, offer map[string]interface{}, tracks []map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("/sessions/%s/tracks/new", sessionID)

	requestBody := map[string]interface{}{
		"sessionDescription": offer,
		"tracks":             tracks,
	}

	return cs.makeCloudflareRequest("POST", url, requestBody)
}

// PullTracks pulls tracks from a remote session
func (cs *CloudflareService) PullTracks(sessionID string, remoteTracks []map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("/sessions/%s/tracks/new", sessionID)

	requestBody := map[string]interface{}{
		"tracks": remoteTracks,
	}

	return cs.makeCloudflareRequest("POST", url, requestBody)
}

// Renegotiate performs renegotiation for a session
func (cs *CloudflareService) Renegotiate(sessionID string, sessionDescription map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("/sessions/%s/renegotiate", sessionID)

	requestBody := map[string]interface{}{
		"sessionDescription": sessionDescription,
	}

	return cs.makeCloudflareRequest("PUT", url, requestBody)
}

// CloseTracks closes tracks in a session
func (cs *CloudflareService) CloseTracks(sessionID string, tracks []map[string]string, force bool, sessionDescription map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("/sessions/%s/tracks/close", sessionID)

	requestBody := map[string]interface{}{
		"tracks":             tracks,
		"force":              force,
		"sessionDescription": sessionDescription,
	}

	return cs.makeCloudflareRequest("PUT", url, requestBody)
}

// GetSessionState gets the current state of a session
func (cs *CloudflareService) GetSessionState(sessionID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("/sessions/%s", sessionID)
	return cs.makeCloudflareRequest("GET", url, nil)
}
