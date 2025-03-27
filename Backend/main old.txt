package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/klauspost/compress/zlib"
	// "golang.org/x/net/websocket"
)

// --- Constants and Configuration ---

var (
	port                   string
	cloudflareAppID        string
	cloudflareAppSecret    string
	secretKey              string
	cloudflareCallsBaseURL string
	cloudflareBasePath     string
	debug                  bool
	RPC_URL                string
	CONTRACT_ADDRESS       string
	ABI_JSON               string
	PRIVATE_KEY            string
)

func initConfig() {
	err := godotenv.Load() // Load .env file
	if err != nil {
		log.Println("Warning: Could not load .env file:", err) // Don't fatal, as we have defaults
	}

	// authRequired = getEnvBool("AUTH_REQUIRED", true)
	port = getEnv("PORT", "5000")
	cloudflareAppID = getEnv("CLOUDFLARE_APP_ID", "")
	cloudflareAppSecret = getEnv("CLOUDFLARE_APP_SECRET", "")
	secretKey = getEnv("JWT_SECRET", "thisisjustademokey")
	cloudflareCallsBaseURL = getEnv("CLOUDFLARE_APPS_URL", "https://rtc.live.cloudflare.com/v1/apps")
	cloudflareBasePath = fmt.Sprintf("%s/%s", cloudflareCallsBaseURL, cloudflareAppID)
	debug = getEnvBool("DEBUG", false)
	RPC_URL = "wss://bsc-testnet.publicnode.com"
	CONTRACT_ADDRESS = getEnv("CONTRACT_ADDRESS", "")
	ABI_JSON = getEnv("ABI_JSON", "")
	PRIVATE_KEY = getEnv("PRIVATE_KEY", "")

	if cloudflareAppID == "" || cloudflareAppSecret == "" {
		log.Fatal("CLOUDFLARE_APP_ID and CLOUDFLARE_APP_SECRET must be set")
	}

	// Thêm log để kiểm tra biến môi trường
	log.Println("Biến môi trường đã tải:")
	log.Println("CLOUDFLARE_APP_ID:", cloudflareAppID)
	log.Println("CLOUDFLARE_APP_SECRET:", cloudflareAppSecret)
	log.Println("CLOUDFLARE_APPS_URL:", cloudflareCallsBaseURL)
	log.Println("RPC_URL:", RPC_URL)
	log.Println("CONTRACT_ADDRESS:", CONTRACT_ADDRESS)
}

// Helper function to get environment variables with default values
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func init() {
	// Initialize maps
	rooms.m = make(map[string]*Room)
	users.m = make(map[string]*User)
	wsConnections.m = make(map[string]map[string]*websocket.Conn)
}

// Helper function to get boolean environment variables
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue // Return default if parsing fails
	}
	return b
}

// --- Data Structures ---

type Room struct {
	RoomId       string                 `json:"roomId"` // Thêm trường này
	Name         string                 `json:"name"`
	Metadata     map[string]interface{} `json:"metadata"`
	Participants []*Participant         `json:"participants"`
	CreatedAt    *big.Int               `json:"createdAt"`
	sync.RWMutex                        // Protects concurrent access to the room
}

type Participant struct {
	UserID          string   `json:"userId"`
	SessionID       string   `json:"sessionId"`
	CreatedAt       *big.Int `json:"createdAt"`
	PublishedTracks []string `json:"publishedTracks"`
}

type User struct {
	UserID      string `json:"userId"`
	Username    string `json:"username"`
	IsModerator bool   `json:"isModerator"`
	Role        string `json:"role"`
}

type SessionResponse struct {
	SessionId     string        `json:"sessionId"`
	OtherSessions []SessionInfo `json:"otherSessions"`
}

type SessionInfo struct {
	UserId          string   `json:"userId"`
	SessionId       string   `json:"sessionId"`
	PublishedTracks []string `json:"publishedTracks"`
}

// Define TrackPullAnswer event struct
type TrackPullAnswerEvent struct {
	RoomId       string         `abi:"roomId"`
	Requester    common.Address `abi:"requester"`
	SessionId    string         `abi:"sessionId"`
	ResponseData []byte         `abi:"responseData"`
}

type LeftRoomEvent struct {
	RoomId    string         `abi:"roomId"`
	User      common.Address `abi:"user"`
	SessionId string         `abi:"sessionId"`
}

// Use a concurrent-safe map for rooms.
var rooms = struct {
	sync.RWMutex
	m map[string]*Room
}{m: make(map[string]*Room)}

var users = struct {
	sync.RWMutex
	m map[string]*User
}{m: make(map[string]*User)}

var wsConnections = struct {
	sync.RWMutex
	m map[string]map[string]*websocket.Conn
}{m: make(map[string]map[string]*websocket.Conn)}

// --- Cloudflare API Interaction Functions ---

func createCloudflareSession() (string, error) {
	url := fmt.Sprintf("%s/sessions/new", cloudflareBasePath)
	log.Printf("[Cloudflare API] Creating new session: %s", url)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("[Cloudflare API Error] Failed to create request: %v", err)
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+cloudflareAppSecret)
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

func renegotiateWithCloudflare(sessionId string, body map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/sessions/%s/renegotiate", cloudflareBasePath, sessionId)

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cloudflareAppSecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

func publishTracks(c *gin.Context, roomId, sessionId string, address common.Address, offer map[string]interface{}, tracks []struct {
	TrackName string `json:"trackName"`
	Mid       string `json:"mid"`
	Location  string `json:"location"`
}) {
	// Check if the room exists
	rooms.RLock()
	room, ok := rooms.m[roomId]
	rooms.RUnlock()
	if !ok {
		log.Printf("[Error] Room %s not found", roomId)
		return
	}

	var participant *Participant
	room.RLock()
	for _, p := range room.Participants {
		if strings.EqualFold(p.UserID, address.Hex()) {
			participant = p
			break
		}
	}

	room.RUnlock()
	if participant == nil {
		log.Printf("[Error] Participant with session ID %s not found in room %s", sessionId, roomId)
		// them participant vao room
		participant = &Participant{
			UserID:          address.Hex(),
			SessionID:       sessionId,
			CreatedAt:       big.NewInt(time.Now().Unix()),
			PublishedTracks: []string{},
		}
		room.Lock()
		room.Participants = append(room.Participants, participant)
		room.Unlock()
	}

	room.Lock()
	// add session id to participant
	participant.SessionID = sessionId
	room.Unlock()

	// Prepare track data for Cloudflare API
	trackData := make([]map[string]interface{}, len(tracks))
	trackNames := make([]string, len(tracks))
	for i, t := range tracks {
		trackData[i] = map[string]interface{}{
			"trackName": t.TrackName,
			"mid":       t.Mid,
			"location":  t.Location,
		}
		trackNames[i] = t.TrackName
	}

	// Call Cloudflare API
	requestBody := map[string]interface{}{
		"sessionDescription": offer,
		"tracks":             trackData,
	}
	requestBodyString, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("[Error] Failed to marshal request body: %v", err)
		return
	}
	requestBodyStr := string(requestBodyString)

	// Get the Cloudflare API URL
	url := fmt.Sprintf("%s/sessions/%s/tracks/new", cloudflareBasePath, sessionId)
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("[Error] Failed to marshal request: %v", err)
		return
	}
	// Log the request we're about to send to Cloudflare
	log.Printf("Sending track publish request to Cloudflare: %s", string(jsonData))
	log.Printf("Cloudflare API URL: %s", url)

	// Create and configure the HTTP request
	cfReq, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBodyStr)))
	if err != nil {
		log.Printf("[Error] Failed to create request: %v", err)
		return
	}

	cfReq.Header.Set("Authorization", "Bearer "+cloudflareAppSecret)
	cfReq.Header.Set("Content-Type", "application/json")
	cfReq.Header.Set("Accept", "*/*")
	cfReq.Header.Set("Accept-Encoding", "identity")

	// log.Printf("Cloudflare API request headers: %v", cfReq.Header)
	// log.Printf("Cloudflare API request body: %s", string(jsonData))

	// Execute the HTTP request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(cfReq)
	if err != nil {
		log.Printf("[Error] Failed to send request to Cloudflare: %v", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Error] Failed to read response body: %v", err)
		return
	}

	// // Log the complete response for debugging
	// log.Printf("Cloudflare API response (status %d): %s", resp.StatusCode, string(body))

	// Check if we got an error response
	if resp.StatusCode >= 400 {
		log.Printf("[Error] Cloudflare API returned error status %d: %s", resp.StatusCode, string(body))
		return
	}

	// Parse the JSON response
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("[Error] Failed to parse response JSON: %v", err)
		return
	}

	// Update participant's published tracks
	room.Lock()
	for _, t := range tracks {
		if !contains(participant.PublishedTracks, t.TrackName) {
			participant.PublishedTracks = append(participant.PublishedTracks, t.TrackName)
			participant.SessionID = sessionId
		}
	}
	room.Unlock()

	log.Printf("Track published event broadcasted for session %s with tracks: %v", sessionId, trackNames)

	// Only return JSON response if this was called from an API endpoint
	if c != nil {
		c.JSON(http.StatusOK, data)
	} else {
		log.Printf("Published tracks successfully !")
		// send answer to smart constract
		answerData, err := json.Marshal(data)
		if err != nil {
			log.Printf("[Error] Failed to marshal Cloudflare answer: %v", err)
			return
		}

		// Nén dữ liệu answer
		var compressed bytes.Buffer
		w := zlib.NewWriter(&compressed)
		w.Write(answerData)
		w.Close()

		// Gửi answer về qua smart contract
		log.Printf("[Smart Contract] Sending TrackPublishedAnswer to contract for session %s", sessionId)
		publishTrackAnswerOnContract(
			roomId,
			address,
			compressed.Bytes(),
		)
	}
}

func contains(s1 []string, s2 string) bool {
	for _, v := range s1 {
		if v == s2 {
			return true
		}
	}
	return false
}

type SessionRenegotiatedEvent struct {
	SessionId string `abi:"sessionId"`
	SDP       []byte `abi:"sdp"`     // Changed to []byte for compressed data
	SDPType   string `abi:"sdpType"` // Changed field name to match contract
}

func renegotiateSessionSM(event *SessionRenegotiatedEvent) (map[string]interface{}, error) {
	log.Printf("[Smart Contract] Processing session renegotiation for session %s", event.SessionId)

	// Decompress the SDP data if it's compressed
	sdpStr, err := decompressTrackData(event.SDP)
	if err != nil {
		log.Printf("[Error] Failed to decompress SDP data: %v", err)
		return nil, err
	}

	body := map[string]interface{}{
		"sessionDescription": map[string]string{
			"sdp":  sdpStr,
			"type": event.SDPType,
		},
	}

	// Call Cloudflare API for renegotiation
	data, err := renegotiateWithCloudflare(event.SessionId, body)
	if err != nil {
		if data != nil && data["errorCode"] != nil {
			log.Printf("[Error] Cloudflare renegotiation failed: %v", data["errorDescription"])
			return data, err
		}
		log.Printf("[Error] Renegotiation failed: %v", err)
		return nil, err
	}

	// Prepare answer data for smart contract
	answerData, err := json.Marshal(data)
	if err != nil {
		log.Printf("[Error] Failed to marshal renegotiation answer: %v", err)
		return data, nil
	}

	// Nén dữ liệu answer
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	w.Write(answerData)
	w.Close()

	// Gửi answer về qua smart contract
	log.Printf("[Smart Contract] Sending renegotiation answer to contract for session %s", event.SessionId)
	publishRenegotiateAnswerOnContract(
		event.SessionId,
		compressed.Bytes(),
	)

	log.Printf("[Success] Session renegotiation completed for session %s", event.SessionId)
	return data, nil
}

func leaveRoomSM(event *LeftRoomEvent) {
	roomId := event.RoomId
	userId := event.User.Hex()

	rooms.RLock()
	room, ok := rooms.m[roomId]
	rooms.RUnlock()
	if !ok {
		log.Printf("[Error] Room %s not found", roomId)
		return
	}

	participantIndex := -1

	room.Lock()
	for i, p := range room.Participants {
		if p.UserID == userId {
			participantIndex = i
			break
		}
	}

	if participantIndex != -1 {
		room.Participants = append(room.Participants[:participantIndex], room.Participants[participantIndex+1:]...)

		// If room is empty, delete it
		if len(room.Participants) == 0 {
			rooms.Lock()
			delete(rooms.m, roomId)
			rooms.Unlock()
		}
	}
	room.Unlock()
}

func processJoinRoomRequest(roomId string, userAddress common.Address) {
	log.Printf("[Join Room] Processing join request for room %s from user %s", roomId, userAddress.Hex())

	// Check if room exists in our local cache
	rooms.RLock()
	room, exists := rooms.m[roomId]
	rooms.RUnlock()

	if !exists {
		log.Printf("[Join Room] Room %s not found, creating new room", roomId)
		room = &Room{
			RoomId:       roomId,
			Name:         "New Room",
			Metadata:     make(map[string]interface{}),
			Participants: make([]*Participant, 0),
			CreatedAt:    big.NewInt(time.Now().Unix()),
		}

		rooms.Lock()
		rooms.m[roomId] = room
		rooms.Unlock()
	}

	// // Create a Cloudflare session
	// sessionID, err := createCloudflareSession()
	// if err != nil {
	// 	log.Printf("[Join Room] Failed to create Cloudflare session: %v", err)
	// 	return
	// }

	// log.Printf("[Join Room] Created Cloudflare session: %s", sessionID)

	// Add participant to the room locally
	participant := &Participant{
		UserID:          userAddress.Hex(),
		SessionID:       "",
		CreatedAt:       big.NewInt(time.Now().Unix()),
		PublishedTracks: make([]string, 0),
	}

	room.Lock()
	room.Participants = append(room.Participants, participant)
	room.Unlock()

	log.Printf("[Join Room] Added participant to room %s: %s", roomId, userAddress.Hex())

}

// Function to update session ID on the smart contract
func updateSessionIDOnContract(roomId string, userAddress common.Address, sessionID string) {
	log.Printf("[Smart Contract] Updating session ID on contract for room %s, user %s, sessionID %s",
		roomId, userAddress.Hex(), sessionID)

	// Connect to the blockchain if not already connected
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		log.Printf("[Smart Contract] Failed to connect to blockchain: %v", err)
		return
	}

	// Convert contract address
	contractAddress := common.HexToAddress(CONTRACT_ADDRESS)

	// Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(ABI_JSON))
	if err != nil {
		log.Printf("[Smart Contract] Failed to parse ABI: %v", err)
		return
	}

	// Load private key
	privateKey, err := crypto.HexToECDSA(PRIVATE_KEY)
	if err != nil {
		log.Printf("[Smart Contract] Failed to parse private key: %v", err)
		return
	}

	// Get public key and address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Printf("[Smart Contract] Error converting public key")
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Get nonce for transaction
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Printf("[Smart Contract] Failed to get nonce: %v", err)
		return
	}

	// Get gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Printf("[Smart Contract] Failed to get gas price: %v", err)
		return
	}

	// Pack the data for the function call
	data, err := parsedABI.Pack("setSessionIDByBackend", roomId, userAddress, sessionID)
	if err != nil {
		log.Printf("[Smart Contract] Failed to pack transaction data: %v", err)
		return
	}

	// Set gas limit
	gasLimit := uint64(300000)

	// Create the transaction
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, data)

	// BSC Testnet chain ID is 97
	chainID := big.NewInt(97)

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		log.Printf("[Smart Contract] Failed to sign transaction: %v", err)
		return
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Printf("[Smart Contract] Failed to send transaction: %v", err)
		return
	}

	log.Printf("[Smart Contract] Session ID update transaction sent. Hash: %s", signedTx.Hash().Hex())
}

// Define TrackPublishedCompressed event struct
type TrackPublishedCompressedEvent struct {
	RoomId      string         `abi:"roomId"`
	Participant common.Address `abi:"participant"`
	SessionId   string         `abi:"sessionId"`
	TrackData   []byte         `abi:"trackData"`
}

// Function to handle TrackPublishedCompressed event

func handleTrackPublishedCompressed(event *TrackPublishedCompressedEvent) {
	// Decompress the data
	data, err := decompressTrackData(event.TrackData)
	if err != nil {
		log.Printf("[Error] Failed to decompress track data: %v", err)
		return
	}

	// Log raw JSON after decompression for debugging
	log.Printf("[Debug] Decompressed track data: %s", data)

	// Parse the decompressed JSON data
	var trackPublishData struct {
		Offer  map[string]interface{} `json:"offer"`
		Tracks []struct {
			TrackName string `json:"trackName"`
			Mid       string `json:"mid"`
			Location  string `json:"location"`
		} `json:"tracks"`
	}

	// Chuyển string thành []byte trước khi unmarshal
	if err := json.Unmarshal([]byte(data), &trackPublishData); err != nil {
		log.Printf("[Error] Failed to parse track data: %v", err)
		log.Printf("[Error] Failed JSON: %s", data)
		return
	}
	sessionID := ""
	// Nếu event.SessionId rỗng thì gọi hàm createCloudflareSession() còn không thì sự dụng giá trị sẵn có của nó
	if event.SessionId == "" {
		sessionID, err = createCloudflareSession()
	} else {
		sessionID = event.SessionId
	}

	log.Printf("[Debug] Publishing tracks for session %s via smart contract: %+v",
		sessionID, trackPublishData)

	// Call publishTracks with nil context (it's not from an HTTP request)
	publishTracks(nil, event.RoomId, sessionID, event.Participant, trackPublishData.Offer, trackPublishData.Tracks)
	updateSessionIDOnContract(event.RoomId, event.Participant, sessionID)
}

// Helper function to decompress track data
func decompressTrackData(compressedData []byte) (string, error) {
	// Create a new zlib reader
	r, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return "", fmt.Errorf("failed to create zlib reader: %v", err)
	}
	defer r.Close()

	// Read all decompressed data
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return "", fmt.Errorf("failed to decompress data: %v", err)
	}

	return buf.String(), nil
}

// Function to send track published answer back to smart contract
func publishTrackAnswerOnContract(roomId string, participant common.Address, compressedData []byte) {
	log.Printf("[Smart Contract] Publishing track answer for room %s, participant %s",
		roomId, participant.Hex())

	// Connect to the blockchain if not already connected
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		log.Printf("[Smart Contract] Failed to connect to blockchain: %v", err)
		return
	}

	// Convert contract address
	contractAddress := common.HexToAddress(CONTRACT_ADDRESS)

	// Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(ABI_JSON))
	if err != nil {
		log.Printf("[Smart Contract] Failed to parse ABI: %v", err)
		return
	}

	// Load private key
	privateKey, err := crypto.HexToECDSA(PRIVATE_KEY)
	if err != nil {
		log.Printf("[Smart Contract] Failed to parse private key: %v", err)
		return
	}

	// Get public key and address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Printf("[Smart Contract] Error converting public key")
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Get nonce for transaction
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Printf("[Smart Contract] Failed to get nonce: %v", err)
		return
	}

	// Get gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Printf("[Smart Contract] Failed to get gas price: %v", err)
		return
	}

	// Pack the data for the function call
	data, err := parsedABI.Pack("publishTrackAnswer", roomId, participant, compressedData)

	if err != nil {
		log.Printf("[Smart Contract] Failed to pack transaction data: %v", err)
		return
	}

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	})
	if err != nil {
		log.Printf("[Smart Contract] Failed to estimate gas: %v", err)
		return
	}

	// Create the transaction
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, data)

	// BSC Testnet chain ID is 97
	chainID := big.NewInt(97)

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		log.Printf("[Smart Contract] Failed to sign transaction: %v", err)
		return
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Printf("[Smart Contract] Failed to send transaction: %v", err)
		return
	}

	log.Printf("[Smart Contract] Track answer publication transaction sent. Hash: %s", signedTx.Hash().Hex())
}

// Define TrackPullRequestCompressed event struct
type TrackPullRequestCompressedEvent struct {
	RoomId          string         `abi:"roomId"`
	Requester       common.Address `abi:"requester"`
	SessionId       string         `abi:"sessionId"`
	RemoteSessionId string         `abi:"remoteSessionId"`
	TrackName       string         `abi:"trackName"`
	RequestData     []byte         `abi:"requestData"`
}

// Function to handle TrackPullRequestCompressed event
func handleTrackPullRequestCompressed(event *TrackPullRequestCompressedEvent) {
	log.Printf("[Smart Contract Event] Received TrackPullRequestCompressed event for room %s from %s",
		event.RoomId, event.Requester.Hex())

	log.Printf("[Pull Request] Session %s requesting to pull track %s from session %s",
		event.SessionId, event.TrackName, event.RemoteSessionId)

	// Decompress request data if provided (may not be needed for simple pull requests)
	var requestData map[string]interface{}
	if len(event.RequestData) > 0 {
		decompressedData, err := decompressTrackData(event.RequestData)
		if err != nil {
			log.Printf("[Error] Failed to decompress pull request data: %v", err)
			// Continue anyway as the essential data is in the event params
		} else {
			if err := json.Unmarshal([]byte(decompressedData), &requestData); err != nil {
				log.Printf("[Warning] Failed to parse decompressed request data: %v", err)
				// Continue anyway as this is optional
			}
		}
	}

	// Prepare the pull request for Cloudflare with the correct endpoint URL
	pullURL := fmt.Sprintf("%s/sessions/%s/tracks/new", cloudflareBasePath, event.SessionId)
	log.Printf("[Debug] Pull URL: %s", pullURL)

	// Prepare the request body with the correct structure for Cloudflare API
	reqBody := map[string]interface{}{
		"tracks": []map[string]interface{}{
			{
				"location":  "remote",
				"sessionId": event.RemoteSessionId,
				"trackName": event.TrackName,
			},
		},
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("[Error] Failed to marshal pull request body: %v", err)
		return
	}
	log.Printf("[Debug] Pull request body: %s", string(reqBodyBytes))

	// Make HTTP request to Cloudflare API
	req, err := http.NewRequest("POST", pullURL, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		log.Printf("[Error] Failed to create pull request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cloudflareAppSecret)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Error] Failed to execute pull request: %v", err)

		// Prepare error response
		errorResponse := map[string]interface{}{
			"errorCode":        500,
			"errorDescription": fmt.Sprintf("Failed to pull track: %v", err),
		}

		responseBytes, _ := json.Marshal(errorResponse)
		// Compress the error response
		var compressed bytes.Buffer
		w := zlib.NewWriter(&compressed)
		w.Write(responseBytes)
		w.Close()

		// Send the error response back to the smart contract
		publishTrackPullAnswerOnContract(
			event.RoomId,
			event.Requester,
			event.SessionId,
			compressed.Bytes(),
		)

		return
	}
	defer resp.Body.Close()

	// Read and process the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Error] Failed to read pull response body: %v", err)
		return
	}
	log.Printf("[Debug] Cloudflare response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		log.Printf("[Error] Pull request failed with status %d: %s", resp.StatusCode, string(body))

		// Prepare error response
		errorResponse := map[string]interface{}{
			"errorCode":        resp.StatusCode,
			"errorDescription": fmt.Sprintf("Cloudflare returned error: %s", string(body)),
		}

		responseBytes, _ := json.Marshal(errorResponse)
		// Compress the error response
		var compressed bytes.Buffer
		w := zlib.NewWriter(&compressed)
		w.Write(responseBytes)
		w.Close()

		// Send the error response back to the smart contract
		publishTrackPullAnswerOnContract(
			event.RoomId,
			event.Requester,
			event.SessionId,
			compressed.Bytes(),
		)

		return
	}

	log.Printf("[Pull Success] Successfully pulled track %s from session %s for session %s",
		event.TrackName, event.RemoteSessionId, event.SessionId)

	// Parse the response
	var pullResponse map[string]interface{}
	if err := json.Unmarshal(body, &pullResponse); err != nil {
		log.Printf("[Error] Failed to parse pull response: %v", err)
		return
	}

	// Check if immediate renegotiation is required
	if requiresRenegotiation, ok := pullResponse["requiresImmediateRenegotiation"].(bool); ok && requiresRenegotiation {
		log.Printf("[Info] Track pull requires immediate renegotiation")
		// You could add this information to the response to the smart contract if needed
	}

	// Compress the response
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	if _, err := w.Write(body); err != nil {
		log.Printf("[Error] Failed to compress pull response: %v", err)
		return
	}
	if err := w.Close(); err != nil {
		log.Printf("[Error] Failed to close zlib writer: %v", err)
		return
	}

	// Send the response back to the smart contract
	publishTrackPullAnswerOnContract(
		event.RoomId,
		event.Requester,
		event.SessionId,
		compressed.Bytes(),
	)
}

// Function to send track pull answer back to smart contract
func publishTrackPullAnswerOnContract(roomId string, requester common.Address, sessionId string, compressedData []byte) {
	log.Printf("[Smart Contract] Publishing track pull answer for room %s, requester %s, sessionId %s",
		roomId, requester.Hex(), sessionId)

	// Connect to the blockchain if not already connected
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		log.Printf("[Error] Failed to connect to Ethereum client: %v", err)
		return
	}

	// Convert contract address
	contractAddress := common.HexToAddress(CONTRACT_ADDRESS)

	// Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(ABI_JSON))
	if err != nil {
		log.Printf("[Error] Failed to parse ABI: %v", err)
		return
	}

	// Load private key
	privateKey, err := crypto.HexToECDSA(PRIVATE_KEY)
	if err != nil {
		log.Printf("[Error] Failed to parse private key: %v", err)
		return
	}

	// Get public key and address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Printf("[Error] Failed to get ECDSA public key")
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Get nonce for transaction
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Printf("[Error] Failed to get nonce: %v", err)
		return
	}

	// Get gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Printf("[Error] Failed to get gas price: %v", err)
		return
	}

	// Pack the data for the function call
	data, err := parsedABI.Pack("publishTrackPullAnswer", roomId, requester, sessionId, compressedData)
	if err != nil {
		log.Printf("[Error] Failed to pack data: %v", err)
		return
	}

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &contractAddress,
		Data: data,
	})
	if err != nil {
		log.Printf("[Smart Contract] Failed to estimate gas: %v", err)
		return
	}

	// Create the transaction
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, data)

	// BSC Testnet chain ID is 97
	chainID := big.NewInt(97)

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		log.Printf("[Error] Failed to sign transaction: %v", err)
		return
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Printf("[Error] Failed to send transaction: %v", err)
		return
	}

	log.Printf("[Smart Contract] Track pull answer transaction sent. Hash: %s", signedTx.Hash().Hex())
}

// Function to send renegotiation answer back to smart contract
func publishRenegotiateAnswerOnContract(sessionId string, sdp []byte) {
	log.Printf("[Smart Contract] Publishing renegotiation answer for session %s", sessionId)

	// Connect to the blockchain
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		log.Printf("[Smart Contract] Failed to connect to blockchain: %v", err)
		return
	}

	// Convert contract address
	contractAddress := common.HexToAddress(CONTRACT_ADDRESS)

	// Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(ABI_JSON))
	if err != nil {
		log.Printf("[Smart Contract] Failed to parse ABI: %v", err)
		return
	}

	// Load private key
	privateKey, err := crypto.HexToECDSA(PRIVATE_KEY)
	if err != nil {
		log.Printf("[Smart Contract] Failed to parse private key: %v", err)
		return
	}

	// Get public key and address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Printf("[Smart Contract] Error converting public key")
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Get nonce for transaction
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Printf("[Smart Contract] Failed to get nonce: %v", err)
		return
	}

	// Get gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Printf("[Smart Contract] Failed to get gas price: %v", err)
		return
	}

	// Pack the data for the function call
	data, err := parsedABI.Pack("AnswerRenegotiate", sessionId, sdp)
	if err != nil {
		log.Printf("[Smart Contract] Failed to pack transaction data: %v", err)
		return
	}

	// Estimate gas
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &contractAddress,
		Data: data,
	})
	if err != nil {
		log.Printf("[Smart Contract] Failed to estimate gas: %v", err)
		return
	}

	// Create the transaction
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, data)

	// BSC Testnet chain ID is 97
	chainID := big.NewInt(97)

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		log.Printf("[Smart Contract] Failed to sign transaction: %v", err)
		return
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Printf("[Smart Contract] Failed to send transaction: %v", err)
		return
	}

	log.Printf("[Smart Contract] Renegotiation answer transaction sent. Hash: %s", signedTx.Hash().Hex())
}

func main() {
	initConfig() // Initialize configuration

	r := gin.Default()

	// Tạo Phòng và lắng nghe smrtract
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		log.Fatalf("Không thể kết nối đến BSC Testnet: %v", err)
	}
	fmt.Println("Đã kết nối đến BSC Testnet!")

	contractAddress := common.HexToAddress(CONTRACT_ADDRESS)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	contractAbi, err := abi.JSON(strings.NewReader(ABI_JSON))
	if err != nil {
		log.Fatalf("Không thể phân tích ABI: %v", err)
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	fmt.Println("Đã kết nối đến BSC Testnet!", sub)
	if err != nil {
		log.Fatalf("Không thể đăng ký lắng nghe sự kiện: %v", err)
	}
	fmt.Println("Đang lắng nghe các sự kiện ...")

	type RoomCreatedEvent struct {
		RoomId    string         `abi:"roomId"`
		Name      string         `abi:"name"`
		Metadata  string         `abi:"metadata"`
		Owner     common.Address `abi:"owner"`
		Timestamp *big.Int       `abi:"timestamp"`
		IsActive  bool           `abi:"isActive"`
	}
	// Define JoinRoomRequest event struct
	type JoinRoomRequestEvent struct {
		RoomId    string         `abi:"roomId"`
		User      common.Address `abi:"user"`
		Timestamp *big.Int       `abi:"timestamp"`
	}

	go func() {
		for eventLog := range logs {
			// Check event signature to determine which event it is
			if len(eventLog.Topics) > 0 {
				eventSignature := eventLog.Topics[0].Hex()

				switch eventSignature {
				case contractAbi.Events["RoomCreated"].ID.Hex():
					event := new(RoomCreatedEvent)
					if err := contractAbi.UnpackIntoInterface(event, "RoomCreated", eventLog.Data); err != nil {
						log.Printf("Lỗi giải mã sự kiện RoomCreated: %v", err)
						continue
					}

					// Process RoomCreated event
					if eventLog.TxHash != (common.Hash{}) {
						var metadata map[string]interface{}
						if err := json.Unmarshal([]byte(event.Metadata), &metadata); err != nil {
							log.Printf("Lỗi parse metadata: %v", err)
							metadata = make(map[string]interface{})
						}

						createTime := time.Unix(event.Timestamp.Int64(), 0)
						fmt.Printf("Sự kiện RoomCreated:\nRoom ID: %s\nTên: %s\nMetadata: %s\nOwner: %s\nThời gian: %v\nActive: %v\n",
							event.RoomId, event.Name, event.Metadata, event.Owner.Hex(), createTime, event.IsActive)

						room := &Room{
							RoomId:    event.RoomId,
							Name:      event.Name,
							Metadata:  metadata,
							CreatedAt: event.Timestamp,
						}

						rooms.Lock()
						rooms.m[event.RoomId] = room
						rooms.Unlock()
					}

				case contractAbi.Events["JoinRoomRequest"].ID.Hex():
					event := new(JoinRoomRequestEvent)
					if err := contractAbi.UnpackIntoInterface(event, "JoinRoomRequest", eventLog.Data); err != nil {
						log.Printf("Lỗi giải mã sự kiện JoinRoomRequest: %v", err)
						continue
					}

					// Process JoinRoomRequest event with our new function
					if eventLog.TxHash != (common.Hash{}) {
						fmt.Printf("Sự kiện JoinRoomRequest:\nRoom ID: %s\nUser: %s\nThời gian: %v\n",
							event.RoomId, event.User.Hex(), time.Unix(event.Timestamp.Int64(), 0))

						// Process join room request
						processJoinRoomRequest(event.RoomId, event.User)
					}

				case contractAbi.Events["TrackPublishedCompressed"].ID.Hex():
					event := new(TrackPublishedCompressedEvent)
					if err := contractAbi.UnpackIntoInterface(event, "TrackPublishedCompressed", eventLog.Data); err != nil {
						log.Printf("Lỗi giải mã sự kiện TrackPublishedCompressed: %v", err)
						continue
					}

					// Process TrackPublishedCompressed event with our new function
					if eventLog.TxHash != (common.Hash{}) {
						fmt.Printf("Sự kiện TrackPublishedCompressed:\nRoom ID: %s\nParticipant: %s\nTrackData size: %d bytes\n",
							event.RoomId, event.Participant.Hex(), len(event.TrackData))

						// Process track published request
						handleTrackPublishedCompressed(event)
					}

				case contractAbi.Events["TrackPullRequestCompressed"].ID.Hex():
					event := new(TrackPullRequestCompressedEvent)
					if err := contractAbi.UnpackIntoInterface(event, "TrackPullRequestCompressed", eventLog.Data); err != nil {
						log.Printf("Lỗi giải mã sự kiện TrackPullRequestCompressed: %v", err)
						continue
					}

					// Process TrackPullRequestCompressed event with our new function
					if eventLog.TxHash != (common.Hash{}) {
						fmt.Printf("Sự kiện TrackPullRequestCompressed:\nRoom ID: %s\nRequester: %s\nSession ID: %s\nRemote Session ID: %s\nTrack Name: %s\nRequest Data size: %d bytes\n",
							event.RoomId, event.Requester.Hex(), event.SessionId, event.RemoteSessionId, event.TrackName, len(event.RequestData))

						// Process track pull request
						handleTrackPullRequestCompressed(event)
					}

				case contractAbi.Events["SessionRenegotiated"].ID.Hex():
					event := new(SessionRenegotiatedEvent)
					if err := contractAbi.UnpackIntoInterface(event, "SessionRenegotiated", eventLog.Data); err != nil {
						log.Printf("Lỗi giải mã sự kiện SessionRenegotiated: %v", err)
						continue
					}

					// Process SessionRenegotiated event
					if eventLog.TxHash != (common.Hash{}) {
						fmt.Printf("Sự kiện SessionRenegotiated:\nSession ID: %s\nSDP Type: %s\n",
							event.SessionId, event.SDPType)
						renegotiateSessionSM(event)
					}

				case contractAbi.Events["LeftRoom"].ID.Hex():
					event := new(LeftRoomEvent)
					if err := contractAbi.UnpackIntoInterface(event, "LeftRoom", eventLog.Data); err != nil {
						log.Printf("Lỗi giải mã sự kiện LeftRoom: %v", err)
						continue
					}

					if eventLog.TxHash != (common.Hash{}) {
						fmt.Printf("Sự kiện LeftRoom:\nRoom ID: %s\nUser: %s\n",
							event.RoomId, event.User.Hex())
						// Process left room event
						leaveRoomSM(event)
					}
				}

			}
		}
	}()

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
