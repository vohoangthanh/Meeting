package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	contract "dappmeetingnew/constract"
	"dappmeetingnew/handle"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// Configuration constants
var (
	ethereumNodeURL     string
	contractAddress     string
	cloudflareBaseURL   string
	cloudflareAppID     string
	cloudflareAppSecret string
)

func init() {
	// Load environment variables
	paths := []string{
		".env",       // Current directory
		"../.env",    // Parent directory
		"../../.env", // Two levels up
		"D:/DAppMeetingNew/backend/.env",
	}

	loaded := false
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			loaded = true
			log.Printf("Loaded environment from %s\n", path)
			break
		}
	}

	if !loaded {
		log.Println("Warning: No .env file found. Using default or environment values.")
	}

	// Load configurations with fallback to defaults
	ethereumNodeURL = getEnv("ETHEREUM_NODE_URL", "wss://bsc-testnet-rpc.publicnode.com")
	contractAddress = getEnv("CONTRACT_ADDRESS", "0xEf497BdCD80Aaa420271F5D2eAd9C1E70c2930E0")
	cloudflareBaseURL = getEnv("CLOUDFLARE_BASE_URL", "https://rtc.live.cloudflare.com/v1/apps")
	cloudflareAppID = getEnv("CLOUDFLARE_APP_ID", "977c39eac9a8fc03a471d7da6a7d66e0")
	cloudflareAppSecret = getEnv("CLOUDFLARE_APP_SECRET", "f069115dbeb5847040e7dbb7f9c79772fa923534fe1f74e3a62d54561aa12118")

	log.Printf("Ethereum Node URL: %s\n", ethereumNodeURL)
	log.Printf("Contract Address: %s\n", contractAddress)
	log.Printf("Cloudflare Base URL: %s\n", cloudflareBaseURL)

	if cloudflareAppID == "" || cloudflareAppSecret == "" {
		log.Println("Warning: Cloudflare credentials not set. Please check your environment variables.")
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func main() {
	// Connect to Ethereum node
	client, err := ethclient.Dial(ethereumNodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network: %v", err)
	}
	defer client.Close()
	fmt.Println("Connected to Ethereum node:", ethereumNodeURL)

	// Create a new instance of the contract binding
	address := common.HexToAddress(contractAddress)
	contractInstance, err := contract.NewContract(address, client)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}
	fmt.Println("Contract instance created at address:", contractAddress)

	// Initialize Cloudflare service
	cloudflareService := handle.NewCloudflareService(cloudflareBaseURL, cloudflareAppID, cloudflareAppSecret)
	fmt.Println("Cloudflare service initialized")

	// Initialize SMCallManager for transaction handling
	smCallManager, err := handle.NewSMCallManager(client, address)
	if err != nil {
		log.Fatalf("Failed to initialize SM Call Manager: %v", err)
	}
	defer smCallManager.Close()
	fmt.Println("SM Call Manager initialized")

	// Initialize EventHandler
	eventHandler := handle.NewEventHandler(contractInstance, cloudflareService, smCallManager)
	fmt.Println("Event Handler initialized")

	// Create a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up channel for handling OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Create channels for contract events
	participantJoinedCh := make(chan *contract.ContractParticipantJoined)
	participantLeftCh := make(chan *contract.ContractParticipantLeft)
	trackAddedCh := make(chan *contract.ContractTrackAdded)
	eventToBackendCh := make(chan *contract.ContractEventForwardedToBackend)
	eventToFrontendCh := make(chan *contract.ContractEventForwardedToFrontend)

	// Create watch options
	watchOpts := &bind.WatchOpts{
		Context: ctx,
	}

	// Start watching for ParticipantJoined events
	participantJoinedSub, err := contractInstance.WatchParticipantJoined(watchOpts, participantJoinedCh)
	if err != nil {
		log.Fatalf("Failed to watch ParticipantJoined events: %v", err)
	}
	defer participantJoinedSub.Unsubscribe()

	// Start watching for ParticipantLeft events
	participantLeftSub, err := contractInstance.WatchParticipantLeft(watchOpts, participantLeftCh)
	if err != nil {
		log.Fatalf("Failed to watch ParticipantLeft events: %v", err)
	}
	defer participantLeftSub.Unsubscribe()

	// Start watching for TrackAdded events
	trackAddedSub, err := contractInstance.WatchTrackAdded(watchOpts, trackAddedCh)
	if err != nil {
		log.Fatalf("Failed to watch TrackAdded events: %v", err)
	}
	defer trackAddedSub.Unsubscribe()

	// Start watching for EventForwardedToBackend events
	eventToBackendSub, err := contractInstance.WatchEventForwardedToBackend(watchOpts, eventToBackendCh)
	if err != nil {
		log.Fatalf("Failed to watch EventForwardedToBackend events: %v", err)
	}
	defer eventToBackendSub.Unsubscribe()

	// Start watching for EventForwardedToFrontend events
	eventToFrontendSub, err := contractInstance.WatchEventForwardedToFrontend(watchOpts, eventToFrontendCh)
	if err != nil {
		log.Fatalf("Failed to watch EventForwardedToFrontend events: %v", err)
	}
	defer eventToFrontendSub.Unsubscribe()

	// Combine subscription error channels
	errorCh := make(chan error, 5)
	go forwardErrors(participantJoinedSub.Err(), errorCh)
	go forwardErrors(participantLeftSub.Err(), errorCh)
	go forwardErrors(trackAddedSub.Err(), errorCh)
	go forwardErrors(eventToBackendSub.Err(), errorCh)
	go forwardErrors(eventToFrontendSub.Err(), errorCh)

	fmt.Println("Listening for contract events...")

	// Create a ticker for retrying lost connections
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Main event loop
	for {
		select {
		case err := <-errorCh:
			log.Printf("Error in subscription: %v", err)
			// Consider implementing reconnection logic here

		case event := <-participantJoinedCh:
			eventHandler.HandleParticipantJoined(event)

		case event := <-participantLeftCh:
			eventHandler.HandleParticipantLeft(event)

		case event := <-trackAddedCh:
			eventHandler.HandleTrackAdded(event)

		case event := <-eventToBackendCh:
			eventHandler.HandleEventToBackend(event)

		case event := <-eventToFrontendCh:
			// Just log these events as they're being sent from our backend to frontend
			log.Printf("Event forwarded to frontend - Room: %s, Participant: %s",
				event.RoomId, event.Participant.Hex())

		case <-ticker.C:
			// Periodic health check
			checkConnections(client)

		case <-sigCh:
			fmt.Println("Received termination signal, shutting down...")
			return
		}
	}
}

// forwardErrors forwards errors from subscription error channels to the main error channel
func forwardErrors(from <-chan error, to chan<- error) {
	for err := range from {
		to <- err
	}
}

// checkConnections performs a periodic health check of connections
func checkConnections(client *ethclient.Client) {
	// Check if we're still connected to the blockchain
	_, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Printf("Connection issue detected: %v", err)
		// You could implement reconnection logic here
	}
}
