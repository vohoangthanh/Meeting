package handle

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	contract "dappmeetingnew/constract"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// TransactionRequest represents a transaction request in the queue
type TransactionRequest struct {
	Method       string
	RoomID       string
	Participant  common.Address
	SessionID    string
	EventData    []byte
	ResponseChan chan *TransactionResponse
}

// TransactionResponse contains the result of a transaction
type TransactionResponse struct {
	TxHash common.Hash
	Error  error
}

// SMCallManager manages transactions with a single wallet and queue
type SMCallManager struct {
	client       *ethclient.Client
	contract     *contract.Contract
	privateKey   *ecdsa.PrivateKey
	address      common.Address
	requestQueue []TransactionRequest
	queueSignal  chan struct{}
	quitCh       chan struct{}
	busy         bool
	mu           sync.Mutex
}

// NewSMCallManager creates a new SMCallManager
func NewSMCallManager(client *ethclient.Client, contractAddress common.Address) (*SMCallManager, error) {
	// Load private key from environment
	if err := loadEnv(); err != nil {
		return nil, fmt.Errorf("failed to load environment: %v", err)
	}

	// Get private key from environment
	privateKeyStr := os.Getenv("WALLET_PRIVATE_KEY")
	if privateKeyStr == "" {
		return nil, fmt.Errorf("WALLET_PRIVATE_KEY not found in environment")
	}

	// Create ECDSA private key
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	// Get wallet address from private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Create contract instance
	contractInstance, err := contract.NewContract(contractAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate contract: %v", err)
	}

	// Create manager
	manager := &SMCallManager{
		client:       client,
		contract:     contractInstance,
		privateKey:   privateKey,
		address:      address,
		requestQueue: make([]TransactionRequest, 0),
		queueSignal:  make(chan struct{}, 1),
		quitCh:       make(chan struct{}),
		busy:         false,
	}

	// Start queue processor
	go manager.processQueue()

	return manager, nil
}

// loadEnv loads environment variables from .env files
func loadEnv() error {
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
			break
		}
	}

	if !loaded {
		return fmt.Errorf("no .env file found")
	}

	return nil
}

// ForwardEventToFrontend sends an event to the frontend through the smart contract
func (m *SMCallManager) ForwardEventToFrontend(roomID string, participant common.Address, eventData []byte) (common.Hash, error) {
	respChan := make(chan *TransactionResponse)
	request := TransactionRequest{
		Method:       "ForwardEventToFrontend",
		RoomID:       roomID,
		Participant:  participant,
		EventData:    eventData,
		ResponseChan: respChan,
	}

	// Add request to queue
	m.mu.Lock()
	m.requestQueue = append(m.requestQueue, request)
	m.mu.Unlock()

	// Signal the queue processor
	select {
	case m.queueSignal <- struct{}{}:
	default:
		// Signal already in queue
	}

	// Wait for response
	response := <-respChan
	return response.TxHash, response.Error
}

// SetParticipantSessionID sets the session ID for a participant
func (m *SMCallManager) SetParticipantSessionID(roomID string, participant common.Address, sessionID string) (common.Hash, error) {
	respChan := make(chan *TransactionResponse)
	request := TransactionRequest{
		Method:       "SetParticipantSessionID",
		RoomID:       roomID,
		Participant:  participant,
		SessionID:    sessionID,
		ResponseChan: respChan,
	}

	// Add request to queue
	m.mu.Lock()
	m.requestQueue = append(m.requestQueue, request)
	m.mu.Unlock()

	// Signal the queue processor
	select {
	case m.queueSignal <- struct{}{}:
	default:
		// Signal already in queue
	}

	// Wait for response
	response := <-respChan
	return response.TxHash, response.Error
}

// AddNewTrackAfterPublish adds track information to the smart contract after Cloudflare publishes a track
func (m *SMCallManager) AddNewTrackAfterPublish(roomID string, participant common.Address, sessionID string,
	trackName string, mid string, location string, isPublished bool) (common.Hash, error) {

	respChan := make(chan *TransactionResponse)
	request := TransactionRequest{
		Method:       "AddNewTrackAfterPublish",
		RoomID:       roomID,
		Participant:  participant,
		SessionID:    sessionID,
		EventData:    []byte(fmt.Sprintf("%s|%s|%s|%v", trackName, mid, location, isPublished)),
		ResponseChan: respChan,
	}

	// Add request to queue
	m.mu.Lock()
	m.requestQueue = append(m.requestQueue, request)
	m.mu.Unlock()

	// Signal the queue processor
	select {
	case m.queueSignal <- struct{}{}:
	default:
		// Signal already in queue
	}

	// Wait for response
	response := <-respChan
	return response.TxHash, response.Error
}

// processQueue continuously processes the transaction queue
func (m *SMCallManager) processQueue() {
	for {
		select {
		case <-m.quitCh:
			// Manager is being shut down
			return

		case <-m.queueSignal:
			// Process next request in the queue if not busy
			m.processNextRequest()

		case <-time.After(5 * time.Second):
			// Periodically check the queue in case we missed a signal
			m.mu.Lock()
			if len(m.requestQueue) > 0 && !m.busy {
				m.mu.Unlock()
				select {
				case m.queueSignal <- struct{}{}:
				default:
					// Signal already in queue
				}
			} else {
				m.mu.Unlock()
			}
		}
	}
}

// processNextRequest processes the next request in the queue if possible
func (m *SMCallManager) processNextRequest() {
	m.mu.Lock()

	// If already processing a transaction or no requests, return
	if m.busy || len(m.requestQueue) == 0 {
		m.mu.Unlock()
		return
	}

	// Get the next request and mark as busy
	request := m.requestQueue[0]
	m.requestQueue = m.requestQueue[1:]
	m.busy = true

	m.mu.Unlock()

	// Process the request
	go func(req TransactionRequest) {
		var txHash common.Hash
		var err error

		switch req.Method {
		case "ForwardEventToFrontend":
			txHash, err = m.executeTransaction(func(auth *bind.TransactOpts) (*types.Transaction, error) {
				return m.contract.ForwardEventToFrontend(auth, req.RoomID, req.Participant, req.EventData)
			})
		case "SetParticipantSessionID":
			txHash, err = m.executeTransaction(func(auth *bind.TransactOpts) (*types.Transaction, error) {
				return m.contract.SetParticipantSessionID(auth, req.RoomID, req.Participant, req.SessionID)
			})
		case "AddNewTrackAfterPublish":
			// Parse track data from EventData (format: trackName|mid|location|isPublished)
			parts := strings.Split(string(req.EventData), "|")
			if len(parts) != 4 {
				err = fmt.Errorf("invalid track data format")
				break
			}
			trackName := parts[0]
			mid := parts[1]
			location := parts[2]
			isPublished := parts[3] == "true"

			txHash, err = m.executeTransaction(func(auth *bind.TransactOpts) (*types.Transaction, error) {
				return m.contract.AddNewTrackAfterPublish(auth, req.RoomID, req.Participant, req.SessionID,
					trackName, mid, location, isPublished)
			})
		default:
			err = fmt.Errorf("unknown method: %s", req.Method)
		}

		// Send response
		req.ResponseChan <- &TransactionResponse{
			TxHash: txHash,
			Error:  err,
		}

		// Mark wallet as available again
		m.mu.Lock()
		m.busy = false
		m.mu.Unlock()

		// Check if there are more requests to process
		m.mu.Lock()
		hasMoreRequests := len(m.requestQueue) > 0
		m.mu.Unlock()

		if hasMoreRequests {
			// Signal to process the next request
			select {
			case m.queueSignal <- struct{}{}:
			default:
				// Signal already in queue
			}
		}
	}(request)
}

// executeTransaction sends a transaction using the wallet
func (m *SMCallManager) executeTransaction(txnFunc func(*bind.TransactOpts) (*types.Transaction, error)) (common.Hash, error) {
	// Create transaction options
	auth, err := m.createTransactionOpts()
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create transaction options: %v", err)
	}

	// Send the transaction
	tx, err := txnFunc(auth)
	if err != nil {
		return common.Hash{}, fmt.Errorf("transaction failed: %v", err)
	}

	// Wait for the transaction to be mined
	receipt, err := m.waitForReceipt(tx.Hash())
	if err != nil {
		return tx.Hash(), fmt.Errorf("error waiting for receipt: %v", err)
	}

	// Check transaction status
	if receipt.Status == 0 {
		return tx.Hash(), errors.New("transaction reverted")
	}

	return tx.Hash(), nil
}

// createTransactionOpts creates transaction options for sending transactions
func (m *SMCallManager) createTransactionOpts() (*bind.TransactOpts, error) {
	ctx := context.Background()

	// Get the latest nonce for the wallet address
	nonce, err := m.client.PendingNonceAt(ctx, m.address)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	// Get gas price
	gasPrice, err := m.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %v", err)
	}

	// Get the chain ID from the connected client
	chainID, err := m.client.ChainID(ctx)
	if err != nil {
		// Fallback to BSC Testnet chain ID if we can't get it from the client
		chainID = big.NewInt(97) // 97 is the chain ID for BSC Testnet
		log.Printf("Using fallback chain ID: %d", chainID)
	}

	// Create transaction options with chain ID
	auth, err := bind.NewKeyedTransactorWithChainID(m.privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %v", err)
	}

	auth.GasPrice = gasPrice
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = 3000000 // Set a reasonable gas limit

	return auth, nil
}

// waitForReceipt waits for a transaction to be mined and returns the receipt
func (m *SMCallManager) waitForReceipt(txHash common.Hash) (*types.Receipt, error) {
	ctx := context.Background()
	for {
		receipt, err := m.client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}

		// If we get an error other than "not found", return it
		if err.Error() != "not found" {
			return nil, err
		}

		// Otherwise, wait a bit and try again
		time.Sleep(2 * time.Second)
	}
}

// Close shuts down the manager
func (m *SMCallManager) Close() {
	close(m.quitCh)
}

// GetQueueLength returns the current length of the transaction queue
func (m *SMCallManager) GetQueueLength() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.requestQueue)
}
