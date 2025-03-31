/**
 * Smart Contract Integration Module
 * 
 * This module handles all interactions with the DAppMeeting smart contract
 */

import { CONTRACT_ABI } from './abi.js';
import { auth } from './auth.js';

class SmartContractConnector {
    constructor() {
        // Contract configuration
        this.contractAddress = "0xf8e21a267aF1450FBC5E90e098Ff0a8DBd9AD6DC";
        this.contractABI = CONTRACT_ABI;

        // State variables
        this.contract = null;
        this.provider = null;
        this.signer = null;
        this.isInitialized = false;

        // Event listener maps
        this.eventListeners = new Map();
        this.cleanupFunctions = new Map();

        // Try to initialize immediately if wallet is already connected
        this._tryAutoInitialize();
    }

    /**
     * Try to automatically initialize if wallet is connected
     * @private
     */
    async _tryAutoInitialize() {
        // Check if auth object is already initialized
        if (auth.isConnected) {
            console.log("Wallet already connected, auto-initializing contract connection");
            try {
                await this.initialize();
            } catch (error) {
                console.warn("Auto-initialization failed:", error);
            }
        } else {
            console.log("Wallet not connected, skipping auto-initialization");
        }
    }

    /**
     * Initialize the smart contract connector
     */
    async initialize() {
        try {
            // Wait a moment to make sure auth has finished loading persisted state
            if (!auth.isConnected) {
                // Check if privateKey exists in localStorage and try to reconnect
                const privateKey = localStorage.getItem('privateKey');
                if (privateKey) {
                    console.log("Found private key in storage, attempting to connect wallet first");
                    const result = await auth.connect(privateKey);
                    if (!result.success) {
                        throw new Error("Failed to connect wallet using stored key");
                    }
                } else {
                    throw new Error("Wallet not connected. Please connect wallet first.");
                }
            }

            // Set up provider
            if (window.ethereum) {
                this.provider = new ethers.providers.Web3Provider(window.ethereum);
            } else {
                // Fallback to BSC testnet
                this.provider = new ethers.providers.JsonRpcProvider("https://bsc-testnet-rpc.publicnode.com");
            }

            // Create wallet with stored private key
            const privateKey = localStorage.getItem('privateKey');
            if (!privateKey) {
                throw new Error("Private key not found in local storage");
            }

            // Create wallet and signer
            const wallet = new ethers.Wallet(privateKey, this.provider);
            this.signer = wallet;

            // Create contract instance
            this.contract = new ethers.Contract(
                this.contractAddress,
                this.contractABI,
                this.signer
            );

            console.log("Smart contract connector initialized successfully");
            this.isInitialized = true;
            return true;
        } catch (error) {
            console.error("Failed to initialize smart contract connector:", error);
            return false;
        }
    }

    /**
     * Ensures the connector is initialized
     */
    async ensureInitialized() {
        if (!this.isInitialized) {
            return await this.initialize();
        }
        return true;
    }

    /**
     * Create a new room in the smart contract
     * @param {string} roomId - The room ID
     */
    async createRoom(roomId) {
        await this.ensureInitialized();

        try {
            const tx = await this.contract.createRoom(roomId);
            await tx.wait();
            console.log(`Room created with ID: ${roomId}`);
            return { success: true, roomId };
        } catch (error) {
            console.error("Failed to create room:", error);
            throw error;
        }
    }

    /**
     * Start a call in the given room
     * @param {string} roomId - The room ID
     */
    async startCallSM(roomId) {
        await this.ensureInitialized();

        try {
            // This is a placeholder - the actual contract might not have this function.
            // This just ensures the room exists in the contract
            const roomExists = await this.checkRoomExists(roomId);

            if (!roomExists) {
                console.log(`Room ${roomId} doesn't exist, creating it...`);
                const tx = await this.contract.createRoom(roomId);
                await tx.wait();
                console.log(`Room created with ID: ${roomId}`);
            }

            return { success: true, roomId };
        } catch (error) {
            console.error("Failed to start call:", error);
            throw error;
        }
    }

    /**
     * Check if a room exists in the contract
     * @param {string} roomId - Room ID to check
     * @returns {Promise<boolean>} - Whether the room exists
     */
    async checkRoomExists(roomId) {
        await this.ensureInitialized();

        try {
            // Call the rooms mapping to check if room exists
            const roomInfo = await this.contract.rooms(roomId);
            // If roomId in the response is empty, room doesn't exist
            return roomInfo && roomInfo.roomId !== "";
        } catch (error) {
            console.error("Failed to check room existence:", error);
            return false;
        }
    }

    /**
     * Join a room in the smart contract
     * @param {string} roomId - The room ID
     * @param {string} name - User's name
     * @param {Array} initialTracks - Initial tracks array
     * @param {Object} sessionDescription - Session description for WebRTC
     * @returns {Promise<Object>} - Join result with session ID
     */
    async joinRoom(roomId, name, initialTracks = [], sessionDescription = {}) {
        await this.ensureInitialized();

        try {
            // For joining we need to prepare the SessionDescription as bytes
            const sessionDescBytes = ethers.utils.toUtf8Bytes(JSON.stringify(sessionDescription));

            // Call the joinRoom function with the necessary parameters
            const tx = await this.contract.joinRoom(
                roomId,
                name || "Anonymous",
                initialTracks.length > 0 ? initialTracks : [],
                sessionDescBytes
            );

            // Wait for the transaction to be mined
            const receipt = await tx.wait();
            console.log("Join room transaction receipt:", receipt);

            // Create a session ID (in a real implementation, this would come from the backend)
            const sessionID = `session-${Date.now()}-${Math.floor(Math.random() * 1000)}`;

            return {
                success: true,
                roomId,
                sessionID,
                transactionHash: receipt.transactionHash
            };
        } catch (error) {
            console.error("Failed to join room:", error);
            throw error;
        }
    }

    /**
     * Leave a room in the smart contract
     * @param {string} roomId - The room ID
     * @param {string} sessionId - The session ID
     * @returns {Promise<Object>} - Leave result
     */
    async leaveRoomHandle(roomId, sessionId) {
        await this.ensureInitialized();

        try {
            // Call the leaveRoom function
            const tx = await this.contract.leaveRoom(roomId);
            const receipt = await tx.wait();

            console.log("Leave room transaction receipt:", receipt);
            return {
                success: true,
                roomId,
                sessionId,
                transactionHash: receipt.transactionHash
            };
        } catch (error) {
            console.error("Failed to leave room:", error);
            throw error;
        }
    }

    /**
     * Send notification to a room
     * @param {string} roomId - The room ID
     * @param {string} message - Message to send
     */
    async sendNotificationToRoom(roomId, message) {
        await this.ensureInitialized();

        try {
            // Use forwardEventToBackend to send a message to all participants
            const eventData = ethers.utils.toUtf8Bytes(JSON.stringify({
                type: "notification",
                message,
                timestamp: Date.now()
            }));

            const tx = await this.contract.forwardEventToBackend(roomId, eventData);
            await tx.wait();

            return { success: true, roomId, message };
        } catch (error) {
            console.error("Failed to send notification:", error);
            throw error;
        }
    }

    /**
     * Compress and encode an object using zlib and base64
     * @param {Object} data - Data to compress
     * @returns {string} - Compressed data as base64 string
     */
    compressData(data) {
        try {
            // Convert data to JSON string
            const jsonString = JSON.stringify(data);
            // Convert string to Uint8Array
            const uint8Array = new TextEncoder().encode(jsonString);

            // Compress using pako (zlib implementation)
            const compressed = pako.deflate(uint8Array);

            // Convert compressed data to base64 and add zlib: prefix for backend
            return "zlib:" + btoa(String.fromCharCode.apply(null, compressed));
        } catch (error) {
            console.error("Error compressing data:", error);
            return JSON.stringify(data); // Fallback to uncompressed
        }
    }

    /**
     * Publish tracks to the smart contract
     * @param {string} roomId - The room ID
     * @param {Object} trackData - Track data including offer and tracks
     * @param {string} sessionId - The session ID
     */
    async publishTracksCompressed(roomId, trackData, sessionId) {
        await this.ensureInitialized();

        try {
            // Compress trackData to save gas
            const compressedData = this.compressData(trackData);
            if (sessionId === undefined || sessionId === null || sessionId === "") {
                sessionId = "unknown-session";
            }

            // Create the event data with compressed track information
            const eventData = ethers.utils.toUtf8Bytes(JSON.stringify({
                type: "publish-track",
                sessionId,
                compressedData,
                timestamp: Date.now()
            }));
            // calculate the gas limit
            const gasLimit = await this.contract.estimateGas.forwardEventToBackend(roomId, eventData);
            console.log("Estimated gas limit for publish tracks:", gasLimit.toString());

            // Send to the smart contract
            const tx = await this.contract.forwardEventToBackend(roomId, eventData, {
                gasLimit: gasLimit.add(100000) // Add a buffer to the estimated gas limit 
            });

            await tx.wait();

            console.log("Published tracks to smart contract");
            return { success: true, roomId, sessionId };
        } catch (error) {
            console.error("Failed to publish tracks:", error);
            throw error;
        }
    }

    /**
     * Listen for track published answer
     * @param {string} roomId - Room ID
     * @param {Function} callback - Callback function
     * @returns {Function} - Cleanup function
     */
    async listenForTrackPublishedAnswer(roomId, callback) {
        await this.ensureInitialized();

        try {
            // Listen for EventForwardedToFrontend events
            // Filter only by participant address, which is indexed - not by roomId
            const filter = this.contract.filters.EventForwardedToFrontend(
                null, // Don't filter by roomId since it's not indexed
                null
            );

            const listener = (eventRoomId, participant, eventData) => {
                // Only process if this event is for our address and room
                if (participant.toLowerCase() !== auth.userAddress.toLowerCase()) {
                    return;
                }

                // Check roomId manually in the listener since we can't filter on it
                if (roomId && eventRoomId !== roomId) {
                    return;
                }

                try {
                    // Parse the event data
                    const decodedData = JSON.parse(ethers.utils.toUtf8String(eventData));

                    // Check if this is a publish-track-response event
                    if (decodedData.type === "publish-track-response") {
                        callback(decodedData.cloudflareResponse);
                    }
                } catch (error) {
                    console.error("Error processing track published answer:", error);
                }
            };

            this.contract.on(filter, listener);

            // Return function to remove the listener
            return () => {
                this.contract.off(filter, listener);
                return true; // Listener was active and has been removed
            };
        } catch (error) {
            console.error("Error setting up track published answer listener:", error);
            return () => false; // Listener was never active
        }
    }

    /**
     * Pull tracks from a remote session
     * @param {string} roomId - Room ID
     * @param {string} sessionId - Local session ID
     * @param {string} remoteSessionId - Remote session ID to pull from
     * @param {string} trackName - Track name to pull
     */
    async pullTracksCompressed(roomId, sessionId, remoteSessionId, trackName) {
        await this.ensureInitialized();

        try {
            // Create the pull track request
            const pullData = {
                type: "pull-track",
                sessionId,
                tracks: [{
                    remoteSessionId,
                    trackName
                }],
                timestamp: Date.now()
            };

            // Compress data to save gas
            const compressedData = this.compressData(pullData);

            // Create the event data
            const eventData = ethers.utils.toUtf8Bytes(JSON.stringify({
                type: "pull-track",
                sessionId,
                compressedData,
                timestamp: Date.now()
            }));

            // Send to the smart contract
            const tx = await this.contract.forwardEventToBackend(roomId, eventData);
            await tx.wait();

            console.log(`Requested track pull: ${trackName} from ${remoteSessionId}`);
            return { success: true };
        } catch (error) {
            console.error("Failed to pull tracks:", error);
            throw error;
        }
    }

    /**
     * Listen for track pull answer
     * @param {string} roomId - Room ID
     * @param {string} sessionId - Session ID
     * @param {Function} callback - Callback function
     * @returns {Function} - Cleanup function
     */
    listenForTrackPullAnswer(roomId, sessionId, callback) {
        return this._listenForEventType(
            roomId,
            "pull-track-response",
            callback
        );
    }

    /**
     * Listen for renegotiation answer
     * @param {string} sessionId - Session ID
     * @param {Function} callback - Callback function
     * @returns {Function} - Cleanup function
     */
    ListenToAnswerRenegotiate(sessionId, callback) {
        return this._listenForEventType(
            null, // Room ID not needed for this filter
            "renegotiation-response",
            callback
        );
    }

    /**
     * Renegotiate WebRTC connection
     * @param {string} sessionId - Session ID
     * @param {string} sdp - SDP offer
     * @param {string} type - SDP type
     */
    async callToRenegotiate(sessionId, sdp, type) {
        await this.ensureInitialized();

        try {
            // Find the room ID from localStorage - we need this for the event
            const roomId = localStorage.getItem('roomId');
            if (!roomId) {
                throw new Error("No roomId found in localStorage");
            }

            // Create renegotiation data
            const renegotiationData = {
                type: "renegotiation",
                sessionId,
                sessionDescription: {
                    sdp,
                    type
                },
                timestamp: Date.now()
            };

            // Compress data
            const compressedData = this.compressData(renegotiationData);

            // Create the event data
            const eventData = ethers.utils.toUtf8Bytes(JSON.stringify({
                type: "renegotiation",
                sessionId,
                compressedData,
                timestamp: Date.now()
            }));

            // Send to the smart contract
            const tx = await this.contract.forwardEventToBackend(roomId, eventData);
            await tx.wait();

            console.log("Renegotiation request sent");
            return { success: true };
        } catch (error) {
            console.error("Failed to send renegotiation request:", error);
            throw error;
        }
    }

    /**
     * Notify other participants about available tracks
     * @param {string} roomId - Room ID
     * @param {string} userAddress - User's wallet address
     * @param {string} sessionId - Session ID
     * @param {Array} trackNames - Track names array
     */
    async notifyTrackPullComplete(roomId, userAddress, sessionId, trackNames) {
        await this.ensureInitialized();

        try {
            // Create notification data
            const notificationData = {
                type: "track-pull-complete",
                sessionId,
                trackNames,
                timestamp: Date.now()
            };

            // Create the event data
            const eventData = ethers.utils.toUtf8Bytes(JSON.stringify(notificationData));

            // Send to the smart contract
            const tx = await this.contract.forwardEventToBackend(roomId, eventData);
            await tx.wait();

            console.log("Track pull complete notification sent");
            return { success: true };
        } catch (error) {
            console.error("Failed to notify track pull complete:", error);
            throw error;
        }
    }

    /**
     * Listen for track pull complete notifications
     * @param {string} roomId - Room ID
     * @param {string} sessionId - Session ID
     * @param {Function} callback - Callback function
     * @returns {Function} - Cleanup function
     */
    listenForTrackPullComplete(roomId, sessionId, callback) {
        return this._listenForEventType(
            roomId,
            "track-pull-complete",
            callback
        );
    }

    /**
     * Listen for participant leave room events
     * @param {string} roomId - Room ID
     * @param {Function} callback - Callback function
     * @returns {Function} - Cleanup function
     */
    listenForParticipantLeaveRoom(roomId, callback) {
        return this._listenForEventType(
            roomId,
            "participant-left",
            callback
        );
    }

    /**
     * Helper method to listen for specific event types
     * @param {string} roomId - Room ID to filter on (optional)
     * @param {string} eventType - Event type to filter
     * @param {Function} callback - Callback function
     * @returns {Function} - Cleanup function
     */
    async _listenForEventType(roomId, eventType, callback) {
        await this.ensureInitialized();

        const listenerKey = `${roomId || 'global'}-${eventType}-${Date.now()}`;

        try {
            // Don't filter by any parameters since they're not indexed in the contract
            const filter = this.contract.filters.EventForwardedToFrontend(
                null, // Don't filter by roomId (not indexed)
                null  // Don't filter by participant (also not indexed)
            );

            const listener = (eventRoomId, participant, eventData) => {
                try {
                    // Check participant address manually in the listener
                    if (participant.toLowerCase() !== auth.userAddress.toLowerCase()) {
                        return;
                    }

                    // Check roomId manually if provided
                    if (roomId && eventRoomId !== roomId) {
                        return;
                    }

                    // Parse the event data
                    const decodedData = JSON.parse(ethers.utils.toUtf8String(eventData));

                    // Check if this matches our event type
                    if (decodedData.type === eventType) {
                        callback(decodedData);
                    }
                } catch (error) {
                    console.error(`Error processing ${eventType} event:`, error);
                }
            };

            // Store the listener and add it to the contract
            this.eventListeners.set(listenerKey, listener);
            this.contract.on(filter, listener);

            // Create cleanup function
            const cleanup = () => {
                if (this.eventListeners.has(listenerKey)) {
                    this.contract.off(filter, this.eventListeners.get(listenerKey));
                    this.eventListeners.delete(listenerKey);
                    return true; // Listener was active and has been removed
                }
                return false; // Listener was not active
            };

            // Store cleanup function
            this.cleanupFunctions.set(listenerKey, cleanup);

            return cleanup;
        } catch (error) {
            console.error(`Error setting up ${eventType} listener:`, error);
            return () => false; // Listener was never active
        }
    }

    /**
 * Get participants from the contract
 * @param {string} roomId - Room ID
 * @returns {Promise<Array>} - Array of participants
 */
    async getParticipantsFromContract(roomId) {
        await this.ensureInitialized();

        try {
            // Get detailed participant information from the contract
            const participantDetails = await this.contract.getRoomParticipantsDetails(roomId);

            if (!participantDetails || participantDetails.length === 0) {
                console.log('No participants found in room', roomId);
                return [];
            }
            console.log('Participant details:', participantDetails);
            // Format the participant details into the expected structure
            const participants = participantDetails.map(participant => {
                return {
                    sessionId: participant.sessionID || 'unknown-session',
                    name: participant.name || 'Anonymous',
                    walletAddress: participant.walletAddress,
                    // Map tracks to publishedTracks array
                    publishedTracks: participant.tracks.map(track => track.trackName).filter(name => name !== '')
                };
            });

            console.log('Retrieved participants from contract:', participants);
            return participants;
        } catch (error) {
            console.error("Failed to get participants from contract:", error);
            // In case of failure, return at least the current user
            const currentUser = {
                sessionId: localStorage.getItem('sessionId') || 'unknown-session',
                name: localStorage.getItem('username') || 'Anonymous',
                walletAddress: auth.userAddress,
                publishedTracks: []
            };
            return [currentUser];
        }
    }
}

// Create singleton instance
const connector = new SmartContractConnector();

export default connector;
