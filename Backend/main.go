// filepath: d:\DAppMeetingNew\ProxyCalltoCL\main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

// Message defines the structure of WebSocket messages
type Message struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

// Response defines the structure of WebSocket responses
type Response struct {
	Action  string                 `json:"action"`
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (customize for production)
	},
}

var cfService *CloudflareService

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("Client connected: %s", conn.RemoteAddr())

	for {
		// Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		// Parse the incoming message
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			sendErrorResponse(conn, "", "Invalid message format")
			continue
		}

		// Process the message based on the action
		handleMessage(conn, msg)
	}

	log.Printf("Client disconnected: %s", conn.RemoteAddr())
}

func handleMessage(conn *websocket.Conn, msg Message) {
	var response Response

	switch strings.ToLower(msg.Action) {
	case "create_session":
		sessionID, err := cfService.CreateSession()
		if err != nil {
			sendErrorResponse(conn, msg.Action, err.Error())
			return
		}
		response = Response{
			Action:  msg.Action,
			Success: true,
			Data: map[string]interface{}{
				"sessionId": sessionID,
			},
		}

	case "publish_tracks":
		sessionID, ok := msg.Data["sessionId"].(string)
		if !ok {
			sendErrorResponse(conn, msg.Action, "sessionId is required")
			return
		}

		offer, ok := msg.Data["offer"].(map[string]interface{})
		if !ok {
			sendErrorResponse(conn, msg.Action, "offer is required")
			return
		}

		tracksData, ok := msg.Data["tracks"]
		if !ok {
			sendErrorResponse(conn, msg.Action, "tracks are required")
			return
		}

		// Convert tracks to the right format
		tracksBytes, _ := json.Marshal(tracksData)
		var tracks []map[string]interface{}
		json.Unmarshal(tracksBytes, &tracks)

		result, err := cfService.PublishTracks(sessionID, offer, tracks)
		if err != nil {
			sendErrorResponse(conn, msg.Action, err.Error())
			return
		}
		response = Response{
			Action:  msg.Action,
			Success: true,
			Data:    result,
		}

	case "pull_tracks":
		sessionID, ok := msg.Data["sessionId"].(string)
		if !ok {
			sendErrorResponse(conn, msg.Action, "sessionId is required")
			return
		}

		tracksData, ok := msg.Data["remoteTracks"]
		if !ok {
			sendErrorResponse(conn, msg.Action, "remoteTracks are required")
			return
		}

		// Convert tracks to the right format
		tracksBytes, _ := json.Marshal(tracksData)
		var remoteTracks []map[string]interface{}
		json.Unmarshal(tracksBytes, &remoteTracks)

		result, err := cfService.PullTracks(sessionID, remoteTracks)
		if err != nil {
			sendErrorResponse(conn, msg.Action, err.Error())
			return
		}
		response = Response{
			Action:  msg.Action,
			Success: true,
			Data:    result,
		}

	case "renegotiate":
		sessionID, ok := msg.Data["sessionId"].(string)
		if !ok {
			sendErrorResponse(conn, msg.Action, "sessionId is required")
			return
		}

		sessionDescription, ok := msg.Data["sessionDescription"].(map[string]interface{})
		if !ok {
			sendErrorResponse(conn, msg.Action, "sessionDescription is required")
			return
		}

		result, err := cfService.Renegotiate(sessionID, sessionDescription)
		if err != nil {
			sendErrorResponse(conn, msg.Action, err.Error())
			return
		}
		response = Response{
			Action:  msg.Action,
			Success: true,
			Data:    result,
		}

	case "close_tracks":
		sessionID, ok := msg.Data["sessionId"].(string)
		if !ok {
			sendErrorResponse(conn, msg.Action, "sessionId is required")
			return
		}

		tracksData, ok := msg.Data["tracks"]
		if !ok {
			sendErrorResponse(conn, msg.Action, "tracks are required")
			return
		}

		// Convert tracks to the right format
		tracksBytes, _ := json.Marshal(tracksData)
		var tracks []map[string]string
		json.Unmarshal(tracksBytes, &tracks)

		force := false
		if forceVal, ok := msg.Data["force"].(bool); ok {
			force = forceVal
		}

		sessionDescription := make(map[string]interface{})
		if sdVal, ok := msg.Data["sessionDescription"].(map[string]interface{}); ok {
			sessionDescription = sdVal
		}

		result, err := cfService.CloseTracks(sessionID, tracks, force, sessionDescription)
		if err != nil {
			sendErrorResponse(conn, msg.Action, err.Error())
			return
		}
		response = Response{
			Action:  msg.Action,
			Success: true,
			Data:    result,
		}

	case "get_session_state":
		sessionID, ok := msg.Data["sessionId"].(string)
		if !ok {
			sendErrorResponse(conn, msg.Action, "sessionId is required")
			return
		}

		result, err := cfService.GetSessionState(sessionID)
		if err != nil {
			sendErrorResponse(conn, msg.Action, err.Error())
			return
		}
		response = Response{
			Action:  msg.Action,
			Success: true,
			Data:    result,
		}

	default:
		sendErrorResponse(conn, msg.Action, "Unknown action")
		return
	}

	// Send the response
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

func sendErrorResponse(conn *websocket.Conn, action, errorMsg string) {
	response := Response{
		Action:  action,
		Success: false,
		Error:   errorMsg,
	}
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending error response: %v", err)
	}
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	// Initialize Cloudflare service
	baseURL := os.Getenv("CLOUDFLARE_BASE_URL")
	appID := os.Getenv("CLOUDFLARE_APP_ID")
	appSecret := os.Getenv("CLOUDFLARE_APP_SECRET")

	if baseURL == "" || appID == "" || appSecret == "" {
		log.Fatal("CLOUDFLARE_BASE_URL, CLOUDFLARE_APP_ID, and CLOUDFLARE_APP_SECRET environment variables are required")
	}

	cfService = NewCloudflareService(baseURL, appID, appSecret)

	// Setup WebSocket handler
	http.HandleFunc("/ws", handleWebSocket)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("WebSocket server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
