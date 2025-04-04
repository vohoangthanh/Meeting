package handle

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
	// log result
	log.Printf("[Cloudflare API Response] Status: %d, Body: %s", resp.StatusCode, string(respBody))
	return result, nil
}

// CreateSession creates a new session
func (cs *CloudflareService) CreateSession() (string, error) {
	cloudflareBasePath := fmt.Sprintf("%s/%s", "https://rtc.live.cloudflare.com/v1/apps", cs.appID)

	url := fmt.Sprintf("%s/sessions/new", cloudflareBasePath)
	log.Printf("[Cloudflare API] Creating new session: %s", url)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("[Cloudflare API Error] Failed to create request: %v", err)
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+cs.appSecret)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Cloudflare API Error] Failed to execute request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Cloudflare API Error] Failed to read response body: %v", err)
		return "", err
	}

	log.Printf("[Cloudflare API Response] Status: %d, Body: %s", resp.StatusCode, string(body))

	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Printf("[Cloudflare API Error] Failed to parse response: %v", err)
		return "", err
	}

	sessionID, ok := responseData["sessionId"].(string)
	if !ok {
		log.Printf("[Cloudflare API Error] Session ID not found in response")
		return "", fmt.Errorf("sessionId not found in response: %s", string(body))
	}

	log.Printf("[Cloudflare API Success] Created session: %s", sessionID)
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
func (cs *CloudflareService) PullTracks(sessionID string, tracks []map[string]interface{}) (map[string]interface{}, error) {
	log.Printf("[Cloudflare API] Pulling tracks for session %s", sessionID)

	// // First check session state
	// state, err := cs.GetSessionState(sessionID)
	// if err != nil {
	// 	return nil, fmt.Errorf("error getting session state: %v", err)
	// }

	// // Log session state
	// log.Printf("[Cloudflare API] Session %s state: %+v", sessionID, state)

	// // Check if session is connected
	// if state == nil || state["connectionState"] != "connected" {
	// 	return nil, fmt.Errorf("Session %s is not in connected state. Current state: %v",
	// 		sessionID, state["connectionState"])
	// }

	// Validate tracks format and prepare request
	for _, track := range tracks {
		if _, ok := track["trackName"]; !ok {
			return nil, fmt.Errorf("trackName must be present in track data")
		}
		if _, ok := track["location"]; !ok {
			track["location"] = "remote"
		}
		if _, ok := track["sessionId"]; !ok {
			return nil, fmt.Errorf("sessionId must be present when pulling remote tracks")
		}
	}

	url := fmt.Sprintf("%s/%s/sessions/%s/tracks/new", cs.baseURL, cs.appID, sessionID)
	log.Printf("[Cloudflare API] Making request to: %s", url)

	// Create request body
	requestBody := map[string]interface{}{
		"tracks": tracks,
	}

	// Log request body
	requestBytes, _ := json.Marshal(requestBody)
	log.Printf("[Cloudflare API] Request body: %s", string(requestBytes))

	// Make API call
	response, err := cs.makeCloudflareRequest("POST", fmt.Sprintf("/sessions/%s/tracks/new", sessionID), requestBody)
	if err != nil {
		log.Printf("[Cloudflare API Error] Failed to pull tracks: %v", err)
		return nil, fmt.Errorf("failed to pull tracks: %v", err)
	}

	// Log successful response
	log.Printf("[Cloudflare API Success] Pulled tracks successfully for session %s", sessionID)

	return response, nil
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

// decompress decompresses zlib-compressed base64-encoded data
func (cs *CloudflareService) decompressZlib(compressedB64 string) (string, error) {
	// Decode base64
	compressedData, err := base64.StdEncoding.DecodeString(compressedB64)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	// Create zlib reader
	zlibReader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return "", fmt.Errorf("failed to create zlib reader: %v", err)
	}
	defer zlibReader.Close()

	// Read decompressed data
	decompressed, err := io.ReadAll(zlibReader)
	if err != nil {
		return "", fmt.Errorf("failed to decompress data: %v", err)
	}

	return string(decompressed), nil
}

// decompressData is a more general decompression function that can handle both
// base64+zlib compression and direct zlib compression
func (cs *CloudflareService) decompressData(data []byte) (string, error) {
	// Check if the data is in string format
	dataStr := string(data)

	// If data starts with "zlib:" prefix, it's base64-encoded zlib data from our frontend
	if strings.HasPrefix(dataStr, "zlib:") {
		// Extract the base64 part
		compressedB64 := strings.TrimPrefix(dataStr, "zlib:")
		log.Printf("Found zlib-prefixed data, decompressing with base64 decoder")
		return cs.decompressZlib(compressedB64)
	}

	// Try to decompress directly as zlib
	zlibReader, err := zlib.NewReader(bytes.NewReader(data))
	if err == nil {
		defer zlibReader.Close()

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, zlibReader); err != nil {
			return "", fmt.Errorf("failed to decompress zlib data: %v", err)
		}

		log.Printf("Successfully decompressed data directly with zlib")
		return buf.String(), nil
	}

	// If we got here, try treating the data as JSON first
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err == nil {
		// The data is valid JSON, check if it has a compressedData field
		if compressedStr, ok := result["compressedData"].(string); ok && strings.HasPrefix(compressedStr, "zlib:") {
			// Extract the base64 part
			compressedB64 := strings.TrimPrefix(compressedStr, "zlib:")
			log.Printf("Found nested zlib-prefixed data in JSON, decompressing with base64 decoder")
			return cs.decompressZlib(compressedB64)
		}
	}

	// If no compression method worked, return data as-is
	log.Printf("No compression format detected, returning data as-is")
	return dataStr, nil
}

// HandleParticipantEvent processes events for a participant, creates a session and publishes tracks
func (cs *CloudflareService) HandleParticipantEvent(roomID string, participant string, tracks []interface{}, sessionDescription interface{}) (string, map[string]interface{}, error) {
	// Create session
	sessionID, err := cs.CreateSession()
	if err != nil {
		return "", nil, fmt.Errorf("failed to create Cloudflare session: %v", err)
	}

	// Format session description for Cloudflare
	var offer map[string]interface{}

	if sdpStr, ok := sessionDescription.([]byte); ok {
		// Parse SDP from bytes
		sdpString := string(sdpStr)

		// Check if the SDP is zlib compressed
		if strings.HasPrefix(sdpString, "zlib:") {
			// Extract the base64 part
			compressedB64 := strings.TrimPrefix(sdpString, "zlib:")

			// Decompress the SDP
			decompressedJson, err := cs.decompressZlib(compressedB64)
			if err != nil {
				return sessionID, nil, fmt.Errorf("failed to decompress SDP: %v", err)
			}

			// Parse decompressed SDP
			if err := json.Unmarshal([]byte(decompressedJson), &offer); err != nil {
				return sessionID, nil, fmt.Errorf("failed to parse decompressed SDP: %v", err)
			}
		} else {
			// Regular SDP without compression
			if err := json.Unmarshal(sdpStr, &offer); err != nil {
				return sessionID, nil, fmt.Errorf("failed to parse session description: %v", err)
			}
		}
	} else if sdpMap, ok := sessionDescription.(map[string]interface{}); ok {
		// Session description is already a map
		offer = sdpMap
	} else {
		return sessionID, nil, fmt.Errorf("invalid session description format")
	}

	// Format tracks for Cloudflare API
	cfTracks := make([]map[string]interface{}, 0)
	for _, track := range tracks {
		if trackMap, ok := track.(map[string]interface{}); ok {
			cfTracks = append(cfTracks, trackMap)
		}
	}

	// Publish tracks
	response, err := cs.PublishTracks(sessionID, offer, cfTracks)
	if err != nil {
		return sessionID, nil, fmt.Errorf("failed to publish tracks: %v", err)
	}

	return sessionID, response, nil
}
