package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	contract "dappmeetingnew/constract" // Update this import path to match your module

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	// Update these values with your own configuration
	ethereumNodeURL = "wss://bsc-testnet-rpc.publicnode.com"       // Use WebSocket endpoint for subscriptions
	contractAddress = "0xEf497BdCD80Aaa420271F5D2eAd9C1E70c2930E0" // Replace with your deployed contract address
)

func main() {
	// Connect to Ethereum node
	client, err := ethclient.Dial(ethereumNodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network: %v", err)
	}
	defer client.Close()
	fmt.Println("Connected to Ethereum node")

	// Create a new instance of the contract binding
	address := common.HexToAddress(contractAddress)
	instance, err := contract.NewContract(address, client)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}
	fmt.Println("Contract instance created")

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

	// Create filter options (from latest block onwards)
	// filterOpts := &bind.FilterOpts{
	// 	Start:   0, // Filter from latest block
	// 	End:     nil,
	// 	Context: ctx,
	// }

	// Create watch options
	watchOpts := &bind.WatchOpts{
		Context: ctx,
	}

	// Start watching for ParticipantJoined events
	participantJoinedSub, err := instance.WatchParticipantJoined(watchOpts, participantJoinedCh)
	if err != nil {
		log.Fatalf("Failed to watch ParticipantJoined events: %v", err)
	}
	defer participantJoinedSub.Unsubscribe()

	// Start watching for ParticipantLeft events
	participantLeftSub, err := instance.WatchParticipantLeft(watchOpts, participantLeftCh)
	if err != nil {
		log.Fatalf("Failed to watch ParticipantLeft events: %v", err)
	}
	defer participantLeftSub.Unsubscribe()

	// Start watching for TrackAdded events
	trackAddedSub, err := instance.WatchTrackAdded(watchOpts, trackAddedCh)
	if err != nil {
		log.Fatalf("Failed to watch TrackAdded events: %v", err)
	}
	defer trackAddedSub.Unsubscribe()

	// Start watching for EventForwardedToBackend events
	eventToBackendSub, err := instance.WatchEventForwardedToBackend(watchOpts, eventToBackendCh)
	if err != nil {
		log.Fatalf("Failed to watch EventForwardedToBackend events: %v", err)
	}
	defer eventToBackendSub.Unsubscribe()

	// Start watching for EventForwardedToFrontend events
	eventToFrontendSub, err := instance.WatchEventForwardedToFrontend(watchOpts, eventToFrontendCh)
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
			handleParticipantJoined(event)

		case event := <-participantLeftCh:
			handleParticipantLeft(event)

		case event := <-trackAddedCh:
			handleTrackAdded(event)

		case event := <-eventToBackendCh:
			handleEventToBackend(event)

		case event := <-eventToFrontendCh:
			handleEventToFrontend(event)

		case <-ticker.C:
			// Periodic check or maintenance if needed

		case <-sigCh:
			fmt.Println("Received termination signal, shutting down...")
			return
		}
	}
}

// forwardErrors forwards errors from one channel to another
func forwardErrors(from <-chan error, to chan<- error) {
	for err := range from {
		to <- err
	}
}

// handleParticipantJoined processes ParticipantJoined events
func handleParticipantJoined(event *contract.ContractParticipantJoined) {
	fmt.Printf("Participant Joined - Room: %s, Address: %s\n",
		event.RoomId, event.Participant.Hex())
	fmt.Printf("Number of initial tracks: %d\n", len(event.InitialTracks))

	// Add your business logic for handling participant joins
	// For example, update your backend state, notify other participants, etc.
}

// handleParticipantLeft processes ParticipantLeft events
func handleParticipantLeft(event *contract.ContractParticipantLeft) {
	fmt.Printf("Participant Left - Room: %s, Address: %s\n",
		event.RoomId, event.Participant.Hex())

	// Add your business logic for handling participant leaves
}

// handleTrackAdded processes TrackAdded events
func handleTrackAdded(event *contract.ContractTrackAdded) {
	fmt.Printf("Track Added - Room: %s, Participant: %s, Track Name: %s\n",
		event.RoomId, event.Participant.Hex(), event.TrackName)

	// Add your business logic for handling new tracks
}

// handleEventToBackend processes EventForwardedToBackend events
func handleEventToBackend(event *contract.ContractEventForwardedToBackend) {
	fmt.Printf("Event To Backend - Room: %s, Sender: %s, Data: %s\n",
		event.RoomId, event.Sender.Hex(), event.EventData)

	// Process the event data from the frontend
	// This could include commands, status updates, etc.
}

// handleEventToFrontend processes EventForwardedToFrontend events
func handleEventToFrontend(event *contract.ContractEventForwardedToFrontend) {
	fmt.Printf("Event To Frontend - Room: %s, Participant: %s, Data: %s\n",
		event.RoomId, event.Participant.Hex(), event.EventData)

	// Here you would typically forward this to your WebSocket/API service
	// to relay to the appropriate frontend client
}

// Example function to query room participant count (not used in the event listening)
func getParticipantCount(contractInstance *contract.Contract, roomId string) (*big.Int, error) {
	opts := &bind.CallOpts{Context: context.Background()}
	return contractInstance.GetRoomParticipantsCount(opts, roomId)
}
