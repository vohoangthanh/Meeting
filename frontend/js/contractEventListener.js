/**
 * Contract Event Listener
 * 
 * Module xử lý việc lắng nghe các sự kiện từ smart contract
 * Hoạt động độc lập với luồng chính và sử dụng mẫu thiết kế Observer
 */

import { CONTRACT_ABI } from './abi.js';

class ContractEventListener {
    constructor() {
        // Cấu hình
        this.contractAddress = "0x7255F5d5DedAe6f7fABeCA17aB176d083b362A9c";
        this.wsProvider = null;
        this.contract = null;
        this.isListening = false;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectInterval = 3000; // 3 giây
        
        // Lưu trữ callback cho từng loại sự kiện
        this.eventHandlers = {
            ParticipantJoined: new Set(),
            ParticipantLeft: new Set(),
            TrackAdded: new Set(),
            EventForwardedToBackend: new Set(),
            EventForwardedToFrontend: new Set()
        };
        
        this.currentSessionId = null;
        this.hasSessionId = false;
        
        // WebSocket endpoint
        this.wsEndpoint = "wss://bsc-testnet-rpc.publicnode.com";
    }
    
    /**
     * Khởi tạo WebSocket provider và contract
     * @private
     */
    _initializeProvider() {
        try {
            // Đóng kết nối cũ nếu có
            if (this.wsProvider) {
                this.wsProvider._websocket.close();
            }
            
            // Khởi tạo WebSocket provider mới
            this.wsProvider = new ethers.providers.WebSocketProvider(this.wsEndpoint);
            
            // Xử lý sự kiện đóng WebSocket để tự động kết nối lại
            this.wsProvider._websocket.onclose = () => {
                console.log("WebSocket connection closed. Attempting to reconnect...");
                this.isListening = false;
                this._attemptReconnect();
            };
            
            // Tạo instance của contract
            this.contract = new ethers.Contract(
                this.contractAddress,
                CONTRACT_ABI,
                this.wsProvider
            );
            
            this.reconnectAttempts = 0;
            return true;
        } catch (error) {
            console.error("Error initializing WebSocket provider:", error);
            return false;
        }
    }
    
    /**
     * Thử kết nối lại khi WebSocket bị ngắt
     * @private
     */
    _attemptReconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            console.error(`Failed to reconnect after ${this.maxReconnectAttempts} attempts`);
            return;
        }
        
        setTimeout(() => {
            this.reconnectAttempts++;
            console.log(`Reconnect attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts}...`);
            
            if (this._initializeProvider()) {
                console.log("Reconnected successfully!");
                this.startListening();
            } else {
                this._attemptReconnect();
            }
        }, this.reconnectInterval);
    }
    
    /**
     * Bắt đầu lắng nghe các sự kiện
     */
    async startListening() {
        if (this.isListening) return;
        
        if (!this.contract) {
            if (!this._initializeProvider()) {
                console.error("Failed to initialize provider, cannot start listening");
                return;
            }
        }
        
        // Bắt đầu lắng nghe từng loại sự kiện
        this._setupEventListener("ParticipantJoined");
        this._setupEventListener("ParticipantLeft");
        this._setupEventListener("TrackAdded");
        this._setupEventListener("EventForwardedToBackend");
        this._setupEventListener("EventForwardedToFrontend");
        
        this.isListening = true;
        console.log("Started listening to contract events via WebSocket");
    }
    
    /**
     * Set current session ID for track comparison
     * @param {string} sessionId - Session ID của người dùng hiện tại
     */
    setCurrentSessionId(sessionId) {
        this.currentSessionId = sessionId;
    }
    
    /**
     * Thiết lập listener cho một loại sự kiện cụ thể
     * @private
     * @param {string} eventName - Tên sự kiện
     */
    _setupEventListener(eventName) {
        this.contract.on(eventName, (...args) => {
            const event = args[args.length - 1]; // Event object is always last argument
            const params = args.slice(0, -1); // Remove event object from params

            if (eventName === 'TrackAdded') {
                const [roomId, participant, trackName, sessionId] = params;
                console.log('TrackAdded event detected:', {
                    roomId,
                    participant,
                    trackName,
                    sessionId,
                    currentSessionId: this.currentSessionId
                });

                // kiem tra address
                var userAdress = localStorage.getItem("wallet_address");
                if (userAdress && userAdress.toLowerCase() === participant.toLowerCase()) {
                    this.hasSessionId = true;
                    this.currentSessionId = sessionId;
                    localStorage.setItem("sessionId", sessionId);
                }

                // // Chỉ pull track nếu không phải là track của bản thân
                // if (sessionId !== this.currentSessionId) {
                //     console.log('Pulling track from other participant:', trackName);
                //     // Thông báo cho tất cả các callback đã đăng ký
                //     this.eventHandlers[eventName].forEach(callback => {
                //         callback(roomId, participant, trackName, sessionId, event);
                //     });
                // } else {
                //     console.log('Ignoring own track:', trackName);
                // }
            } else {
                // Xử lý các sự kiện khác như cũ
                this.eventHandlers[eventName].forEach(callback => {
                    callback(...params, event);
                });
            }
        });
    }
    
    /**
     * Dừng lắng nghe các sự kiện
     */
    stopListening() {
        if (!this.isListening || !this.contract) return;
        
        // Hủy đăng ký tất cả các sự kiện
        this.contract.removeAllListeners();
        
        // Đóng kết nối WebSocket
        if (this.wsProvider && this.wsProvider._websocket) {
            this.wsProvider._websocket.close();
        }
        
        this.isListening = false;
        console.log("Stopped listening to contract events");
    }
    
    /**
     * Đăng ký callback cho một sự kiện
     * @param {string} eventName - Tên sự kiện
     * @param {Function} callback - Hàm callback khi sự kiện xảy ra
     * @returns {Function} - Hàm để hủy đăng ký callback
     */
    on(eventName, callback) {
        if (!this.eventHandlers[eventName]) {
            this.eventHandlers[eventName] = new Set();
        }
        
        this.eventHandlers[eventName].add(callback);
        
        // Đảm bảo đã bắt đầu lắng nghe
        if (!this.isListening) {
            this.startListening();
        }
        
        // Trả về hàm để hủy đăng ký callback
        return () => {
            this.eventHandlers[eventName].delete(callback);
        };
    }
    
    /**
     * Lắng nghe sự kiện EventForwardedToFrontend
     * @param {string} roomId - ID phòng (không bắt buộc)
     * @param {string} address - Địa chỉ người dùng
     * @param {Function} callback - Hàm callback
     * @returns {Function} - Hàm để hủy đăng ký callback
     */
    listenForEventsToFrontend(roomId, address, callback) {
        return this.on("EventForwardedToFrontend", async (eventRoomId, participant, eventData, event) => {
            // Filter by user address
            if (address && participant.toLowerCase() !== address.toLowerCase()) {
                return;
            }
            
            // Filter by roomId if provided
            if (roomId && eventRoomId !== roomId) {
                return;
            }
            
            try {
                let data;
                // If data is in hex format (from contract), convert to Uint8Array
                if (typeof eventData === 'string' && eventData.startsWith('0x')) {
                    eventData = eventData.slice(2); // Remove '0x'
                    const byteArray = new Uint8Array(eventData.match(/.{1,2}/g).map(byte => parseInt(byte, 16)));
                    data = byteArray;
                } else if (eventData instanceof Uint8Array) {
                    data = eventData;
                }

                // Try to decompress the data
                try {
                    const decompressed = pako.inflate(data);
                    const decoder = new TextDecoder('utf-8');
                    const jsonString = decoder.decode(decompressed);
                    const decodedData = JSON.parse(jsonString);
                    console.log('Successfully decompressed event data:', decodedData);
                    callback(decodedData, event);
                } catch (decompressionError) {
                    console.warn('Decompression failed, trying direct decoding:', decompressionError);
                    // If decompression fails, try to decode as UTF-8 directly
                    const decoder = new TextDecoder('utf-8');
                    const jsonString = decoder.decode(data);
                    const decodedData = JSON.parse(jsonString);
                    callback(decodedData, event);
                }
            } catch (error) {
                console.error("Error processing EventForwardedToFrontend:", error);
            }
        });
    }

    /**
     * Lắng nghe sự kiện kiểu TrackPublishedAnswer
     * @param {string} roomId - ID phòng
     * @param {string} address - Địa chỉ người dùng
     * @param {Function} callback - Hàm callback nhận cloudflareResponse
     * @returns {Function} - Hàm để hủy đăng ký callback
     */
    listenForTrackPublishedAnswer(roomId, address, callback) {
        return this.listenForEventsToFrontend(roomId, address, (decodedData) => {
            if (decodedData.type === "publish-track-response") {
                callback(decodedData.cloudflareResponse);
            }
        });
    }
    
    /**
     * Lắng nghe sự kiện kiểu TrackPullComplete
     * @param {string} roomId - ID phòng
     * @param {string} sessionId - ID phiên
     * @param {Function} callback - Hàm callback
     * @returns {Function} - Hàm để hủy đăng ký callback
     */
    listenForTrackPullComplete(roomId, address, callback) {
        return this.listenForEventsToFrontend(roomId, address, (decodedData) => {
            if (decodedData.type === "track-pull-complete") {
                callback(decodedData);
            }
        });
    }
    
    /**
     * Lắng nghe sự kiện kiểu ParticipantLeft
     * @param {string} roomId - ID phòng
     * @param {Function} callback - Hàm callback
     * @returns {Function} - Hàm để hủy đăng ký callback
     */
    listenForParticipantLeaveRoom(roomId, callback) {
        return this.on("ParticipantLeft", (eventRoomId, participant, event) => {
            if (roomId && eventRoomId !== roomId) {
                return;
            }
            
            callback({
                type: "participant-left",
                roomId: eventRoomId,
                participant,
                sessionId: participant // Sử dụng địa chỉ làm sessionId mặc định
            }, event);
        });
    }
    
    /**
     * Lắng nghe sự kiện TrackAdded
     * @param {string} roomId - ID phòng
     * @param {Function} callback - Callback function để xử lý khi có track mới
     * @returns {Function} - Hàm để hủy đăng ký callback
     */
    listenForTrackAdded(roomId, callback) {
        return this.on("TrackAdded", (eventRoomId, participant, trackName, sessionId, event) => {
            // Nếu có roomId được chỉ định, chỉ xử lý sự kiện cho phòng đó
            if (roomId && eventRoomId !== roomId) {
                return;
            }
            
            // Gọi callback với thông tin track
            callback(participant, trackName, sessionId);
        });
    }
    
    /**
     * Listen for track pull answer events
     * @param {string} roomId - ID phòng
     * @param {string} address - Địa chỉ ví người dùng
     * @param {Function} callback - Callback function
     * @returns {Function} - Function để hủy đăng ký callback
     */
    listenForTrackPullAnswer(roomId, address, callback) {
        return this.listenForEventsToFrontend(roomId, address, (decodedData) => {
            if (decodedData.type === 'pull-track-response') {
                try {
                    console.log("Received pull track response:", decodedData.type);
                    // The data should already be decompressed by listenForEventsToFrontend
                    callback(decodedData.cloudflareResponse);
                } catch (error) {
                    console.error("Error processing track pull answer:", error);
                }
            }
        });
    }

    /**
     * Helper method to decompress data
     * @param {string|Uint8Array} compressedData - The compressed data to decompress
     * @returns {string|null} - Decompressed data as string or null if failed
     * @private
     */
    async _decompressData(compressedData) {
        try {
            // If data is in hex string format (starts with '0x'), convert to Uint8Array
            if (typeof compressedData === "string" && compressedData.startsWith("0x")) {
                compressedData = compressedData.slice(2); // Remove '0x'
                const byteArray = new Uint8Array(compressedData.match(/.{1,2}/g).map(byte => parseInt(byte, 16)));

                // Decompress data using pako
                const data = pako.inflate(byteArray);

                // Convert byte array to UTF-8 string
                const decoder = new TextDecoder("utf-8");
                const strData = decoder.decode(data);

                console.log("Decompressed Data:", strData);

                return strData;
            } else {
                console.error("Invalid compressedData format:", compressedData);
                return null;
            }
        } catch (error) {
            console.error("Error in _decompressData:", error);
            return null;
        }
    }

    /**
     * Listen for specific event type from backend
     * @param {string} eventType - Type of event to listen for
     * @param {string} roomId - Room ID
     * @param {string} userAddress - User's wallet address 
     * @param {Function} callback - Callback function for event
     */
    listenForEventType(eventType, roomId, userAddress, callback) {
        return this.on("EventForwardedToFrontend", (eventRoomId, participant, eventData) => {
            // Filter by roomId and address
            if (roomId && eventRoomId !== roomId) return;
            if (userAddress && participant.toLowerCase() !== userAddress.toLowerCase()) return;
            
            try {
                // Parse and decompress data
                let data;
                if (typeof eventData === 'string' && eventData.startsWith('0x')) {
                    const dataBytes = ethers.utils.arrayify(eventData);
                    const decompressed = pako.inflate(dataBytes, {to: 'string'});
                    data = JSON.parse(decompressed);
                } else {
                    data = JSON.parse(eventData);
                }

                // Log received event for debugging
                console.log(`Received ${eventType} event:`, data);

                // Only call callback if event type matches
                if (data.type === eventType) {
                    callback(data);
                }
            } catch (error) {
                console.error(`Error processing ${eventType} event:`, error);
            }
        });
    }

    /**
     * Listen for renegotiation response
     * @param {string} roomId - Room ID
     * @param {string} userAddress - User's wallet address
     * @param {Function} callback - Callback function
     */
    listenForRenegotiationResponse(roomId, userAddress, callback) {
        console.log('Setting up renegotiation response listener');
        return this.listenForEventType('renegotiation-response', roomId, userAddress, (data) => {
            console.log('Renegotiation response received:', data);
            callback(data);
        });
    }
}

// Tạo và export instance singleton
const eventListener = new ContractEventListener();
export default eventListener;